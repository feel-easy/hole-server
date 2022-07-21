package game

// *  *  *  *  *  *  *  *  *  *
// *  *  X  *  *  *  *  *  *  *
// X  X  X  X  X  *  *  *  *  *
// *  *  X  *  *  *  *  *  *  *
// *  X  X  X  *  *  *  *  *  *
// *  *  *  *  *  *  *  *  *  *
// *  *  *  O  *  *  *  *  *  *
// *  O  O  O  O  O  *  *  *  *
// *  *  *  O  *  *  *  *  *  *
// *  *  O  O  O  *  *  *  *  *

//    0  1  2  3  4  5  6  7  8  9
// 0  O  O  X  O  O  O  O  O  O  O
// 1  X  X  X  X  X  O  O  O  O  O
// 2  O  O  X  O  O  O  O  *  O  O
// 3  O  X  X  X  O  *  *  *  *  *
// 4  O  O  O  O  O  O  O  *  O  O
// 5  O  O  O  O  O  O  *  *  *  O
// 6  O  O  O  O  O  O  O  O  O  O
// 7  O  O  O  O  O  O  O  O  O  O
// 8  O  O  O  O  O  O  O  O  O  O
// 9  O  O  O  O  O  O  O  O  O  O
type Game struct {
	players *PlayerIterator
	deck    *Deck
	pile    *Pile
}

func (g *Game) Players() *PlayerIterator {
	return g.players
}

func (g *Game) Deck() *Deck {
	return g.deck
}

func (g *Game) Pile() *Pile {
	return g.pile
}

func (g *Game) Next() *playerController {
	player := g.Players().Next()
	return player
}

func New(players []Player) *Game {
	return &Game{
		players: newPlayerIterator(players),
		deck:    NewDeck(),
		pile:    NewPile(),
	}
}

func (g *Game) Current() *playerController {
	return g.players.Current()
}

func (g Game) ExtractState(player *playerController) *State {
	return &State{}
}
