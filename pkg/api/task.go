package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"scheduler/pkg/db"
)

const tasksLimit = 50

type tasksResponse struct {
	Tasks []*db.Task `json:"tasks"`
}

func (a *App) taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		a.addTaskHandler(w, r)
	case http.MethodGet:
		a.getTaskHandler(w, r)
	case http.MethodPut:
		a.updateTaskHandler(w, r)
	case http.MethodDelete:
		a.deleteTaskHandler(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (a *App) tasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	search := strings.TrimSpace(r.FormValue("search"))

	tasks, err := a.store.GetTasks(tasksLimit, search)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, tasksResponse{
		Tasks: tasks,
	})
}

func (a *App) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := readTaskFromBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateTask(task, false); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	id, err := a.store.AddTask(task)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"id": strconv.FormatInt(id, 10),
	})
}

func (a *App) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.FormValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "не указан идентификатор")
		return
	}

	task, err := a.store.GetTask(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (a *App) updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	task, err := readTaskFromBody(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateTask(task, true); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := a.store.UpdateTask(task); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}

func (a *App) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.FormValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "не указан идентификатор")
		return
	}

	if err := a.store.DeleteTask(id); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}

func (a *App) doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimSpace(r.FormValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "не указан идентификатор")
		return
	}

	task, err := a.store.GetTask(id)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if strings.TrimSpace(task.Repeat) == "" {
		if err := a.store.DeleteTask(id); err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		writeJSON(w, http.StatusOK, map[string]string{})
		return
	}

	next, err := NextDate(normalizeDate(time.Now()), task.Date, task.Repeat)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := a.store.UpdateDate(id, next); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{})
}

func readTaskFromBody(r *http.Request) (*db.Task, error) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		return nil, err
	}

	task.ID = strings.TrimSpace(task.ID)
	task.Date = strings.TrimSpace(task.Date)
	task.Title = strings.TrimSpace(task.Title)
	task.Comment = strings.TrimSpace(task.Comment)
	task.Repeat = strings.TrimSpace(task.Repeat)

	return &task, nil
}

func validateTask(task *db.Task, needID bool) error {
	if needID && task.ID == "" {
		return simpleError("не указан идентификатор")
	}

	if task.Title == "" {
		return simpleError("не указан заголовок задачи")
	}

	return checkDate(task)
}

func checkDate(task *db.Task) error {
	now := normalizeDate(time.Now())

	if task.Date == "" {
		task.Date = formatDate(now)
	}

	taskDate, err := parseDate(task.Date)
	if err != nil {
		return err
	}

	var next string
	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if taskDate.Before(now) {
		if task.Repeat == "" {
			task.Date = formatDate(now)
		} else {
			task.Date = next
		}
	}

	return nil
}

type simpleError string

func (e simpleError) Error() string {
	return string(e)
}
