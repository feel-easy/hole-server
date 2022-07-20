package models

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

type RoomGame interface {
	delete()
}

type Game interface {
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
