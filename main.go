package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Expected Usage: taskly <cmd> <parameters>")
		os.Exit(1)
	}
	cmd := os.Args[1]
	switch {
	case cmd == "add":
		if len(os.Args) != 3 {
			fmt.Fprintln(os.Stderr, "Task Title is Required")
			os.Exit(1)
		}
		task, err := createTask(os.Args[2])
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error while creating a task:", err)
			os.Exit(1)
		}
		fmt.Println("Task created: ", task)
	case cmd == "ls":
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

		fmt.Fprintln(w, "ID\tDone\tTitle")
		fmt.Fprintln(w, "--\t----\t----\t")

		tasks, err := getTasks()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error while getting the tasks")
			os.Exit(1)
		}
		for _, task := range tasks {
			var done string
			if task.IsCompleted {
				done = "[X]"
			} else {
				done = "[ ]"
			}
			fmt.Fprintf(w, "%d\t%s\t%s\n", task.ID, done, task.Title)
		}
		w.Flush()
	case cmd == "edit":
		if len(os.Args) != 4 {
			fmt.Fprintln(os.Stderr, "Expected usage: taskly edit <id> <value>")
			os.Exit(1)
		}
		var id int
		fmt.Sscanf(os.Args[2], "%d", &id)
		new_title := os.Args[3]

		task, err := updateTask(id, new_title)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(task)
	}
	// taskly add <task> => add a new task
	// taskly ls => lists all the tasks
	// taskly edit <id> <task> => edit a task with id
	// taskly delete <id> => delete a task with id
	// taskly check <id> => mark a task with id as complete
	// taskly uncheck <id> => mark a task with id as incomplete

}
