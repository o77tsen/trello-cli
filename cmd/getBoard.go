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

// getBoardCmd represents the getBoard command
var getBoardCmd = &cobra.Command{
	Use:   "getBoard",
	Short: "Get board data from your trello",
	Long:  `Get board data from your trello`,
	Run: func(cmd *cobra.Command, args []string) {
		getBoard()
	},
}

type Board struct {
	Name       string
	Desc       string
	LabelNames []string
	Lists      []string
	URL        string
}

func init() {
	rootCmd.AddCommand(getBoardCmd)
}

func getBoard() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file:", err)
	}

	appKey := os.Getenv("TRELLO_KEY")
	token := os.Getenv("TRELLO_TOKEN")
	boardId := os.Getenv("TRELLO_BOARD_ID")

	client := trello.NewClient(appKey, token)

	board, err := client.GetBoard(boardId, trello.Defaults())
	if err != nil {
		log.Fatal(err)
	}

	jsonData, err := json.MarshalIndent(board, "", "    ")
	if err != nil {
		log.Fatal("Error converting to JSON:", err)
	}

	fmt.Println(string(jsonData))
}