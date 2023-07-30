package trelloClient

import (
    "log"
    "os"

    "github.com/adlio/trello"
    "github.com/joho/godotenv"
)

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading env file:", err)
    }
}

func NewTrelloClient() *trello.Client {
    appKey := os.Getenv("TRELLO_KEY")
    token := os.Getenv("TRELLO_TOKEN")

    client := trello.NewClient(appKey, token)
    return client
}

func GetBoardID() string {
    boardId := os.Getenv("TRELLO_BOARD_ID")
    return boardId
}