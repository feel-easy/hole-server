package models

var users = make(map[int]*User)
var rooms = make(map[int]*Room)

func getRoom(roomID int) *Room {
	return rooms[roomID]
}
func delRoom(roomID int) {
	delete(rooms, roomID)
}

func getUser(userID int) *User {
	return users[userID]
}
