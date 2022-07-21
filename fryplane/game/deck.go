package game

import (
	"math/rand"
)

type Deck struct {
	planes []int
}

func NewDeck() *Deck {
	deck := &Deck{}
	fillDeck(deck)
	return deck
}

func (d *Deck) NoPlanes() bool {
	return len(d.planes) == 0
}

func (d *Deck) DrawOne() int {
	return d.Draw(1)[0]
}

func (d *Deck) Draw(amount int) []int {
	planes := d.planes[0:amount]
	d.planes = d.planes[amount:]
	return planes
}

func (d *Deck) BottomDrawOne() int {
	plane := d.planes[len(d.planes)-1]
	d.planes = d.planes[:len(d.planes)-1]
	return plane
}

func fillDeck(deck *Deck) {
	planes := make([]int, 0, 144)

	shuffleCards(planes)
	deck.planes = append(deck.planes, planes...)
}

func shuffleCards(planes []int) {
	rand.Shuffle(len(planes), func(i, j int) { planes[i], planes[j] = planes[j], planes[i] })
}
