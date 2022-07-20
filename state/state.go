package state

import (
	"strings"

	"github.com/feel-easy/hole-server/consts"
	"github.com/feel-easy/hole-server/models"
	"github.com/feel-easy/hole-server/utils"
	"github.com/feel-easy/hole-server/utils/logs"
)

var states = map[consts.StateID]State{}

func init() {
	register(consts.StateWelcome, &welcome{})
	register(consts.StateHome, &home{})
}

func register(id consts.StateID, state State) {
	states[id] = state
}

type State interface {
	Next(user *models.User) (consts.StateID, error)
	Exit(user *models.User) consts.StateID
}

func Run(user *models.User) {
	user.State(consts.StateWelcome)
	defer func() {
		if err := recover(); err != nil {
			utils.PrintStackTrace(err)
		}
		logs.Infof("user %s state machine break up.\n", user)
	}()
	for {
		state := states[user.GetState()]
		stateId, err := state.Next(user)
		if err != nil {
			if err1, ok := err.(consts.Error); ok {
				if err1.Exit {
					stateId = state.Exit(user)
				}
			} else {
				logs.Error(err)
				state.Exit(user)
				break
			}
		}
		if stateId > 0 {
			user.State(stateId)
		}
	}
}

func isExit(signal string) bool {
	signal = strings.ToLower(signal)
	return isX(signal, "exit", "e")
}

func isLs(signal string) bool {
	return isX(signal, "ls")
}

func isX(signal string, x ...string) bool {
	signal = strings.ToLower(signal)
	for _, v := range x {
		if v == signal {
			return true
		}
	}
	return false
}
