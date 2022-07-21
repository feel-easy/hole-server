package game

type Player interface {
	PlayerID() int
	NickName() string
	Play(planes []int, gameState State) (int, error)
	drew(planes []int, gameState State) (int, []int, error)
}
