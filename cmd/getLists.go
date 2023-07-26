/*
Copyright © 2023 o77tsen
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/adlio/trello"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// getListsCmd represents the getLists command
var getListsCmd = &cobra.Command{
	Use:   "getLists",
	Short: "Get lists data from your trello board",
	Long:  `Get lists data from your trello board`,
	Run: func(cmd *cobra.Command, args []string) {
		getLists()
	},
}

func init() {
	rootCmd.AddCommand(getListsCmd)
}

func getLists() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file:", err)
		os.Exit(1)
	}

	appKey := os.Getenv("TRELLO_KEY")
	token := os.Getenv("TRELLO_TOKEN")
	boardId := os.Getenv("TRELLO_BOARD_ID")

	client := trello.NewClient(appKey, token)

	board, err := client.GetBoard(boardId, trello.Defaults())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	jsonData, err := json.MarshalIndent(lists, "", "    ")
	if err != nil {
		log.Fatal("Error converting to JSON:", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonData))
}
