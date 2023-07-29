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

// getCardCmd represents the getCard command
var getCardsCmd = &cobra.Command{
	Use:   "getCards",
	Short: "Get all cards from your trello board",
	Long:  `Get all cards from your trello board`,
	Run: func(cmd *cobra.Command, args []string) {
		getCards()
	},
}

func init() {
	rootCmd.AddCommand(getCardsCmd)
}

func getCards() {
	board, err := trelloInstance.GetBoard(trelloClient.GetBoardID())
	if err != nil {
		log.Fatal(err)
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Fatal(err)
	}

	for _, card := range cards {
		if !card.Closed {
			cyan := promptui.Styler(promptui.FGCyan)

			fmt.Printf("%s %s\n", cyan("-"), card.Name)
		}
	}
}
