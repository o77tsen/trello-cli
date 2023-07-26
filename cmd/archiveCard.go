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

// archiveCardCmd represents the archiveCard command
var archiveCardCmd = &cobra.Command{
	Use:   "archiveCard",
	Short: "Archive a card from your trello board",
	Long:  `Archive a card from your trello board`,
	Run: func(cmd *cobra.Command, args []string) {
		archiveCard()
	},
}

type CardToArchive struct {
	ID   string `json:"id"`
	Name string `json:"Name"`
}

func init() {
	rootCmd.AddCommand(archiveCardCmd)
}

func archiveCard() {
	client := trelloClient.NewTrelloClient()

	boardId := os.Getenv("TRELLO_BOARD_ID")

	board, err := client.GetBoard(boardId, trello.Defaults())

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var cardDataList []CardToArchive

	for _, card := range cards {
		if !card.Closed {
			getCard := CardToArchive{
				ID:   card.ID,
				Name: card.Name,
			}

			cardDataList = append(cardDataList, getCard)
		}
	}

	if len(cardDataList) == 0 {
		fmt.Println("There are no cards to archive.")
		os.Exit(1)
	}

	selectedCardIdx, cardID, err := promptSelectArchive(cardDataList)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = archiveSingleCard(client, cardID)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Printf("Success: archived card %s", cardDataList[selectedCardIdx].Name)
}

func promptSelectArchive(cards []CardToArchive) (int, string, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "🚀 {{ .Name | cyan }}",
		Inactive: " {{ .Name | cyan }}",
		Selected: "You are viewing: {{ .Name | cyan }}",
	}

	prompt := promptui.Select{
		Label:     "Select a card to archive",
		Items:     cards,
		Templates: templates,
		Size:      6,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return -1, "", err
	}

	return idx, cards[idx].ID, nil
}

func archiveSingleCard(client *trello.Client, cardID string) error {
	card, err := client.GetCard(cardID, trello.Defaults())
	if err != nil {
		return err
	}

	return card.Archive()
}
