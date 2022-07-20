package menu

import (
	"bytes"

	"github.com/feel-easy/hole-server/consts"
	"github.com/feel-easy/hole-server/models"
)

type Home struct{}

func (*Home) Next(user *models.User) (consts.StateID, error) {
	buf := bytes.Buffer{}
	buf.WriteString("1.Join\n")
	buf.WriteString("2.New\n")
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
		return consts.StateJoin, nil
	case 2:
		return consts.StateCreate, nil
	}
	return 0, user.WriteError(consts.ErrorsInputInvalid)
}

func (*Home) Exit(user *models.User) consts.StateID {
	return 0
}
