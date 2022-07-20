package consts

import "time"

type StateID int

const (
	_ StateID = iota
	StateWelcome
	StateLogin
	StateRegister
	StateHome
	StateJoin
	StateCreate
	StateWaiting
	StateGame
	StateRunFastGame
	StateUnoGame
	StateMahjong
)

type GameType int

const (
	_ GameType = iota
	Mahjong
)

type RoomState int

const (
	_ RoomState = iota
	Waiting
	Running
)

const (
	IsStart = "INTERACTIVE_SIGNAL_START"
	IsStop  = "INTERACTIVE_SIGNAL_STOP"

	MinPlayers = 2
	MaxPlayers = 3

	PlayTimeout        = 40 * time.Second
	PlayMahjongTimeout = 30 * time.Second
)

type Error struct {
	Code int
	Msg  string
	Exit bool
}

func (e Error) Error() string {
	return e.Msg
}

func NewErr(code int, exit bool, msg string) Error {
	return Error{Code: code, Exit: exit, Msg: msg}
}

var (
	ErrorsExist                  = NewErr(1, true, "Exist. ")
	ErrorsChanClosed             = NewErr(1, true, "Chan closed. ")
	ErrorsTimeout                = NewErr(1, false, "Timeout. ")
	ErrorsInputInvalid           = NewErr(1, false, "Input invalid. ")
	ErrorsChatUnopened           = NewErr(1, false, "Chat disabled. ")
	ErrorsAuthFail               = NewErr(1, true, "Auth fail. ")
	ErrorsRoomInvalid            = NewErr(1, true, "Room invalid. ")
	ErrorsGameTypeInvalid        = NewErr(1, false, "Game type invalid. ")
	ErrorsRoomPlayersIsFull      = NewErr(1, false, "Room players is fill. ")
	ErrorsRoomPassword           = NewErr(1, false, "Sorry! Password incorrect! ")
	ErrorsJoinFailForRoomRunning = NewErr(1, false, "Join fail, room is running. ")
	ErrorsGamePlayersInvalid     = NewErr(1, false, "Game players invalid. ")
	ErrorsPokersFacesInvalid     = NewErr(1, false, "Pokers faces invalid. ")
	ErrorsHaveToPlay             = NewErr(1, false, "Have to play. ")
	ErrorsMustHaveToPlay         = NewErr(1, false, "There is a hand that can be played and must be played. ")
	ErrorsEndToPlay              = NewErr(1, false, "Can only come out at the end. ")

	GameTypes = map[GameType]string{
		Mahjong: "Mahjong",
	}
	GameTypesIds = []GameType{Mahjong}
	RoomStates   = map[RoomState]string{
		Waiting: "Waiting",
		Running: "Running",
	}
)
