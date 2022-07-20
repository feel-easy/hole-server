package state

import (
	"bytes"
	"fmt"

	"github.com/feel-easy/hole-server/consts"
	"github.com/feel-easy/hole-server/models"
)

type create struct{}

func (*create) Next(user *models.User) (consts.StateID, error) {
	gameType, err := askForGameType(user)
	if err != nil {
		return 0, err
	}
	// 创建房间
	room := models.CreateRoom(user.ID, consts.GameType(gameType))
	err = user.WriteString(fmt.Sprintf("Create room successful, id : %d\n", room.ID))
	if err != nil {
		return 0, user.WriteError(err)
	}
	err = models.JoinRoom(room.ID, user.ID)
	if err != nil {
		return 0, user.WriteError(err)
	}
	return consts.StateWaiting, nil
}

func (*create) Exit(user *models.User) consts.StateID {
	return consts.StateHome
}

// 询问游戏类型
func askForGameType(user *models.User) (gameType int, err error) {
	buf := bytes.Buffer{}
	buf.WriteString("Please select game type\n")
	for _, id := range consts.GameTypesIds {
		buf.WriteString(fmt.Sprintf("%d.%s\n", id, consts.GameTypes[id]))
	}
	err = user.WriteString(buf.String())
	if err != nil {
		return 0, user.WriteError(err)
	}
	gameType, err = user.AskForInt()
	if err != nil {
		return 0, user.WriteError(err)
	}
	// 游戏类型输入非法
	if _, ok := consts.GameTypes[consts.GameType(gameType)]; !ok {
		return 0, user.WriteError(consts.ErrorsGameTypeInvalid)
	}
	return
}
