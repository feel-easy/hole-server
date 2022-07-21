package game

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"github.com/feel-easy/hole-server/consts"
	"github.com/feel-easy/hole-server/models"
	"github.com/feel-easy/hole-server/utils/logs"
	"github.com/feel-easy/uno/card/color"
	"github.com/feel-easy/uno/event"
	"github.com/feel-easy/uno/game"
)

type Uno struct{}

func (g *Uno) Next(user *models.User) (consts.StateID, error) {
	room := models.GetRoom(user.RoomID)
	if room == nil {
		return 0, user.WriteError(consts.ErrorsExist)
	}
	game := room.RoomGame.(*models.UnoGame)
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf(
		"WELCOME TO %s%s%s!!!\n",
		color.Red.Paint("U"),
		color.Yellow.Paint("N"),
		color.Blue.Paint("O"),
	))
	buf.WriteString(fmt.Sprintf("Your Cards: %s\n", game.Game.GetPlayerCards(user.ID)))
	_ = user.WriteString(buf.String())
	for {
		if room.State == consts.Waiting {
			return consts.StateWaiting, nil
		}
		state := <-game.States[user.ID]
		switch state {
		case stateFirstCard:
			if msg := game.Game.PlayFirstCard(); msg != "" {
				models.Broadcast(room.ID, msg)
			}
			pc := game.Game.Players().Next()
			game.States[pc.ID()] <- statePlay
		case statePlay:
			err := handlePlayUno(room, user, game)
			if err != nil {
				logs.Error(err)
				return 0, err
			}
		case stateWaiting:
			return consts.StateWaiting, nil
		default:
			return 0, consts.ErrorsChanClosed
		}
	}
}

func (g *Uno) Exit(user *models.User) consts.StateID {
	room := models.GetRoom(user.RoomID)
	if room == nil {
		return consts.StateUnoGame
	}
	models.LeaveRoom(room.ID, user.ID)
	return consts.StateUnoGame
}

func handlePlayUno(room *models.Room, user *models.User, game *models.UnoGame) error {
	p := game.Game.Current()
	if p.ID() != user.ID {
		game.States[p.ID()] <- statePlay
		return nil
	}
	if !game.HavePlay(user) {
		pc := game.Game.Players().Next()
		game.States[pc.ID()] <- statePlay
	}
	gameState := game.Game.ExtractState(p)
	card, err := p.Play(gameState, game.Game.Deck())
	if err != nil || card == nil {
		event.PlayerPassed.Emit(event.PlayerPassedPayload{
			PlayerName: p.Name(),
		})
		pc := game.Game.Players().Next()
		game.States[pc.ID()] <- statePlay
		return err
	}
	game.Game.Pile().Add(card)
	event.CardPlayed.Emit(event.CardPlayedPayload{
		PlayerName: p.Name(),
		Card:       card,
	})
	if msg := game.Game.PerformCardActions(card); msg != "" {
		models.Broadcast(room.ID, msg)
	}
	if p.NoCards() || game.NeedExit() {
		models.Broadcast(room.ID, fmt.Sprintf("%s wins! \n", p.Name()))
		room.Lock()
		room.State = consts.Waiting
		room.Unlock()
		for _, userId := range game.Users {
			game.States[userId] <- stateWaiting
		}
		return nil
	}
	pc := game.Game.Players().Next()
	game.States[pc.ID()] <- statePlay
	return nil
}

func InitUnoGame(room *models.Room) (*models.UnoGame, error) {
	users := make([]int, 0)
	roomUsers := room.Users
	unoPlayers := make([]game.Player, 0)
	states := map[int]chan int{}
	for userId := range roomUsers {
		user := models.GetUser(userId)
		users = append(users, user.ID)
		unoPlayers = append(unoPlayers, models.NewUnoPlayer(user))
		states[userId] = make(chan int, 1)
	}
	rand.Seed(time.Now().UnixNano())
	unoGame := game.New(unoPlayers)
	unoGame.DealStartingCards()
	states[int(unoGame.Current().ID())] <- stateFirstCard
	return &models.UnoGame{
		Room:   room,
		Users:  users,
		States: states,
		Game:   unoGame,
	}, nil
}
