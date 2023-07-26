/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/adlio/trello"
	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// moveCardCmd represents the moveCard command
var moveCardCmd = &cobra.Command{
	Use:   "moveCard",
	Short: "Move a card to another list from your trello",
	Long: `Create a card to another list from your trello`,
	Run: func(cmd *cobra.Command, args []string) {
		moveCard()
	},
}

func init() {
	rootCmd.AddCommand(moveCardCmd)
}

func moveCard() {
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

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	cardNames := make([]string, len(cards))
	for i, card := range cards {
		cardNames[i] = card.Name
	}

	promptCard := promptui.Select{
		Label: "Select a card you want to move",
		Items: cardNames,
	}

	cardIdx, _, err := promptCard.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	selectedCard := cards[cardIdx]

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
		Label: "Select the list to move a card to",
		Items: listNames,
	}

	idx, _, err := promptList.Run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	selectedList := lists[idx]

	err = selectedCard.MoveToList(selectedList.ID, trello.Defaults())
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	fmt.Printf("Success: moved card %s to list %s\n", selectedCard.Name, selectedList.Name)
}