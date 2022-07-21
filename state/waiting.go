package state

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/feel-easy/hole-server/consts"
	"github.com/feel-easy/hole-server/models"
	"github.com/feel-easy/hole-server/state/game"
)

type waiting struct{}

func (s *waiting) Next(user *models.User) (consts.StateID, error) {
	room := models.GetRoom(user.RoomID)
	if room == nil {
		return 0, consts.ErrorsExist
	}
	//_type 对接类别
	_type, access, err := waitingForStart(user, room)
	if err != nil {
		return 0, err
	}
	if access {
		switch room.Type {
		case consts.Mahjong:
			return consts.StateMahjong, nil
		case consts.Uno:
			return consts.StateUnoGame, nil
		default:
			return _type, nil
		}
	}
	return s.Exit(user), nil
}

func (*waiting) Exit(user *models.User) consts.StateID {
	room := models.GetRoom(user.RoomID)
	if room != nil {
		isOwner := room.Creator == user.ID
		models.LeaveRoom(room.ID, user.ID)
		models.Broadcast(room.ID, fmt.Sprintf("%s exited room! room current has %d users\n", user.Name, room.UserNumber()))
		if isOwner {
			newOwner := models.GetUser(room.Creator)
			models.Broadcast(room.ID, fmt.Sprintf("%s become new owner\n", newOwner.Name))
		}
	}
	return consts.StateHome
}

func waitingForStart(user *models.User, room *models.Room) (consts.StateID, bool, error) {
	access := false
	//对局类别
	_type := consts.StateGame
	user.StartTransaction()
	defer user.StopTransaction()
	for {
		signal, err := user.AskForStringWithoutTransaction(time.Second)
		if err != nil && err != consts.ErrorsTimeout {
			return consts.StateWaiting, access, err
		}
		if room.State == consts.Running {
			if room.Type == 4 {
				_type = consts.StateRunFastGame
			}
			access = true
			break
		}
		signal = strings.ToLower(signal)
		if signal == "ls" || signal == "v" {
			viewRoomUsers(room, user)
		} else if (signal == "start" || signal == "s") && room.Creator == user.ID && room.UserNumber() > 1 {
			//跑得快限制必须三人
			if room.Type == 4 && room.UserNumber() != 3 {
				err := user.WriteError(consts.ErrorsGamePlayersInvalid)
				if err != nil {
					return consts.StateWaiting, false, err
				}
				continue
			}
			access = true
			room.Lock()
			switch room.Type {
			default:
			case consts.Mahjong:
				room.RoomGame, err = game.InitMahjongGame(room)
			case consts.Uno:
				room.RoomGame, err = game.InitUnoGame(room)
			}
			if err != nil {
				room.Unlock()
				_ = user.WriteError(err)
				return consts.StateWaiting, access, err
			}
			room.State = consts.Running
			room.Unlock()

			break
		} else if strings.HasPrefix(signal, "set ") && room.Creator == user.ID {

			user.BroadcastChat(fmt.Sprintf("%s say: %s\n", user.Name, signal))
		} else if len(signal) > 0 {
			user.BroadcastChat(fmt.Sprintf("%s say: %s\n", user.Name, signal))
		}
	}
	return _type, access, nil
}

func viewRoomUsers(room *models.Room, currUser *models.User) {
	buf := bytes.Buffer{}

	buf.WriteString(fmt.Sprintf("Room ID: %d\n", room.ID))
	buf.WriteString(fmt.Sprintf("%-20s%-10s%-10s\n", "Name", "Score", "Title"))
	for userId := range room.Users {
		title := "user"
		if userId == room.Creator {
			title = "owner"
		}
		user := models.GetUser(userId)
		buf.WriteString(fmt.Sprintf("%-20s%-10d%-10s\n", user.Name, user.Score, title))
	}
	buf.WriteString("\nSettings:\n")
	pwd := room.Password
	if pwd != "" {
		if room.Creator != currUser.ID {
			pwd = "********"
		}
	} else {
		pwd = "off"
	}
	buf.WriteString(fmt.Sprintf("%-5s%-20v\n", "pwd", pwd))
	_ = currUser.WriteString(buf.String())
}
