package main

var world *World

type World struct {
	Rooms     map[string]*Room
	StartRoom *Room
}

func NewWorld() *World {
	return &World{
		Rooms: make(map[string]*Room),
	}
}

func (w *World) SetStartRoom(name string) {
	if room, ok := w.Rooms[name]; ok {
		w.StartRoom = room
	}
}

func (w *World) AddRoom(name string, events map[string][]string, exits map[string]bool, items []*Item, target string) {
	if w.Rooms == nil {
		w.Rooms = make(map[string]*Room)
	}
	room := NewRoom(name, events, exits, items, target)
	w.Rooms[name] = room
}
