package main

var globalID uint64

type Room struct {
	id     int
	Name   string
	Events map[string][]string
	Exits  map[string]bool
	Items  []*Item
	Target string
}

const (
	RoomKitchen = "кухня"
	RoomStreet  = "улица"
	RoomHome    = "домой"
	RoomHall    = "коридор"
	RoomRoom    = "комната"
)

func NewRoom(name string, events map[string][]string, exits map[string]bool, items []*Item, target string) *Room {
	globalID++

	if exits == nil {
		exits = make(map[string]bool)
	}
	if items == nil {
		items = make([]*Item, 0)
	}

	return &Room{
		id:     int(globalID),
		Name:   name,
		Events: events,
		Exits:  exits,
		Items:  items,
		Target: target,
	}
}

func (r Room) CanGo(dest string) bool {
	open, ok := r.Exits[dest]
	return ok && open
}
