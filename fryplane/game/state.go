package game

import (
	"fmt"
	"sort"
	"strings"

	"github.com/feel-easy/mahjong/tile"
)

type State struct {
	CurrentPlayerHand []int
	PlayerSequence    []*playerController
}

func (s State) String() string {
	var lines []string
	drew := s.CurrentPlayerHand[len(s.CurrentPlayerHand)-1]
	sort.Ints(s.CurrentPlayerHand)

	lines = append(lines, fmt.Sprintf("Your drew: %s ", tile.Tile(drew)))
	lines = append(lines, fmt.Sprintf("Your hand: %s \n", tile.ToTileString(s.CurrentPlayerHand)))
	return strings.Join(lines, "\n")
}
