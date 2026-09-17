package main

import (
	"bytes"
	"cmp"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
)

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"is_completed"`
}

const tasksFile = ".data/tasks.json"

func getNextTaskID(tasks []Task) int {
	if len(tasks) == 0 {
		return 1
	}
	maxTask := slices.MaxFunc(tasks, func(a, b Task) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return maxTask.ID + 1
}

func createTask(title string) (Task, error) {
	tasks, err := getTasks()
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(os.Stderr, "Unexpected error while reading the tasks.json file")
		return Task{}, err
	}
	new_task := Task{
		ID:          getNextTaskID(tasks),
		Title:       title,
		IsCompleted: false,
	}
	tasks = append(tasks, new_task)
	save_err := saveTask(tasks)
	if save_err != nil {
		return Task{}, save_err
	}
	return new_task, nil
}

func saveTask(tasks []Task) error {
	if err := os.MkdirAll(".data", 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tasksFile, data, 0644)
}

func getTasks() ([]Task, error) {
	data, err := os.ReadFile(tasksFile)
	if errors.Is(err, os.ErrNotExist) {
		return []Task{}, nil // No Tasks
	}
	if err != nil {
		return nil, err
	}
	if len(bytes.TrimSpace(data)) == 0 {
		return []Task{}, nil
	}
	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	return tasks, err
}

func getTaskByID(tasks []Task, id int) Task {
	var filtered []Task
	for _, task := range tasks {
		if task.ID == id {
			filtered = append(filtered, task)
		}
	}
	if len(filtered) == 0 {
		return Task{}
	}
	return filtered[0]
}

func updateTask(id int, new_title string) (Task, error) {
	tasks, getTaskErr := getTasks()
	if getTaskErr != nil {
		fmt.Fprintln(os.Stderr, "Error while getting the tasks")
		return Task{}, getTaskErr
	}
	taskToUpdate := getTaskByID(tasks, id)
	if taskToUpdate == (Task{}) {
		fmt.Fprintf(os.Stdout, "No task with id %d found", id)
		err := errors.New("No task with id")
		return Task{}, err
	}
	taskToUpdate.Title = new_title
	fmt.Println(tasks)
	saveTask(tasks)
	return taskToUpdate, nil
}
