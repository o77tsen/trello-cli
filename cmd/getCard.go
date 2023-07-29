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
var getCardCmd = &cobra.Command{
	Use:   "getCard",
	Short: "Select a card to view its data (name, desc, URL, labels)",
	Long:  `Select a card to view its data (name, desc, URL, labels)`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := rootCmd.Flags().GetString("name")
		getCard(name)
	},
}

func init() {
	rootCmd.AddCommand(getCardCmd)
}

type CardToGet struct {
	ID     string   `json:"id"`
	Name   string   `json:"Name"`
	Desc   string   `json:"Desc"`
	URL    string   `json:"url"`
	Labels []string `json:"labels"`
}

func getCard(name string) {
	board, err := trelloInstance.GetBoard(trelloClient.GetBoardID())
	if err != nil {
		log.Fatal(err)
	}

	cards, err := board.GetCards(trello.Defaults())
	if err != nil {
		log.Fatal(err)
	}

	cardGetDataList := filterGetCards(cards)
	if len(cardGetDataList) == 0 {
		fmt.Println("There are no cards to view.")
		return
	}

	if name != "" {
		getCardByName(name, cards)
	} else {
		getCardBySelect(cardGetDataList)
	}
}

func filterGetCards(cards []*trello.Card) []CardToGet {
	var cardDataList []CardToGet

	for _, card := range cards {
		if !card.Closed {
			cardToGet := CardToGet{
				ID:   card.ID,
				Name: card.Name,
			}

			cardDataList = append(cardDataList, cardToGet)
		}
	}

	return cardDataList
}

func getCardByName(name string, cards []*trello.Card) {
	var selectedCard *trello.Card
	var labels []string

	for _, card := range cards {
		for _, label := range card.Labels {
			labels = append(labels, label.Name)
		}

		if card.Name == name {
			selectedCard = card
			break
		}
	}

	if selectedCard == nil {
		fmt.Printf("Card name `%s` could not be found.\n", name)
		return
	}

	cardData := CardToGet{
		ID:   selectedCard.ID,
		Name: selectedCard.Name,
		Desc: selectedCard.Desc,
		URL:  selectedCard.URL,
		Labels: labels,
	}

	displayCardData(&cardData)
}

func getCardBySelect(cards []CardToGet) {
	selectedCardId, _, err := promptSelectGet(cards)
	if err != nil {
		log.Fatal(err)
	}

	selectedCard := &cards[selectedCardId]

	displayCardData(selectedCard)
}

func promptSelectGet(cards []CardToGet) (int, string, error) {
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "🚀 {{ .Name | cyan }}",
		Inactive: " {{ .Name | cyan }}",
		Selected: "You are viewing: {{ .Name | cyan }}",
	}

	prompt := promptui.Select{
		Label:     "Select a card to view",
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

func displayCardData(card *CardToGet) {
	cyan := promptui.Styler(promptui.FGCyan)

	fmt.Printf("%s: %s\n%s: %s\n%s: %s\n%s: %s\n%s: %v\n", cyan("ID"), card.ID, cyan("Name"), card.Name, cyan("Desc"), card.Desc, cyan("URL"), card.URL, cyan("Labels"), card.Labels)
}
