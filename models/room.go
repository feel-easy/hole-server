package models

import (
	"sync"
	"time"

	"github.com/feel-easy/hole-server/consts"
	"github.com/feel-easy/hole-server/utils/logs"
)

type Room struct {
	sync.Mutex

	ID         int              `json:"id"`
	Type       consts.GameType  `json:"type"`
	Game       Game             `json:"roomGame"`
	State      consts.RoomState `json:"state"`
	Users      map[int]*User    `json:"users"`
	Robots     int              `json:"robots"`
	Creator    int              `json:"creator"`
	ActiveTime time.Time        `json:"activeTime"`
	MaxUsers   int              `json:"maxUsers"`
	Password   string           `json:"password"`
	EnableChat bool             `json:"enableChat"`
	Banker     int              `json:"banker"`
}

func (room *Room) UserNumber() int {
	return len(room.Users)
}

func (room *Room) removeUser(user *User) {
	if room == nil || user == nil {
		return
	}
	room.ActiveTime = time.Now()
	if _, ok := room.Users[user.ID]; ok {
		user.RoomID = 0
		delete(room.Users, user.ID)
		if len(room.Users) > 0 && room.Creator == user.ID {
			for k := range room.Users {
				room.Creator = k
				break
			}
		}
	}
	if len(room.Users) == 0 {
		room.delete()
	}
}

func (room *Room) Cancel() {
	if room.ActiveTime.Add(24 * time.Hour).Before(time.Now()) {
		logs.Infof("room %d is timeout 24 hours, removed.\n", room.ID)
		room.delete()
		return
	}
	living := false
	for id := range room.Users {
		if getUser(id).online {
			living = true
			break
		}
	}
	if !living {
		logs.Infof("room %d is not living, removed.\n", room.ID)
		room.delete()
	}
}

func (room *Room) broadcast(msg string, exclude ...int) {
	room.ActiveTime = time.Now()
	excludeSet := map[int]bool{}
	for _, exc := range exclude {
		excludeSet[exc] = true
	}

	for userId := range room.Users {
		if user := getUser(userId); user != nil && !excludeSet[userId] {
			_ = user.WriteString(">> " + msg)
		}
	}
}

func (room *Room) delete() {
	if room != nil {
		if room.Game != nil {
			room.Game.delete()
		}
		delRoom(room.ID)
	}
}
