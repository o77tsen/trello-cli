/*
Copyright © 2023 o77tsen
*/
package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/adlio/trello"
	"github.com/manifoldco/promptui"
	"github.com/o77tsen/trello-cli/client"
	"github.com/spf13/cobra"
)

// createCardCmd represents the createCard command
var createCardCmd = &cobra.Command{
	Use:   "createCard",
	Short: "Create a new card in a specific list on your trello board",
	Long:  `Create a new card in a specific list on your trello board`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := rootCmd.Flags().GetString("name")
		list, _ := rootCmd.Flags().GetString("list")
		createCard(name, list)
	},
}

var cardLabels []string
var cardDesc string

type GetListData struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func init() {
	createCardCmd.Flags().StringVarP(&cardDesc, "desc", "d", "", "Provide a description for this card")
	createCardCmd.Flags().StringSliceVarP(&cardLabels, "labels", "t", nil, "Provide labels for this card (separate with commas)")
	rootCmd.AddCommand(createCardCmd)
}

func createCard(name string, list string) {
	board, err := trelloInstance.GetBoard(trelloClient.GetBoardID())
	if err != nil {
		log.Fatal(err)
	}

	if list != "" && name != "" {
		newCard, err := createCardDirect(board, name, cardDesc, list, cardLabels)
		if err != nil {
			log.Fatal(err)
		}

		err = trelloInstance.CreateCard(newCard, trello.Defaults())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Success: created %s in list %s\n", newCard.Name, list)
	} else {
		selectedList := fetchListByName(board)
		if selectedList == nil {
			fmt.Printf("List `%s` could not be found.", list)
			return
		}

		newCard, err := createCardObj(board, selectedList)
		if err != nil {
			log.Fatal(err)
		}

		err = trelloInstance.CreateCard(newCard, trello.Defaults())
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Success: created card `%s` in list `%s`\n", newCard.Name, selectedList.Name)
	}
}

func fetchListByName(board *trello.Board) *trello.List {
	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		log.Fatal(err)
	}

	var getListData []GetListData
	for _, list := range lists {
		listData := GetListData{
			ID: list.ID,
			Name: list.Name,
		}

		getListData = append(getListData, listData)
	}

	selectedListId, _, err := promptSelectList(getListData)
	if err != nil {
		log.Fatal(err)
	}

	return lists[selectedListId]
}

func createCardDirect(board *trello.Board, cardTitle, cardDesc, listName string, cardLabels []string) (*trello.Card, error) {
	cardLabelsInput := []string{}

	if len(cardLabels) > 0 {
		for _, label := range cardLabels {
			labels := strings.Split(label, ", ")
			for _, lbl := range labels {
				trimmedLabel := strings.TrimSpace(lbl)
				labelID := getLabelID(board, trimmedLabel)
				if labelID != "" {
					cardLabelsInput = append(cardLabelsInput, labelID)
				}
			}
		}
	}

	selectedList := findListByName(board, listName)
	if selectedList == nil {
		return nil, fmt.Errorf("List not found: %s", listName)
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
		Name:     cardTitle,
		Desc:     cardDesc,
		Pos:      pos,
		IDList:   selectedList.ID,
		IDLabels: cardLabelsInput,
	}

	return newCard, nil
}

func findListByName(board *trello.Board, listName string) *trello.List {
	lists, err := board.GetLists(trello.Defaults())
	if err != nil {
		log.Fatal(err)
	}

	for _, list := range lists {
		if list.Name == listName {
			return list
		}
	}
	return nil
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
