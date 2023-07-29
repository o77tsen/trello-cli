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

// moveCardCmd represents the moveCard command
var moveCardCmd = &cobra.Command{
	Use:   "moveCard",
	Short: "Move a card to another list from your trello",
	Long:  `Create a card to another list from your trello`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := rootCmd.Flags().GetString("name")
		list, _ := rootCmd.Flags().GetString("list")
		moveCard(name, list)
	},
}

func init() {
	rootCmd.AddCommand(moveCardCmd)
}

type CardToMove struct {
	ID   string `json:"id"`
	Name string `json:"Name"`
}

func moveCard(name string, list string) {
	board, err := trelloInstance.GetBoard(trelloClient.GetBoardID())
	if err != nil {
		log.Fatal(err)
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Fatal(err)
	}

	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		log.Fatal(err)
	}

	cardMoveDataList := filterMoveCards(cards)
	if len(cardMoveDataList) == 0 {
		fmt.Println("There are no cards to move.")
		return
	}

	if name != "" {
		moveCardByName(name, list, lists, cardMoveDataList)
	} else {
		moveCardBySelect(cardMoveDataList, lists)
	}
}

func filterMoveCards(cards []*trello.Card) []CardToMove {
	var cardDataList []CardToMove

	for _, card := range cards {
		if !card.Closed {
			cardToMove := CardToMove{
				ID:   card.ID,
				Name: card.Name,
			}

			cardDataList = append(cardDataList, cardToMove)
		}
	}

	return cardDataList
}

func findListByInput(name string, lists []*trello.List) *trello.List {
	for _, listData := range lists {
		if listData.Name == name {
			return listData
		}
	}

	return nil
}

func findCardByName(name string, cards []CardToMove) *CardToMove {
	for _, cardData := range cards {
		if cardData.Name == name {
			return &cardData
		}
	}

	return nil
}

func moveCardByName(name string, list string, lists []*trello.List, cards []CardToMove) {
	selectedList := findListByInput(list, lists)
	if selectedList == nil {
		fmt.Printf("List `%s` could not be found.\n", list)
		return
	}

	selectedCard := findCardByName(name, cards)
	if selectedCard == nil {
		fmt.Printf("Card name `%s` could not be found.\n", name)
		return
	}

	err := moveSingleCard(trelloInstance, selectedList.ID, selectedCard.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Success: moved card `%s` to `%s`\n", selectedCard.Name, selectedList.Name)
}

func moveCardBySelect(cards []CardToMove, lists []*trello.List) {
	selectedCardId, cardID, err := promptSelectMoveCard(cards)
	if err != nil {
		log.Fatal(err)
	}

	selectedListId, _, err := promptSelectMoveList(lists)
	if err != nil {
		log.Fatal(err)
	}

	err = moveSingleCard(trelloInstance, lists[selectedListId].ID, cardID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Success: moved card `%s` to `%s`\n", cards[selectedCardId].Name, lists[selectedListId].Name)
}

func promptSelectMoveCard(cards []CardToMove) (int, string, error) {
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
		Size:      6,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return -1, "", err
	}

	return idx, cards[idx].ID, nil
}

func promptSelectMoveList(lists []*trello.List) (int, string, error) {
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
		Size:      6,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return -1, "", err
	}

	return idx, lists[idx].ID, nil
}

func moveSingleCard(client *trello.Client, listID string, cardID string) error {
	card, err := client.GetCard(cardID, trello.Defaults())
	if err != nil {
		return err
	}

	args := trello.Arguments{
		"idList": listID,
	}

	return card.Update(args)
}