package consts

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

const (
	IsStart = "INTERACTIVE_SIGNAL_START"
	IsStop  = "INTERACTIVE_SIGNAL_STOP"

	MinPlayers = 2
	MaxPlayers = 3

	RoomStateWaiting = 1
	RoomStateRunning = 2
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

	GameTypes    = map[int]string{}
	GameTypesIds = []int{}
	RoomStates   = map[int]string{
		RoomStateWaiting: "Waiting",
		RoomStateRunning: "Running",
	}
)
