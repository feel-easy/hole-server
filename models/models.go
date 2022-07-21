package models

type RoomGame interface {
	delete()
}

func Broadcast(roomId int, msg string, exclude ...int) {
	room := getRoom(roomId)
	if room == nil {
		return
	}
	room.Lock()
	defer room.Unlock()
	room.broadcast(msg, exclude...)
}
