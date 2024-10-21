package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
)

func main() {

	//loadTasks()
	app := &cli.App{
		Name:  "task",
		Usage: "manage your TODOs",
		Commands: []*cli.Command{
			{
				Name:   "add",
				Usage:  "Add a new task",
				Action: addTask,
			},
			{
				Name:   "do",
				Usage:  "Mark a task as complete",
				Action: completeTask,
			},
			{
				Name:   "list",
				Usage:  "List all tasks",
				Action: listTasks,
			},
		},
		Action: defaultAction,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	//saveTasks()

}

func defaultAction(c *cli.Context) error {
	fmt.Println("Welcome to the TODO app!")
	fmt.Println("Use 'task add <description>' to add a new task.")
	fmt.Println("Use 'task do <ID>' to mark a task as complete.")
	fmt.Println("Use 'task list' to see all tasks.")
	return nil
}

func addTask(c *cli.Context) error {
	taskDescription := c.Args().Slice()
	if len(taskDescription) == 0 {
		fmt.Println("No task description provided.")
		return nil
	}
	fullDescription := strings.Join(taskDescription, " ")

	db, err := openDB("tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = createBucket(db, "Tasks")
	if err != nil {
		log.Fatal(err)
	}

	putValue(db, "Tasks", fullDescription)

	fmt.Printf("Added %q to your task list.\n", fullDescription)
	return nil
}

func listTasks(c *cli.Context) error {

	db, err := openDB("tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = createBucket(db, "Tasks")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("You have the following tasks:\n")
	listAllValue(db, "Tasks")

	return nil
}

func completeTask(c *cli.Context) error {
	if c.NArg() == 0 {
		return cli.Exit("No task ID provided.", 1)
	}

	idStr := c.Args().Get(0)
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		return cli.Exit("Invalid task ID format.", 1)
	}

	db, err := openDB("tasks.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = createBucket(db, "Tasks")
	if err != nil {
		log.Fatal(err)
	}
	removeValue(db, "Tasks", idStr)

	fmt.Printf("You have completed the task with ID %d.\n", id)
	return nil
}

/*
type Task struct {
	ID          int
	Description string
	Completed   bool
}

var tasks []Task

const tasksFile = "tasks.json"

func addTask(c *cli.Context) error {
	taskDescription := c.Args().Slice()
	if len(taskDescription) == 0 {
		fmt.Println("No task description provided.")
		return nil
	}
	fullDescription := strings.Join(taskDescription, " ")

	id := len(tasks) + 1
	newTask := Task{
		ID:          id,
		Description: fullDescription,
		Completed:   false,
	}
	tasks = append(tasks, newTask)

	fmt.Printf("Added %q to your task list.\n", fullDescription)
	return nil
}

func completeTask(c *cli.Context) error {
	if c.NArg() == 0 {
		return cli.Exit("No task ID provided.", 1)
	}

	idStr := c.Args().Get(0)
	var id int
	if _, err := fmt.Sscanf(idStr, "%d", &id); err != nil {
		return cli.Exit("Invalid task ID format.", 1)
	}

	if id < 1 || id > len(tasks) {
		return cli.Exit("Task ID out of range.", 1)
	}

	tasks[id-1].Completed = true
	fmt.Printf("You have completed the task with ID %d.\n", id)
	return nil
}

func listTasks(c *cli.Context) error {
	fmt.Printf("You have the following tasks:\n")
	for i := range tasks {
		if !tasks[i].Completed {
			fmt.Printf("%d. %q \n", i+1, tasks[i].Description)
		}
	}
	return nil
}

func saveTasks() {
	file, err := os.OpenFile(tasksFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(tasks); err != nil {
		log.Fatal(err)
	}
}

func loadTasks() {
	file, err := os.Open(tasksFile)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		log.Fatal(err)
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if fi.Size() == 0 {
		return
	}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&tasks); err != nil {
		log.Fatal(err)
	}
}
*/
