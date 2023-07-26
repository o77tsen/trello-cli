/*
Copyright © 2023 o77tsen
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/adlio/trello"
	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// createCardCmd represents the createCard command
var createCardCmd = &cobra.Command{
	Use:   "createCard",
	Short: "Create a new card in a specific list on your trello board",
	Long:  `Create a new card in a specific list on your trello board`,
	Run: func(cmd *cobra.Command, args []string) {
		createCard()
	},
}

type GetListData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func init() {
	rootCmd.AddCommand(createCardCmd)
}

func createCard() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading env file:", err)
		os.Exit(1)
	}

	appKey := os.Getenv("TRELLO_KEY")
	token := os.Getenv("TRELLO_TOKEN")
	boardId := os.Getenv("TRELLO_BOARD_ID")

	client := trello.NewClient(appKey, token)

	board, err := client.GetBoard(boardId, trello.Defaults())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	var getListData []GetListData
	for _, list := range lists {
		listData := GetListData{
			ID:   list.ID,
			Name: list.Name,
		}

		getListData = append(getListData, listData)
	}

	selectedListIdx, _, err := promptSelectList(getListData)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	selectedList := lists[selectedListIdx]

	newCard, err := createCardObj(board, selectedList)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = client.CreateCard(newCard, trello.Defaults())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Printf("Success: created card %s in list %s\n", newCard.Name, selectedList.Name)
}

func createCardObj(board *trello.Board, selectedList *trello.List) (*trello.Card, error) {
	promptCardName := promptui.Prompt{
		Label: "Provide a name for this card ",
	}

	cardName, err := promptCardName.Run()
	if err != nil {
		return nil, err
	}

	promptCardDesc := promptui.Prompt{
		Label: "Provide a description for this card",
	}

	cardDesc, err := promptCardDesc.Run()
	if err != nil {
		return nil, err
	}

	cardDesc = strings.ReplaceAll(cardDesc, "\\n", "\n")

	promptCardLabels := promptui.Prompt{
		Label: "Provide labels for this card (Separate with commas)",
	}

	cardLabelsInput, err := promptCardLabels.Run()
	if err != nil {
		return nil, err
	}

	cardLabels := []string{}
	if cardLabelsInput != "" {
		labels := strings.Split(cardLabelsInput, ", ")
		for _, label := range labels {
			trimmedLabel := strings.TrimSpace(label)
			labelID := getLabelID(board, trimmedLabel)
			if labelID != "" {
				cardLabels = append(cardLabels, labelID)
			}
		}
	}

	cardsInList, err := selectedList.GetCards(trello.Defaults())
	if err != nil {
		return nil, err
	}

	var pos float64
	if len(cardsInList) > 0 {
		pos = cardsInList[0].Pos + 1.0
	} else {
		pos = 1.0
	}

	newCard := &trello.Card{
		Name:     cardName,
		Desc:     cardDesc,
		Pos:      pos,
		IDList:   selectedList.ID,
		IDLabels: cardLabels,
	}

	return newCard, nil
}

func getLabelID(board *trello.Board, labelName string) string {
	labels, err := board.GetLabels(trello.Defaults())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	for _, label := range labels {
		if label.Name == labelName {
			return label.ID
		}
	}

	return ""
}

func promptSelectList(lists []GetListData) (int, string, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "🚀 {{ .Name | cyan }}",
		Inactive: " {{ .Name | cyan }}",
		Selected: "You are viewing: {{ .Name | cyan }}",
	}

	prompt := promptui.Select{
		Label:     "Select a list to create the card in",
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
