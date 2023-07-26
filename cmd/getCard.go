/*
Copyright © 2023 o77tsen
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/adlio/trello"
	"github.com/manifoldco/promptui"
	"github.com/o77tsen/trello-cli/client"
	"github.com/spf13/cobra"
)

// getCardCmd represents the getCard command
var getCardCmd = &cobra.Command{
	Use:   "getCard",
	Short: "Get a card data from your trello",
	Long:  `Get a card data from your trello`,
	Run: func(cmd *cobra.Command, args []string) {
		getCard()
	},
}

type SingleCardData struct {
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
	client := trelloClient.NewTrelloClient()

	boardId := os.Getenv("TRELLO_BOARD_ID")

	board, err := client.GetBoard(boardId, trello.Defaults())

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var cardDataList []SingleCardData

	for _, card := range cards {
		if !card.Closed {
			var labels []string
			for _, label := range card.Labels {
				labels = append(labels, label.Name)
			}

			singleCardData := SingleCardData{
				ID:     card.ID,
				Name:   card.Name,
				Desc:   card.Desc,
				URL:    card.URL,
				Labels: labels,
			}

			cardDataList = append(cardDataList, singleCardData)
		}
	}

	selectedCardIdx, _, err := promptSelectCard(cardDataList)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	displayCardData(cardDataList[selectedCardIdx])
}

func promptSelectCard(cards []SingleCardData) (int, string, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "🚀 {{ .Name | cyan }}",
		Inactive: " {{ .Name | cyan }}",
		Selected: "You are viewing: {{ .Name | cyan }}",
	}

	prompt := promptui.Select{
		Label:     "Select a card to view",
		Items:     cards,
		Templates: templates,
		Size:      10,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return -1, "", err
	}

	return idx, cards[idx].ID, nil
}

func displayCardData(singleCardData SingleCardData) {
	fmt.Printf("Card ID: %s\nName: %s\nDesc: %s\nURL: %s\nLabels: %v\n", singleCardData.ID, singleCardData.Name, singleCardData.Desc, singleCardData.URL, singleCardData.Labels)
}
