package cmd

import (
	"fmt"
	"strconv"
	"task/db"

	"github.com/spf13/cobra"
)

var doCmd = &cobra.Command{
	Use:   "do",
	Short: "marks a task as complete",
	Run: func(cmd *cobra.Command, args []string) {
		var ids []int
		for _, arg := range args {
			id, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("Failed to parse argument", err)
			} else {
				ids = append(ids, id)
			}
		}
		tasks, err := db.AllTasks()

		if err != nil {
			fmt.Println("Something went wrong", err)
			return
		}

		isValid := func(id int) func(int) bool {
			return func(len int) bool {
				return id > 0 && id <= len
			}
		}

		for _, id := range ids {
			if isValid(id)(len(tasks)) {
				task := tasks[id-1]
				err := db.DeleteTask(task.Key)
				if err != nil {
					fmt.Printf("Failed to mark \"%d\" as completed. Error: %s\n", id, err)
				} else {
					fmt.Printf("Marked \"%d\" as completed.\n", id)
				}
			} else {
				fmt.Printf("Invalid task id %d\n", id)
			}
		}

	},
}

func init() {
	RootCmd.AddCommand(doCmd)
}
