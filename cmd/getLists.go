/*
Copyright © 2023 o77tsen
*/
package cmd

import (
	"fmt"
	"log"

	"github.com/adlio/trello"
	"github.com/manifoldco/promptui"
	"github.com/o77tsen/trello-cli/client"
	"github.com/spf13/cobra"
)

// getListsCmd represents the getLists command
var getListsCmd = &cobra.Command{
	Use:   "getLists",
	Short: "Get all lists from your trello board",
	Long:  `Get all lists from your trello board`,
	Run: func(cmd *cobra.Command, args []string) {
		getLists()
	},
}

func init() {
	rootCmd.AddCommand(getListsCmd)
}

func getLists() {
	board, err := trelloInstance.GetBoard(trelloClient.GetBoardID())
	if err != nil {
		log.Fatal(err)
	}

	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		log.Fatal(err)
	}

	cyan := promptui.Styler(promptui.FGCyan)

	for _, list := range lists {
		fmt.Printf("%s %s\n", cyan("-"), list.Name)
	}
}
