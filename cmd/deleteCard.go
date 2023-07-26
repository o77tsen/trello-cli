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

// deleteCardCmd represents the deleteCard command
var deleteCardCmd = &cobra.Command{
	Use:   "deleteCard",
	Short: "Delete a card from your trello board",
	Long:  `Delete a card from your trello board`,
	Run: func(cmd *cobra.Command, args []string) {
		delCard()
	},
}

type GetCard struct {
	ID   string `json:"id"`
	Name string `json:"Name"`
}

func init() {
	rootCmd.AddCommand(deleteCardCmd)
}

func delCard() {
	client := trelloClient.NewTrelloClient()

	boardId := os.Getenv("TRELLO_BOARD_ID")
	
	board, err := client.GetBoard(boardId, trello.Defaults())

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var cardDataList []GetCard

	for _, card := range cards {
		if !card.Closed {
			getCard := GetCard{
				ID:   card.ID,
				Name: card.Name,
			}

			cardDataList = append(cardDataList, getCard)
		}
	}

	if len(cardDataList) == 0 {
		fmt.Println("There are no cards to delete.")
		os.Exit(1)
	}

	selectedCardIdx, cardID, err := promptSelect(cardDataList)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Printf("Are you sure you want to delete %s?", cardDataList[selectedCardIdx].Name)

	if !promptConfirm("Confirm deletion") {
		fmt.Println("Cancelled card deletion.")
		return
	}

	err = deleteCard(client, cardID)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Printf("Success: deleted card %s", cardDataList[selectedCardIdx].Name)
}

func promptSelect(cards []GetCard) (int, string, error) {
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

func deleteCard(client *trello.Client, cardID string) error {
	card, err := client.GetCard(cardID, trello.Defaults())
	if err != nil {
		return err
	}

	return card.Delete()
}
