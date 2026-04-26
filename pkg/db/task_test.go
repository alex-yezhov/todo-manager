package db

import (
	"path/filepath"
	"strconv"
	"testing"
)

func TestAddGetUpdateDeleteTask(t *testing.T) {
	setupTestDB(t)

	task := &Task{
		Date:    "20260426",
		Title:   "Test task",
		Comment: "Test comment",
		Repeat:  "d 1",
	}

	id, err := AddTask(task)
	if err != nil {
		t.Fatalf("AddTask returned error: %v", err)
	}

	got, err := GetTask(int64ToString(id))
	if err != nil {
		t.Fatalf("GetTask returned error: %v", err)
	}

	if got.Title != task.Title {
		t.Fatalf("got title %q, want %q", got.Title, task.Title)
	}
	if got.Comment != task.Comment {
		t.Fatalf("got comment %q, want %q", got.Comment, task.Comment)
	}
	if got.Repeat != task.Repeat {
		t.Fatalf("got repeat %q, want %q", got.Repeat, task.Repeat)
	}
	if got.Date != task.Date {
		t.Fatalf("got date %q, want %q", got.Date, task.Date)
	}

	got.Title = "Updated task"
	got.Comment = "Updated comment"
	got.Repeat = "y"
	got.Date = "20270426"

	if err := UpdateTask(got); err != nil {
		t.Fatalf("UpdateTask returned error: %v", err)
	}

	updated, err := GetTask(got.ID)
	if err != nil {
		t.Fatalf("GetTask after update returned error: %v", err)
	}

	if updated.Title != "Updated task" {
		t.Fatalf("title was not updated")
	}
	if updated.Comment != "Updated comment" {
		t.Fatalf("comment was not updated")
	}
	if updated.Repeat != "y" {
		t.Fatalf("repeat was not updated")
	}
	if updated.Date != "20270426" {
		t.Fatalf("date was not updated")
	}

	if err := UpdateDate(got.ID, "20280426"); err != nil {
		t.Fatalf("UpdateDate returned error: %v", err)
	}

	updatedDate, err := GetTask(got.ID)
	if err != nil {
		t.Fatalf("GetTask after UpdateDate returned error: %v", err)
	}

	if updatedDate.Date != "20280426" {
		t.Fatalf("got date %q, want %q", updatedDate.Date, "20280426")
	}

	if err := DeleteTask(got.ID); err != nil {
		t.Fatalf("DeleteTask returned error: %v", err)
	}

	_, err = GetTask(got.ID)
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestGetTasks(t *testing.T) {
	setupTestDB(t)

	_, _ = AddTask(&Task{
		Date:    "20260428",
		Title:   "Task 3",
		Comment: "",
		Repeat:  "",
	})
	_, _ = AddTask(&Task{
		Date:    "20260426",
		Title:   "Task 1",
		Comment: "alpha",
		Repeat:  "",
	})
	_, _ = AddTask(&Task{
		Date:    "20260427",
		Title:   "Task 2",
		Comment: "beta",
		Repeat:  "",
	})

	tasks, err := GetTasks(10, "")
	if err != nil {
		t.Fatalf("GetTasks returned error: %v", err)
	}

	if len(tasks) != 3 {
		t.Fatalf("got %d tasks, want 3", len(tasks))
	}

	if tasks[0].Title != "Task 1" {
		t.Fatalf("tasks are not sorted by date")
	}
	if tasks[1].Title != "Task 2" {
		t.Fatalf("tasks are not sorted by date")
	}
	if tasks[2].Title != "Task 3" {
		t.Fatalf("tasks are not sorted by date")
	}
}

func TestGetTasks_SearchByText(t *testing.T) {
	setupTestDB(t)

	_, _ = AddTask(&Task{
		Date:    "20260426",
		Title:   "Бассейн",
		Comment: "пойти вечером",
		Repeat:  "",
	})
	_, _ = AddTask(&Task{
		Date:    "20260427",
		Title:   "Магазин",
		Comment: "купить воду",
		Repeat:  "",
	})

	tasks, err := GetTasks(10, "Бассейн")
	if err != nil {
		t.Fatalf("GetTasks returned error: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("got %d tasks, want 1", len(tasks))
	}

	if tasks[0].Title != "Бассейн" {
		t.Fatalf("unexpected task title %q", tasks[0].Title)
	}
}

func TestGetTasks_SearchByDate(t *testing.T) {
	setupTestDB(t)

	_, _ = AddTask(&Task{
		Date:    "20260426",
		Title:   "Task A",
		Comment: "",
		Repeat:  "",
	})
	_, _ = AddTask(&Task{
		Date:    "20260427",
		Title:   "Task B",
		Comment: "",
		Repeat:  "",
	})

	tasks, err := GetTasks(10, "26.04.2026")
	if err != nil {
		t.Fatalf("GetTasks returned error: %v", err)
	}

	if len(tasks) != 1 {
		t.Fatalf("got %d tasks, want 1", len(tasks))
	}

	if tasks[0].Title != "Task A" {
		t.Fatalf("unexpected task title %q", tasks[0].Title)
	}
}

func setupTestDB(t *testing.T) {
	t.Helper()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	oldDB := DB

	t.Setenv(envDBFile, dbPath)

	database, err := InitDB()
	if err != nil {
		t.Fatalf("InitDB returned error: %v", err)
	}

	t.Cleanup(func() {
		_ = database.Close()
		DB = oldDB
	})
}

func int64ToString(v int64) string {
	return strconv.FormatInt(v, 10)
}
