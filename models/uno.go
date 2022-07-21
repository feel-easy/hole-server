package models

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/feel-easy/hole-server/consts"
	"github.com/feel-easy/uno/card"
	"github.com/feel-easy/uno/card/color"
	"github.com/feel-easy/uno/event"
	"github.com/feel-easy/uno/game"
)

type UnoGame struct {
	Room   *Room            `json:"room"`
	Users  []int            `json:"users"`
	States map[int]chan int `json:"states"`
	Game   *game.Game       `json:"game"`
}

func (ug *UnoGame) HavePlay(user *User) bool {
	for _, id := range ug.Users {
		if id == user.ID && user.online {
			return true
		}
	}
	return false
}

func (un *UnoGame) NeedExit() bool {
	return len(un.Room.Users) <= 1
}

func (un *UnoGame) delete() {
	if un != nil {
		for _, state := range un.States {
			close(state)
		}
	}
}

type UnoPlayer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewUnoPlayer(user *User) game.Player {
	return &UnoPlayer{
		ID:   user.ID,
		Name: user.Name,
	}
}

func (up *UnoPlayer) PlayerID() int {
	return up.ID
}

func (up *UnoPlayer) NickName() string {
	return up.Name
}

func contains(cards []card.Card, searchedCard card.Card) bool {
	for _, card := range cards {
		if card.Equal(searchedCard) {
			return true
		}
	}
	return false
}

func (up *UnoPlayer) NotifyCardsDrawn(cards []card.Card) {
	u := getUser(up.ID)
	getUser(u.ID).WriteString(fmt.Sprintf("You drew %s!\n", cards))
}

func (up *UnoPlayer) NotifyNoMatchingCardsInHand(lastPlayedCard card.Card, hand []card.Card) {
	u := getUser(up.ID)
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("%s, none of your cards match %s! \n", u.Name, lastPlayedCard))
	buf.WriteString(fmt.Sprintf("Your hand is %s \n", hand))
	getUser(u.ID).WriteString(buf.String())
}

func (up *UnoPlayer) OnFirstCardPlayed(payload event.FirstCardPlayedPayload) {
	u := getUser(up.ID)
	Broadcast(u.RoomID, fmt.Sprintf("First card is %s\n", payload.Card))
}

func (up *UnoPlayer) OnCardPlayed(payload event.CardPlayedPayload) {
	u := getUser(up.ID)
	Broadcast(u.RoomID, fmt.Sprintf("%s played %s!\n", payload.PlayerName, payload.Card))
}

func (up *UnoPlayer) OnColorPicked(payload event.ColorPickedPayload) {
	u := getUser(up.ID)
	Broadcast(u.RoomID, fmt.Sprintf("%s picked color %s!\n", payload.PlayerName, payload.Color))
}

func (up *UnoPlayer) OnUserPassed(payload event.PlayerPassedPayload) {
	u := getUser(up.ID)
	Broadcast(u.RoomID, fmt.Sprintf("%s passed!\n", payload.PlayerName))
}

func (up *UnoPlayer) PickColor(gameState game.State) color.Color {
	u := getUser(up.ID)
	for {
		u = getUser(u.ID)
		u.WriteString(fmt.Sprintf(
			"Select a color: %s, %s, %s or %s ? \n",
			color.Red,
			color.Yellow,
			color.Green,
			color.Blue,
		))
		colorName, err := u.AskForString(consts.PlayTimeout)
		if err != nil {
			if err == consts.ErrorsTimeout {
				return color.Red
			}
			u.WriteString(fmt.Sprintf("Unknown color '%s' \n", colorName))
			continue
		}
		chosenColor, err := color.ByName(strings.ToLower(colorName))
		if err != nil {
			u.WriteString(fmt.Sprintf("Unknown color '%s' \n", colorName))
			continue
		}
		return chosenColor
	}
}

func (up *UnoPlayer) Play(playableCards []card.Card, gameState game.State) (card.Card, error) {
	u := getUser(up.ID)
	Broadcast(u.RoomID, fmt.Sprintf("It's %s turn! \n", u.Name), u.ID)
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("It's your turn, %s! \n", u.Name))
	buf.WriteString(gameState.String())
	u.WriteString(buf.String())
	runeSequence := runeSequence{}
	cardOptions := make(map[string]card.Card)
	for _, card := range playableCards {
		label := string(runeSequence.next())
		cardOptions[label] = card
	}
	cardSelectionLines := []string{"Select a card to play:"}
	for label, card := range cardOptions {
		cardSelectionLines = append(cardSelectionLines, fmt.Sprintf("%s %s", label, card))
	}
	cardSelectionMessage := strings.Join(cardSelectionLines, " \n ") + " \n "
	for {
		u = getUser(u.ID)
		u.WriteString(cardSelectionMessage)
		selectedLabel, err := u.AskForString(consts.PlayTimeout)
		if err != nil {
			if err == consts.ErrorsTimeout {
				selectedLabel = "A"
			} else {
				return nil, err
			}
		}
		selectedCard, found := cardOptions[strings.ToUpper(selectedLabel)]
		if !found {
			u.BroadcastChat(fmt.Sprintf("%s say: %s\n", u.Name, selectedLabel))
			continue
		}
		if !contains(playableCards, selectedCard) {
			u.WriteString(fmt.Sprintf("Cheat detected! Card %s is not in %s's hand! \n", selectedCard, u.Name))
			continue
		}
		return selectedCard, nil
	}
}
