package models

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	rconsts "github.com/feel-easy/hole-server/consts"
	"github.com/feel-easy/mahjong/card"
	"github.com/feel-easy/mahjong/consts"
	"github.com/feel-easy/mahjong/event"
	"github.com/feel-easy/mahjong/game"
	"github.com/feel-easy/mahjong/tile"
)

const initialRune = 'A'

type runeSequence struct {
	currentRune rune
}

func (s *runeSequence) next() rune {
	if s.currentRune == 0 {
		s.currentRune = initialRune
	}
	currentRune := s.currentRune
	s.currentRune++
	return currentRune
}

type Mahjong struct {
	Room      *Room            `json:"room"`
	PlayerIDs []int            `json:"playerIds"`
	States    map[int]chan int `json:"states"`
	Game      *game.Game       `json:"game"`
}

func (game *Mahjong) delete() {
	if game != nil {
		for _, state := range game.States {
			close(state)
		}
	}
}

type OP struct {
	operation int
	tiles     []int
}
type MahjongPlayer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func NewPlayer(user *User) *MahjongPlayer {
	return &MahjongPlayer{
		ID:   user.ID,
		Name: user.Name,
	}
}

func (p *MahjongPlayer) PlayerID() int {
	return p.ID
}

func (p *MahjongPlayer) NickName() string {
	return p.Name
}

func (mp *MahjongPlayer) OnPlayTile(payload event.PlayTilePayload) {
	u := GetUser(mp.ID)
	u.WriteString(fmt.Sprintf("You play %s ! \n", tile.Tile(payload.Tile)))
	Broadcast(u.RoomID, fmt.Sprintf("%s PlayTile %s !\n", payload.PlayerName, tile.Tile(payload.Tile)), u.ID)
}

func (mp *MahjongPlayer) Take(tiles []int, gameState game.State) (int, []int, error) {
	u := GetUser(mp.ID)
	Broadcast(u.RoomID, fmt.Sprintf("It's %s take mahjong! \n", u.Name), u.ID)
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("It's your take mahjong, %s! \n", u.Name))
	buf.WriteString(gameState.String())
	u.WriteString(buf.String())
	askBuf := bytes.Buffer{}
	tileOptions := make(map[string]*OP)
	runeSequence := runeSequence{}
	if pvs, ok := gameState.SpecialPrivileges[u.ID]; ok {
		for _, pv := range pvs {
			switch pv {
			case consts.GANG:
				askBuf.WriteString("You can 杠!!!\n")
				label := string(runeSequence.next())
				ts := []int{gameState.LastPlayedTile, gameState.LastPlayedTile, gameState.LastPlayedTile}
				tileOptions[label] = &OP{
					operation: consts.GANG,
					tiles:     append(ts, gameState.LastPlayedTile),
				}
				askBuf.WriteString(fmt.Sprintf("%s:%s \n", label, tile.ToTileString(ts)))
			case consts.PENG:
				askBuf.WriteString("You can 碰!!!\n")
				label := string(runeSequence.next())
				ts := []int{gameState.LastPlayedTile, gameState.LastPlayedTile}
				tileOptions[label] = &OP{
					operation: consts.PENG,
					tiles:     append(ts, gameState.LastPlayedTile),
				}
				askBuf.WriteString(fmt.Sprintf("%s:%s \n", label, tile.ToTileString(ts)))
			case consts.CHI:
				askBuf.WriteString("You can 吃!!!\n")
				for _, ts := range card.CanChiTiles(tiles, gameState.LastPlayedTile) {
					label := string(runeSequence.next())
					tileOptions[label] = &OP{
						operation: consts.CHI,
						tiles:     append(ts, gameState.LastPlayedTile),
					}
					askBuf.WriteString(fmt.Sprintf("%s:%s \n", label, tile.ToTileString(ts)))
				}
			}
		}
	}
	label := string(runeSequence.next())
	askBuf.WriteString(fmt.Sprintf("%s:%s \n", label, "no"))
	tileOptions[label] = &OP{
		operation: 0,
		tiles:     []int{},
	}
	for {
		u = getUser(u.ID)
		u.WriteString(askBuf.String())
		selectedLabel, err := u.AskForString(rconsts.PlayMahjongTimeout)
		if err != nil {
			if err == rconsts.ErrorsExist {
				u.WriteString("Don't quit a good game！\n")
				continue
			}
			if err == rconsts.ErrorsTimeout {
				selectedLabel = "A"
			} else {
				return 0, nil, err
			}
		}
		selected, found := tileOptions[strings.ToUpper(selectedLabel)]
		if !found {
			u.BroadcastChat(fmt.Sprintf("%s say: %s\n", u.Name, selectedLabel))
			continue
		}
		return selected.operation, selected.tiles, nil
	}
}

func (mp *MahjongPlayer) Play(tiles []int, gameState game.State) (int, error) {
	u := GetUser(mp.ID)
	Broadcast(u.RoomID, fmt.Sprintf("It's %s turn! \n", u.Name), u.ID)
	buf := bytes.Buffer{}
	buf.WriteString(fmt.Sprintf("It's your turn, %s! \n", u.Name))
	buf.WriteString(gameState.String())
	u.WriteString(buf.String())
	askBuf := bytes.Buffer{}
	askBuf.WriteString("Select a tile to play:\n")
	runeSequence := runeSequence{}
	tileOptions := make(map[string]int)
	sort.Ints(tiles)
	for _, i := range tiles {
		label := string(runeSequence.next())
		tileOptions[label] = i
		askBuf.WriteString(fmt.Sprintf("%s:%s ", label, tile.Tile(i).String()))
	}
	askBuf.WriteString("\n")
	for {
		u = GetUser(u.ID)
		u.WriteString(askBuf.String())
		selectedLabel, err := u.AskForString(rconsts.PlayMahjongTimeout)
		if err != nil {
			if err == rconsts.ErrorsExist {
				u.WriteString("Don't quit a good game！\n")
				continue
			}
			if err == rconsts.ErrorsTimeout {
				selectedLabel = "A"
			} else {
				return 0, err
			}
		}
		selectedCard, found := tileOptions[strings.ToUpper(selectedLabel)]
		if !found {
			u.BroadcastChat(fmt.Sprintf("%s say: %s\n", u.Name, selectedLabel))
			continue
		}
		mp.OnPlayTile(event.PlayTilePayload{
			PlayerName: u.Name,
			Tile:       selectedCard,
		})
		return selectedCard, nil
	}
}
