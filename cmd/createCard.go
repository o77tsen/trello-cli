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
	Short: "Create a card from your trello",
	Long:  `Create a card from your trello`,
	Run: func(cmd *cobra.Command, args []string) {
		createCard()
	},
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

	listNames := make([]string, len(lists))
	for i, list := range lists {
		listNames[i] = list.Name
	}

	promptList := promptui.Select{
		Label: "Select the list to create a card",
		Items: listNames,
	}

	idx, _, err := promptList.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	selectedList := lists[idx]

	promptCardName := promptui.Prompt{
		Label: "Provide a name for this card",
	}

	cardName, err := promptCardName.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	promptCardDesc := promptui.Prompt{
		Label: "Provide a descripion for this card ",
	}

	cardDesc, err := promptCardDesc.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	cardDesc = strings.ReplaceAll(cardDesc, "\\n", "\n")

	promptCardLabels := promptui.Prompt{
		Label: "Provide card labels for this card (Separate them with commas) ",
	}

	cardLabelsInput, err := promptCardLabels.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
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
		log.Fatal(err)
		os.Exit(1)
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

	err = client.CreateCard(newCard, trello.Defaults())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Printf("Success: created card %s in list %s\n", cardName, selectedList.Name)
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
