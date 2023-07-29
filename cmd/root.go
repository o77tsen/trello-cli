/*
Copyright © 2023 o77tsen
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/adlio/trello"
	"github.com/o77tsen/trello-cli/client"
)

var trelloInstance * trello.Client

var rootCmd = &cobra.Command{
	Use:   "trello",
	Short: "Use trello to manage your trello board",
	Long: `Use trello to manage your trello board`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("name", "n", "", "card name")
	rootCmd.PersistentFlags().StringP("list", "l", "", "list name")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	trelloInstance = trelloClient.NewTrelloClient()
}