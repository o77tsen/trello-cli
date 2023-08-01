# trello-cli
this CLI allows you to interact with your trello boards, lists, and cards efficiently and seamlessly. it provides both flag-based and interactive prompts for a cleaner and user-friendly experience.

https://github.com/o77tsen/trello-cli/assets/88957235/39bb124f-cbc6-4989-ac09-8065640cdb5d

<img src="https://i.imgur.com/ASQ47py.png">

## commands
- `archiveCard` - archive a card from your board
- `createCard` - create a new card in a specific list from your board
- `deleteCard` - delete a card from your board
- `getCard` - select a card to view its data (title, description, URL, labels)
- `getCards` - get all cards from your board
- `getLists` - get all lists from your board
- `moveCard` - move a card from one list to another

## installation
1. clone the repo and navigate to the project directory
```shell
git clone https://github.com/o77tsen/trello-cli/
cd trello-cli
```

2. build the executable
```shell
go build
```

## configuration
obtain trello key & token [here](https://developer.atlassian.com/cloud/trello/guides/rest-api/api-introduction/)
```
TRELLO_KEY=
TRELLO_TOKEN=
TRELLO_BOARD_ID=
```

## usage

1. Create a new card
**with flags**
```
./trello-cli createCard -n "Card Name" -d "Card Description\nGoes here" -t "label-1, label-2" -l "List Name"
```

**without flags**
```
./trello-cli createCard
```

2. Delete a card
**with flags**
```
./trello-cli deleteCard -n "Card Name"
```
**without flags**
```
./trello-cli deleteCard
```

3. Archive a card
**with flags**
```
./trello-cli archiveCard -n "Card Name"
```

**without flags**
```
./trello-cli archiveCard
```

4. Move a card
**with flags**
```
./trello-cli moveCard -n "Card Name" -l "List Name"
```

**without flags**
```
./trello-cli moveCard
```
