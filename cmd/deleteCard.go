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

// deleteCardCmd represents the deleteCard command
var deleteCardCmd = &cobra.Command{
	Use:   "deleteCard",
	Short: "Delete a card from your trello board",
	Long:  `Delete a card from your trello board`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := rootCmd.Flags().GetString("name")
		deleteCard(name)
	},
}

func init() {
	rootCmd.AddCommand(deleteCardCmd)
}

type CardToDelete struct {
	ID   string `json:"id"`
	Name string `json:"Name"`
}

func deleteCard(name string) {
	board, err := trelloInstance.GetBoard(trelloClient.GetBoardID())
	if err != nil {
		log.Fatal(err)
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Fatal(err)
	}

	cardDeleteDataList := filterDeleteCards(cards)
	if len(cardDeleteDataList) == 0 {
		fmt.Println("There are no cards to delete.")
		return
	}

	if name != "" {
		deleteCardByName(name, cards)
	} else {
		deleteCardBySelect(cardDeleteDataList)
	}
}

func filterDeleteCards(cards []*trello.Card) []CardToDelete {
	var cardDataList []CardToDelete

	for _, card := range cards {
		if !card.Closed {
			cardToDelete := CardToDelete{
				ID:   card.ID,
				Name: card.Name,
			}

			cardDataList = append(cardDataList, cardToDelete)
		}
	}

	return cardDataList
}

func deleteCardByName(name string, cards []*trello.Card) {
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

	err := deleteSingleCard(trelloInstance, selectedCard.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Success: deleted card `%s`\n", selectedCard.Name)
}

func deleteCardBySelect(cards []CardToDelete) {
	selectedCardId, cardID, err := promptSelectDelete(cards)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Are you sure you want to delete `%s`?\n", cards[selectedCardId].Name)
	
	if !promptConfirm("Confirm deletion") {
		fmt.Println("Cancelled card deletion.")
		return
	}

	err = deleteSingleCard(trelloInstance, cardID)
	// err = deleteSingleCard(trelloInstance, cardID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Success: deleted card `%s`\n", cards[selectedCardId].Name)
}
func promptSelectDelete(cards []CardToDelete) (int, string, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "🚀 {{ .Name | cyan }}",
		Inactive: " {{ .Name | cyan }}",
		Selected: "You are viewing: {{ .Name | cyan }}",
	}

	prompt := promptui.Select{
		Label:     "Select a card to delete",
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

func promptConfirm(msg string) bool {
	confirmPrompt := promptui.Select{
		Label: msg,
		Items: []string{"Confirm", "Cancel"},
	}

	_, result, err := confirmPrompt.Run()
	if err != nil || result == "Cancel" {
		return false
	}

	return true
}

func deleteSingleCard(client *trello.Client, cardID string) error {
	card, err := client.GetCard(cardID, trello.Defaults())
	if err != nil {
		return err
	}

	return card.Delete()
}