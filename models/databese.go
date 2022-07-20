package models

import (
	"sync/atomic"
	"time"

	"github.com/feel-easy/hole-server/consts"
)

var roomIds int64 = 0
var users = make(map[int]*User)
var rooms = make(map[int]*Room)

func GetRooms() map[int]*Room {
	return rooms
}

func getRoom(roomID int) *Room {
	return rooms[roomID]
}

func GetRoom(roomID int) *Room {
	return getRoom(int(roomID))
}

func delRoom(roomID int) {
	delete(rooms, roomID)
}

func GetUser(userID int) *User {
	return getUser(int(userID))
}

func getUser(userID int) *User {
	return users[userID]
}

func setUser(user *User) {
	users[user.ID] = user
}

func CreateRoom(creator int, t consts.GameType) *Room {
	room := &Room{
		ID:         int(atomic.AddInt64(&roomIds, 1)),
		Type:       t,
		State:      consts.Waiting,
		Creator:    creator,
		Users:      map[int]*User{},
		ActiveTime: time.Now(),
		MaxUsers:   consts.MaxPlayers,
		EnableChat: true,
	}
	switch room.Type {
	case consts.Mahjong:
		room.MaxUsers = 4
	}
	rooms[room.ID] = room
	JoinRoom(room.ID, creator)
	return room
}

func LeaveRoom(roomId, userId int) {
	room := getRoom(roomId)
	if room != nil {
		room.Lock()
		defer room.Unlock()
		room.removeUser(getUser(userId))
	}
}

func JoinRoom(roomId, userId int) error {
	user := getUser(userId)
	if user == nil {
		return consts.ErrorsExist
	}
	room := getRoom(roomId)
	if room == nil {
		return consts.ErrorsRoomInvalid
	}
	if room.State == consts.Running {
		return consts.ErrorsJoinFailForRoomRunning
	}
	if room.UserNumber() >= room.MaxUsers {
		return consts.ErrorsRoomPlayersIsFull
	}
	room.Lock()
	defer room.Unlock()
	room.ActiveTime = time.Now()
	usersIds := room.Users
	if usersIds != nil {
		usersIds[userId] = user
		user.RoomID = roomId
	} else {
		room.delete()
		return consts.ErrorsRoomInvalid
	}
	return nil
}
