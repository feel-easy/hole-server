package game

import (
	"bytes"
	"fmt"
	"math/rand"
	"sort"
	"time"

	"github.com/feel-easy/hole-server/consts"
	"github.com/feel-easy/hole-server/models"
	"github.com/feel-easy/mahjong/card"
	mjconsts "github.com/feel-easy/mahjong/consts"
	"github.com/feel-easy/mahjong/event"
	"github.com/feel-easy/mahjong/game"
	"github.com/feel-easy/mahjong/tile"
	"github.com/feel-easy/mahjong/util"
	"github.com/feel-easy/mahjong/win"
)

type Mahjong struct{}

func (g *Mahjong) Next(user *models.User) (consts.StateID, error) {
	room := models.GetRoom(user.RoomID)
	if room == nil {
		return 0, user.WriteError(consts.ErrorsExist)
	}
	game := room.RoomGame.(*models.Mahjong)
	buf := bytes.Buffer{}
	buf.WriteString("WELCOME TO MAHJONG GAME!!! \n")
	buf.WriteString(fmt.Sprintf("%s is Banker! \n", models.GetUser(room.Banker).Name))
	buf.WriteString(fmt.Sprintf("Your Tiles: %s\n", game.Game.GetPlayerTiles(user.ID)))
	_ = user.WriteString(buf.String())
	for {
		if room.State == consts.Waiting {
			return consts.StateWaiting, nil
		}
		state := <-game.States[user.ID]
		switch state {
		case statePlay:
			err := handlePlayMahjong(room, user, game)
			if err != nil {
				return 0, err
			}
		case stateTakeCard:
			err := handleTake(room, user, game)
			if err != nil {
				return 0, err
			}
		case stateWaiting:
			return consts.StateWaiting, nil
		}
	}
}

func (g *Mahjong) Exit(user *models.User) consts.StateID {
	room := models.GetRoom(user.RoomID)
	if room == nil {
		return consts.StateMahjong
	}
	game := room.RoomGame.(*models.Mahjong)
	if game == nil {
		return consts.StateMahjong
	}
	for _, userId := range game.PlayerIDs {
		game.States[userId] <- stateWaiting
	}
	models.Broadcast(user.RoomID, fmt.Sprintf("user %s exit, game over! \n", user.Name))
	models.LeaveRoom(user.RoomID, user.ID)
	room.Lock()
	room.RoomGame = nil
	room.State = consts.Waiting
	room.Unlock()
	return consts.StateMahjong
}

func handleTake(room *models.Room, user *models.User, game *models.Mahjong) error {
	p := game.Game.Current()
	if p.ID() != user.ID {
		game.States[p.ID()] <- stateTakeCard
		return nil
	}
	if game.Game.Deck().NoTiles() {
		models.Broadcast(room.ID, "Game over but no winners!!! \n")
		room.Lock()
		room.RoomGame = nil
		room.State = consts.Waiting
		room.Unlock()
		for _, userId := range game.PlayerIDs {
			game.States[userId] <- stateWaiting
		}
		return nil
	}
	if t, ok := card.HaveGang(p.Hand()); ok {
		p.DarkGang(t)
		p.TryBottomDecking(game.Game.Deck())
		game.States[p.ID()] <- statePlay
		return nil
	}
	if card.CanGang(p.GetShowCardTiles(), p.LastTile()) {
		showCard := p.FindShowCard(p.LastTile())
		showCard.ModifyPongToKong(mjconsts.GANG, false)
		p.TryBottomDecking(game.Game.Deck())
		game.States[p.ID()] <- stateTakeCard
		return nil
	}
	gameState := game.Game.ExtractState(p)
	if len(gameState.SpecialPrivileges) > 0 {
		_, ok, err := p.Take(gameState, game.Game.Deck(), game.Game.Pile())
		if err != nil {
			return err
		}
		if ok {
			game.States[p.ID()] <- statePlay
			return nil
		}
		for {
			if gameState.OriginallyPlayer.ID() == p.ID() {
				p.TryTopDecking(game.Game.Deck())
				game.States[p.ID()] <- statePlay
				return nil
			}
			p = game.Game.Next()
		}
	}
	p.TryTopDecking(game.Game.Deck())
	game.States[p.ID()] <- statePlay
	return nil
}

