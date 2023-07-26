package trelloClient

import (
    "log"
    "os"

    "github.com/adlio/trello"
    "github.com/joho/godotenv"
)

func NewTrelloClient() *trello.Client {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading env file:", err)
        os.Exit(1)
    }

    appKey := os.Getenv("TRELLO_KEY")
    token := os.Getenv("TRELLO_TOKEN")

    client := trello.NewClient(appKey, token)
    return client
}

func GetBoard(client *trello.Client, boardID string) (*trello.Board, error) {
    board, err := client.GetBoard(boardID, trello.Defaults())
    if err != nil {
        return nil, err
    }

    return board, nil
}