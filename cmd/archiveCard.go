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

var archiveCardCmd = &cobra.Command{
	Use:   "archiveCard",
	Short: "Archive a card from your trello board",
	Long:  `Archive a card from your trello board`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := rootCmd.Flags().GetString("name")
		archiveCard(name)
	},
}

func init() {
	rootCmd.AddCommand(archiveCardCmd)
}

type CardToArchive struct {
	ID   string `json:"id"`
	Name string `json:"Name"`
}

func archiveCard(name string) {
	board, err := trelloInstance.GetBoard(trelloClient.GetBoardID())
	if err != nil {
		log.Fatal(err)
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Fatal(err)
	}

	cardArchiveDataList := filterArchiveCards(cards)
	if len(cardArchiveDataList) == 0 {
		fmt.Println("There are no cards to archive.")
		return
	}

	if name != "" {
		archiveCardByName(name, cards)
	} else {
		archiveCardBySelect(cardArchiveDataList)
	}
}

func filterArchiveCards(cards []*trello.Card) []CardToArchive {
	var cardDataList []CardToArchive

	for _, card := range cards {
		if !card.Closed {
			cardToArchive := CardToArchive{
				ID:   card.ID,
				Name: card.Name,
			}

			cardDataList = append(cardDataList, cardToArchive)
		}
	}

	return cardDataList
}

func archiveCardByName(name string, cards []*trello.Card) {
	var selectedCard *trello.Card

	for _, cardData := range cards {
		if cardData.Name == name {
			selectedCard = cardData
			break
		}
	}

	if selectedCard == nil {
		fmt.Printf("Card name `%s` could not be found.\n", name)
		return
	}

	err := archiveSingleCard(trelloInstance, selectedCard.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Success: archived card `%s`\n", selectedCard.Name)
}

func archiveCardBySelect(cards []CardToArchive) {
	selectedCardId, cardID, err := promptSelectArchive(cards)
	if err != nil {
		log.Fatal(err)
	}

	err = archiveSingleCard(trelloInstance, cardID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Success: archived card `%s`\n", cards[selectedCardId].Name)
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
