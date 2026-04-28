package db

import (
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddGetUpdateDeleteTask(t *testing.T) {
	store := setupTestDB(t)

	task := &Task{
		Date:    "20260426",
		Title:   "Test task",
		Comment: "Test comment",
		Repeat:  "d 1",
	}

	id, err := store.AddTask(task)
	require.NoError(t, err)
	require.NotZero(t, id)

	got, err := store.GetTask(int64ToString(id))
	require.NoError(t, err)
	require.NotNil(t, got)

	assert.Equal(t, task.Title, got.Title)
	assert.Equal(t, task.Comment, got.Comment)
	assert.Equal(t, task.Repeat, got.Repeat)
	assert.Equal(t, task.Date, got.Date)

	got.Title = "Updated task"
	got.Comment = "Updated comment"
	got.Repeat = "y"
	got.Date = "20270426"

	err = store.UpdateTask(got)
	require.NoError(t, err)

	updated, err := store.GetTask(got.ID)
	require.NoError(t, err)
	require.NotNil(t, updated)

	assert.Equal(t, "Updated task", updated.Title)
	assert.Equal(t, "Updated comment", updated.Comment)
	assert.Equal(t, "y", updated.Repeat)
	assert.Equal(t, "20270426", updated.Date)

	err = store.UpdateDate(got.ID, "20280426")
	require.NoError(t, err)

	updatedDate, err := store.GetTask(got.ID)
	require.NoError(t, err)
	require.NotNil(t, updatedDate)

	assert.Equal(t, "20280426", updatedDate.Date)

	err = store.DeleteTask(got.ID)
	require.NoError(t, err)

	_, err = store.GetTask(got.ID)
	require.Error(t, err)
}

func TestGetTasks(t *testing.T) {
	store := setupTestDB(t)

	_, err := store.AddTask(&Task{
		Date:    "20260428",
		Title:   "Task 3",
		Comment: "",
		Repeat:  "",
	})
	require.NoError(t, err)

	_, err = store.AddTask(&Task{
		Date:    "20260426",
		Title:   "Task 1",
		Comment: "alpha",
		Repeat:  "",
	})
	require.NoError(t, err)

	_, err = store.AddTask(&Task{
		Date:    "20260427",
		Title:   "Task 2",
		Comment: "beta",
		Repeat:  "",
	})
	require.NoError(t, err)

	tasks, err := store.GetTasks(10, "")
	require.NoError(t, err)
	require.Len(t, tasks, 3)

	assert.Equal(t, "Task 1", tasks[0].Title)
	assert.Equal(t, "Task 2", tasks[1].Title)
	assert.Equal(t, "Task 3", tasks[2].Title)
}

func TestGetTasks_SearchByText(t *testing.T) {
	store := setupTestDB(t)

	_, err := store.AddTask(&Task{
		Date:    "20260426",
		Title:   "Бассейн",
		Comment: "пойти вечером",
		Repeat:  "",
	})
	require.NoError(t, err)

	_, err = store.AddTask(&Task{
		Date:    "20260427",
		Title:   "Магазин",
		Comment: "купить воду",
		Repeat:  "",
	})
	require.NoError(t, err)

	tasks, err := store.GetTasks(10, "Бассейн")
	require.NoError(t, err)
	require.Len(t, tasks, 1)

	assert.Equal(t, "Бассейн", tasks[0].Title)
}

func TestGetTasks_SearchByDate(t *testing.T) {
	store := setupTestDB(t)

	_, err := store.AddTask(&Task{
		Date:    "20260426",
		Title:   "Task A",
		Comment: "",
		Repeat:  "",
	})
	require.NoError(t, err)

	_, err = store.AddTask(&Task{
		Date:    "20260427",
		Title:   "Task B",
		Comment: "",
		Repeat:  "",
	})
	require.NoError(t, err)

	tasks, err := store.GetTasks(10, "26.04.2026")
	require.NoError(t, err)
	require.Len(t, tasks, 1)

	assert.Equal(t, "Task A", tasks[0].Title)
}

func setupTestDB(t *testing.T) *Store {
	t.Helper()

	tmpDir := t.TempDir()
	dbPath := filepath.Join(tmpDir, "test.db")

	t.Setenv(envDBFile, dbPath)

	store, err := InitDB()
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = store.Close()
	})

	return store
}

func int64ToString(v int64) string {
	return strconv.FormatInt(v, 10)
}
