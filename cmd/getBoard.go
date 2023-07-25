/*
Copyright © 2023 o77tsen
*/
package cmd

import (
	// "encoding/json"
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

type LabelNames struct {
	Black  string `json:"black"`
	Green  string `json:"green"`
	Orange string `json:"orange"`
	Pink   string `json:"pink"`
	Purple string `json:"purple"`
	Red    string `json:"red"`
	Sky    string `json:"sky"`
}

type BoardData struct {
	Name       string     `json:"name"`
	Desc       string     `json:"desc"`
	URL        string     `json:"url"`
	ShortURL   string     `json:"shortUrl"`
	LabelNames LabelNames `json:"labelNames"`
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

	boardData := BoardData{
		Name:       board.Name,
		Desc:       board.Desc,
		URL:        board.URL,
		ShortURL:   board.ShortURL,
		LabelNames: LabelNames{
			Black: board.LabelNames.Black,
			Green: board.LabelNames.Green,
			Orange: board.LabelNames.Orange,
			Pink: board.LabelNames.Pink,
			Purple: board.LabelNames.Purple,
			Red: board.LabelNames.Red,
			Sky: board.LabelNames.Sky,
		},
	}

	printBoardDataFormatted(boardData)
}

func printBoardDataFormatted(boardData BoardData) {
	labels := fmt.Sprintf("\n- %s\n- %s\n- %s\n- %s\n- %s\n- %s\n- %s",
		boardData.LabelNames.Black, boardData.LabelNames.Green, boardData.LabelNames.Orange,
		boardData.LabelNames.Pink, boardData.LabelNames.Purple, boardData.LabelNames.Red,
		boardData.LabelNames.Sky,
	)

	fmt.Printf("Board Data for %s\nDesc: %s\nURL: %s\nShort URL: %s\nLabels: %s\n",
		boardData.Name, boardData.Desc, boardData.URL, boardData.ShortURL, labels)
}