/*
Copyright © 2023 o77tsen

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getCardCmd represents the getCard command
var getCardCmd = &cobra.Command{
	Use:   "getCard",
	Short: "Get card data from your trello board",
	Long: `Get card data from your trello board`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getCard called")
	},
}

func init() {
	rootCmd.AddCommand(getCardCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCardCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCardCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
