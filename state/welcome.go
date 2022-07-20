package state

import (
	"bytes"
	"fmt"

	"github.com/feel-easy/hole-server/consts"
	"github.com/feel-easy/hole-server/models"
)

type welcome struct{}

func (*welcome) Next(user *models.User) (consts.StateID, error) {
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("Hi %s, Welcome to hole online!\n", user.Name))
	buf.WriteString("1.Login\n")
	buf.WriteString("2.Register\n")
	err := user.WriteString(buf.String())
	if err != nil {
		return 0, user.WriteError(err)
	}
	selected, err := user.AskForInt()
	if err != nil {
		return 0, user.WriteError(err)
	}
	switch selected {
	case 1:
		return consts.StateLogin, nil
	case 2:
		return consts.StateRegister, nil
	}
	return 0, user.WriteError(consts.ErrorsInputInvalid)
}

func (*welcome) Exit(user *models.User) consts.StateID {
	return 0
}