func handlePlayMahjong(room *models.Room, user *models.User, game *models.Mahjong) error {
	p := game.Game.Current()
	if p.ID() != user.ID {
		game.States[p.ID()] <- statePlay
		return nil
	}
	gameState := game.Game.ExtractState(p)
	if win.CanWin(p.Hand(), p.GetShowCardTiles()) {
		tiles := p.Tiles()
		sort.Ints(tiles)
		models.Broadcast(room.ID, fmt.Sprintf("%s wins! \n%s \n", p.Name(), tile.ToTileString(tiles)))
		room.Lock()
		room.RoomGame = nil
		room.Banker = p.ID()
		room.State = consts.Waiting
		room.Unlock()
		for _, userId := range game.PlayerIDs {
			game.States[userId] <- stateWaiting
		}
		return nil
	}
	if _, ok := card.HaveGang(p.Hand()); ok {
		game.States[p.ID()] <- stateTakeCard
		return nil
	}
	if card.CanGang(p.GetShowCardTiles(), p.LastTile()) {
		game.States[p.ID()] <- stateTakeCard
		return nil
	}
	til, err := p.Play(gameState)
	if err != nil {
		return err
	}
	game.Game.Pile().Add(til)
	game.Game.Pile().SetLastPlayer(p)
	event.TilePlayed.Emit(event.TilePlayedPayload{
		PlayerName: p.Name(),
		Tile:       til,
	})
	pc := game.Game.Next()
	game.Game.Pile().SetOriginallyPlayer(pc)
	gameState = game.Game.ExtractState(p)
	if len(gameState.CanWin) > 0 {
		for _, p := range gameState.CanWin {
			tiles := append(p.Tiles(), gameState.LastPlayedTile)
			sort.Ints(tiles)
			models.Broadcast(room.ID, fmt.Sprintf("%s wins! \n%s \n", p.Name(), tile.ToTileString(tiles)))
		}
		room.Lock()
		room.RoomGame = nil
		room.Banker = gameState.CanWin[rand.Intn(len(gameState.CanWin))].ID()
		room.State = consts.Waiting
		room.Unlock()
		for _, userId := range game.PlayerIDs {
			game.States[userId] <- stateWaiting
		}
		return nil
	}
	if len(gameState.SpecialPrivileges) > 0 {
		pvID := pc.ID()
		flag := false
		for _, i := range []int{mjconsts.GANG, mjconsts.PENG, mjconsts.CHI} {
			for id, pvs := range gameState.SpecialPrivileges {
				if util.IntInSlice(i, pvs) {
					pvID = id
					flag = true
					break
				}
			}
			if flag {
				break
			}
		}
		for {
			if pc.ID() == pvID {
				game.States[pc.ID()] <- stateTakeCard
				return nil
			}
			pc = game.Game.Next()
		}
	}
	game.States[pc.ID()] <- stateTakeCard
	return nil
}

func InitMahjongGame(room *models.Room) (*models.Mahjong, error) {
	roomUsers := room.Users
	playerIDs := make([]int, 0, len(roomUsers))
	mjPlayers := make([]game.Player, 0, len(roomUsers))
	states := map[int]chan int{}
	for userId := range roomUsers {
		user := models.GetUser(userId)
		mjPlayers = append(mjPlayers, models.NewPlayer(user))
		playerIDs = append(playerIDs, user.ID)
		states[userId] = make(chan int, 1)
	}
	rand.Seed(time.Now().UnixNano())
	mahjong := game.New(mjPlayers)
	mahjong.DealStartingTiles()
	if room.Banker == 0 || !util.IntInSlice(room.Banker, playerIDs) {
		room.Banker = playerIDs[rand.Intn(len(playerIDs))]
	}
	for {
		if mahjong.Current().ID() == room.Banker {
			break
		}
		mahjong.Next()
	}
	states[mahjong.Current().ID()] <- stateTakeCard
	return &models.Mahjong{
		Room:      room,
		PlayerIDs: playerIDs,
		States:    states,
		Game:      mahjong,
	}, nil
}
