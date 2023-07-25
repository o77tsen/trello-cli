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

// getCardCmd represents the getCard command
var getCardCmd = &cobra.Command{
	Use:   "getCard",
	Short: "Get card data from your trello board",
	Long:  `Get card data from your trello board`,
	Run: func(cmd *cobra.Command, args []string) {
		getCard()
	},
}

type CardData struct {
	ID     string   `json:"id"`
	Name   string   `json:"Name"`
	Desc   string   `json:"Desc"`
	URL    string   `json:"url"`
	Labels []string `json:"labels"`
}

func init() {
	rootCmd.AddCommand(getCardCmd)
}

func getCard() {
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

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Fatal(err)
	}

	var cardDataList []CardData

	for _, card := range cards {
		if !card.Closed {
			var labels []string
			for _, label := range card.Labels {
				labels = append(labels, label.Name)
			}
			cardData := CardData{
				ID:   card.ID,
				Name: card.Name,
				Desc: card.Desc,
				URL:  card.URL,
				Labels: labels,
			}

			cardDataList = append(cardDataList, cardData)
		}
	}

	if len(cardDataList) > 1 {
		cardDataList = cardDataList[1:]
	}

	jsonData, err := json.MarshalIndent(cardDataList, "", "    ")
	if err != nil {
		log.Fatal("Error converting to JSON:", err)
	}

	fmt.Println(string(jsonData))
}
