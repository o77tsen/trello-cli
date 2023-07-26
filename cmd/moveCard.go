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

// moveCardCmd represents the moveCard command
var moveCardCmd = &cobra.Command{
	Use:   "moveCard",
	Short: "Move a card to another list from your trello",
	Long:  `Create a card to another list from your trello`,
	Run: func(cmd *cobra.Command, args []string) {
		moveCard()
	},
}

type MovedCardData struct {
	ID   string `json:"id"`
	Name string `json:"Name"`
}

func init() {
	rootCmd.AddCommand(moveCardCmd)
}

func moveCard() {
	client := trelloClient.NewTrelloClient()

	boardId := os.Getenv("TRELLO_BOARD_ID")

	board, err := client.GetBoard(boardId, trello.Defaults())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var cardDataList []MovedCardData

	for _, card := range cards {
		if !card.Closed {
			movedCardData := MovedCardData{
				ID:   card.ID,
				Name: card.Name,
			}

			cardDataList = append(cardDataList, movedCardData)
		}
	}

	var listDataList []MovedCardData
	for _, list := range lists {
		listData := MovedCardData{
			ID:   list.ID,
			Name: list.Name,
		}

		listDataList = append(listDataList, listData)
	}

	selectedCardIdx, cardID, err := promptSelectMoveCard(cardDataList)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	selectedCard, err := client.GetCard(cardID, trello.Defaults())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	selectedListIdx, _, err := promptSelectMoveList(listDataList)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	selectedList := lists[selectedListIdx]

	err = selectedCard.MoveToList(selectedList.ID, trello.Defaults())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Printf("Success: moved card %s to list %s\n", cardDataList[selectedCardIdx].Name, selectedList.Name)
}

func promptSelectMoveCard(cards []MovedCardData) (int, string, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "🚀 {{ .Name | cyan }}",
		Inactive: " {{ .Name | cyan }}",
		Selected: "You are viewing: {{ .Name | cyan }}",
	}

	prompt := promptui.Select{
		Label:     "Select a card to move",
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

func promptSelectMoveList(lists []MovedCardData) (int, string, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "🚀 {{ .Name | cyan }}",
		Inactive: " {{ .Name | cyan }}",
		Selected: "You are viewing: {{ .Name | cyan }}",
	}

	prompt := promptui.Select{
		Label:     "Select a list to move the card to",
		Items:     lists,
		Templates: templates,
		Size:      10,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return -1, "", err
	}

	return idx, lists[idx].ID, nil
}
