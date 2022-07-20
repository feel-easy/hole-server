package state

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/feel-easy/hole-server/consts"
	"github.com/feel-easy/hole-server/models"
)

type join struct{}

func (s *join) Next(user *models.User) (consts.StateID, error) {
	buf := bytes.Buffer{}
	rooms := models.GetRooms()
	buf.WriteString(fmt.Sprintf("%-10s%-10s%-10s%-10s\n", "ID", "Type", "Players", "State"))
	for _, room := range rooms {
		pwdFlag := ""
		if room.Password != "" {
			pwdFlag = "*"
		}
		buf.WriteString(fmt.Sprintf("%-10d%-10s%-10d%-10s\n", room.ID, pwdFlag+consts.GameTypes[room.Type], room.UserNumber(), consts.RoomStates[room.State]))
	}
	err := user.WriteString(buf.String())
	if err != nil {
		return 0, user.WriteError(err)
	}
	signal, err := user.AskForString()
	if err != nil {
		return 0, user.WriteError(err)
	}
	if IsExit(signal) {
		return s.Exit(user), nil
	}
	if IsLs(signal) {
		return consts.StateJoin, nil
	}
	roomID, err := strconv.ParseInt(signal, 10, 64)
	roomId := int(roomID)
	if err != nil {
		return 0, user.WriteError(consts.ErrorsRoomInvalid)
	}
	room := models.GetRoom(roomId)
	if room == nil {
		return 0, user.WriteError(consts.ErrorsRoomInvalid)
	}

	//房间存在密码，要求输入密码
	pwd := room.Password
	if pwd != "" {
		err = verifyPassword(user, pwd)
		if err != nil {
			return 0, user.WriteError(err)
		}
	}
	err = models.JoinRoom(roomId, user.ID)
	if err != nil {
		return 0, user.WriteError(err)
	}
	models.Broadcast(roomId, fmt.Sprintf("%s joined room! room current has %d users\n", user.Name, room.UserNumber()))
	return consts.StateWaiting, nil
}

func (*join) Exit(user *models.User) consts.StateID {
	return consts.StateHome
}

// 校验密码
func verifyPassword(user *models.User, pwd string) error {
	err := user.WriteString("Please input room password: \n")
	if err != nil {
		return err
	}
	password, err := user.AskForString()
	if err != nil {
		return err
	}
	if password != pwd {
		return consts.ErrorsRoomPassword
	}
	return nil
}
