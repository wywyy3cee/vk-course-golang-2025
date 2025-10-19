package main

import (
	"fmt"
	"sort"
	"strings"
)

var player *Player

type Player struct {
	Name        string
	State       string
	Actions     []string
	Invertory   map[*Item]bool
	CurrentRoom *Room
	ContainerOn bool
}

const (
	CmdLook           = "осмотреться"
	CmdGo             = "идти"
	CmdPutOn          = "надеть"
	CmdTake           = "взять"
	CmdUse            = "применить"
	CmdExamineTheItem = "осмотреть"
	CmdExit           = "выйти"
	CmdExits          = "выходы"
)

func NewPlayer(name, state string, actions []string, startRoom *Room) *Player {
	return &Player{
		Name:        name,
		State:       state,
		Actions:     actions,
		Invertory:   make(map[*Item]bool),
		CurrentRoom: startRoom,
	}
}

func (p Player) HasContainer() bool {
	for item := range p.Invertory {
		if item.IsContainer {
			return true
		}
	}
	return false
}

func (p *Player) PutOnContainer() string {
	if p.ContainerOn {
		return "контейнер уже надет"
	}

	for item := range p.Invertory {
		if item.IsContainer {
			p.ContainerOn = true
			return fmt.Sprintf("вы надели: %s", item.Name)
		}
	}

	for i, item := range p.CurrentRoom.Items {
		if item.IsContainer {
			p.CurrentRoom.Items = append(p.CurrentRoom.Items[:i], p.CurrentRoom.Items[i+1:]...)
			p.Invertory[item] = true
			item.Availability = false
			p.ContainerOn = true
			return fmt.Sprintf("вы надели: %s", item.Name)
		}
	}

	return "нет подходящего контейнера, чтобы надеть"
}

func (p *Player) AddToInvertory(itemName string) string {
	if !p.ContainerOn {
		return "некуда класть"
	}
	if p.CurrentRoom == nil {
		return "непонятно, где вы находитесь"
	}

	var itemToTake *Item
	for i, item := range p.CurrentRoom.Items {
		if item.Name == itemName && item.Availability {
			itemToTake = item
			p.CurrentRoom.Items = append(p.CurrentRoom.Items[:i], p.CurrentRoom.Items[i+1:]...)
			break
		}
	}

	if itemToTake == nil {
		return "нет такого"
	}

	p.Invertory[itemToTake] = true
	itemToTake.Availability = false
	return fmt.Sprintf("предмет добавлен в инвентарь: %s", itemToTake.Name)
}

func (p Player) DescribeCurrentRoomExits() string {
	if p.CurrentRoom == nil {
		return "вы не находитесь ни в одной комнате."
	}

	room := p.CurrentRoom

	if len(room.Exits) == 0 {
		return "выходов нет — вокруг только стены."
	}

	openExits := []string{}
	for name, open := range room.Exits {
		if open {
			openExits = append(openExits, name)
		}
	}
	openExitsStr := strings.Join(openExits, ", ")

	if len(openExits) > 0 {
		return fmt.Sprintf("можно пройти - %v", openExitsStr)
	}

	return "все двери заперты — выхода нет."
}

func (p Player) DescribeCurrentRoom() string {
	if p.CurrentRoom == nil {
		return "вы нигде не находитесь"
	}

	if p.CurrentRoom.Name == RoomRoom {
		if len(p.CurrentRoom.Items) == 0 {
			return fmt.Sprintf("пустая комната. %s", p.DescribeAllRoomExits())
		}

		var tableItems []string
		var chairItems []string

		for _, item := range p.CurrentRoom.Items {
			if item.Name == "рюкзак" {
				chairItems = append(chairItems, item.Name)
			} else {
				tableItems = append(tableItems, item.Name)
			}
		}

		sort.Strings(tableItems)
		sort.Strings(chairItems)

		var parts []string
		if len(tableItems) > 0 {
			parts = append(parts, fmt.Sprintf("на столе: %s", strings.Join(tableItems, ", ")))
		}
		if len(chairItems) > 0 {
			parts = append(parts, fmt.Sprintf("на стуле: %s", strings.Join(chairItems, ", ")))
		}

		return fmt.Sprintf("%s. %s", strings.Join(parts, ", "), p.DescribeAllRoomExits())
	}

	itemNames := make([]string, 0, len(p.CurrentRoom.Items))
	for _, item := range p.CurrentRoom.Items {
		itemNames = append(itemNames, item.Name)
	}
	sort.Strings(itemNames)

	itemsStr := "ничего интересного"
	if len(itemNames) > 0 {
		itemsStr = strings.Join(itemNames, ", ")
	}

	preposition := "в"
	roomName := p.CurrentRoom.Name

	switch p.CurrentRoom.Name {
	case RoomKitchen:
		roomName = "кухне"
		preposition = "на"
	case RoomStreet:
		preposition = "на"
	}

	if p.HasContainer() && p.CurrentRoom.Target != "" {
		p.CurrentRoom.Target = "надо идти в универ"
	}

	targetPart := ""
	if p.CurrentRoom.Target != "" {
		targetPart = fmt.Sprintf(", %s", p.CurrentRoom.Target)
	}

	if events, ok := p.CurrentRoom.Events[CmdLook]; ok && len(events) > 0 {
		if len(p.CurrentRoom.Items) > 0 || p.CurrentRoom.Name == RoomKitchen || p.CurrentRoom.Name == RoomStreet {
			return fmt.Sprintf("ты находишься %s %s, на столе: %s%s. %s",
				preposition, roomName, itemsStr, targetPart, p.DescribeAllRoomExits())
		}
	}

	return fmt.Sprintf("ты находишься %s %s%s. %s", preposition, roomName, targetPart, p.DescribeAllRoomExits())
}

func (p Player) DescribeAllRoomExits() string {
	if p.CurrentRoom == nil {
		return "вы не находитесь ни в одной комнате."
	}

	room := p.CurrentRoom
	if len(room.Exits) == 0 {
		return "в этой комнате нет дверей."
	}

	allExits := make([]string, 0, len(room.Exits))
	for name := range room.Exits {
		allExits = append(allExits, name)
	}

	order := map[string]int{
		RoomKitchen: 1,
		RoomRoom:    2,
		RoomStreet:  3,
	}

	sort.Slice(allExits, func(i, j int) bool {
		oi, oki := order[allExits[i]]
		oj, okj := order[allExits[j]]
		if oki && okj {
			return oi < oj
		}
		if oki {
			return true
		}
		if okj {
			return false
		}
		return allExits[i] < allExits[j]
	})

	return "можно пройти - " + strings.Join(allExits, ", ")
}

func (p *Player) MoveTo(dest string, world *World) string {
	if p.CurrentRoom == nil {
		return "вы не находитесь ни в одной комнате"
	}

	if dest == RoomHome {
		dest = RoomHall
	} else if dest == RoomHall && p.CurrentRoom.Name == RoomStreet {
		dest = RoomHome
	}

	open, exists := p.CurrentRoom.Exits[dest]
	if !exists {
		if p.CurrentRoom.Name == RoomStreet && dest == RoomHome {
			open = true
		} else {
			return fmt.Sprintf("нет пути в %s", dest)
		}
	}
	if !open {
		return "дверь закрыта"
	}

	roomName := dest
	if dest == "домой" {
		roomName = RoomHall
	}

	room, ok := world.Rooms[roomName]
	if !ok {
		return fmt.Sprintf("комната '%s' не найдена", dest)
	}

	p.CurrentRoom = room

	switch room.Name {
	case RoomStreet:
		return "на улице весна. можно пройти - домой"
	case RoomKitchen:
		return fmt.Sprintf("кухня, ничего интересного. %s", p.DescribeAllRoomExits())
	case RoomRoom:
		return fmt.Sprintf("ты в своей комнате. %s", p.DescribeAllRoomExits())
	}

	if len(room.Items) > 0 {
		names := make([]string, len(room.Items))
		for i, item := range room.Items {
			names[i] = item.Name
		}
		sort.Strings(names)
		return fmt.Sprintf("вы перешли в комнату: %s, на столе: %s", dest, strings.Join(names, ", "))
	}

	return fmt.Sprintf("ничего интересного. %s", p.DescribeAllRoomExits())
}

func (p *Player) ApplyItem(command string) string {
	parts := strings.Fields(strings.TrimSpace(strings.ToLower(command)))
	if len(parts) < 3 {
		return "укажите, что и к чему применить, например: 'применить ключ дверь кухня'"
	}

	if !strings.EqualFold(parts[0], CmdUse) {
		return fmt.Sprintf("неизвестная команда: %s", parts[0])
	}

	itemName := parts[1]
	target := parts[2]
	var dest string
	if len(parts) > 3 {
		dest = parts[3]
	}

	var item *Item
	for i := range p.Invertory {
		if strings.Contains(strings.ToLower(i.Name), strings.ToLower(itemName)) {
			item = i
			break
		}
	}

	if item == nil {
		for i, roomItem := range p.CurrentRoom.Items {
			if strings.Contains(strings.ToLower(roomItem.Name), strings.ToLower(itemName)) && roomItem.Availability {
				item = roomItem
				p.CurrentRoom.Items = append(p.CurrentRoom.Items[:i], p.CurrentRoom.Items[i+1:]...)
				p.Invertory[item] = true
				break
			}
		}
	}

	if item == nil {
		return fmt.Sprintf("нет предмета в инвентаре - %s", itemName)
	}

	if strings.Contains(strings.ToLower(target), "двер") {
		if dest == "" {
			for name, open := range p.CurrentRoom.Exits {
				if !open {
					dest = name
					break
				}
			}
			if dest == "" {
				return "все двери уже открыты — применять нечего"
			}
		}

		open, exists := p.CurrentRoom.Exits[dest]
		if !exists {
			return fmt.Sprintf("в этом направлении нет двери: %s", dest)
		}

		if open {
			return fmt.Sprintf("дверь '%s' уже открыта", dest)
		}

		if strings.Contains(strings.ToLower(item.Name), "ключ") {
			p.CurrentRoom.Exits[dest] = true
			return "дверь открыта"
		}

		return fmt.Sprintf("%s не подходит к двери '%s'", item.Name, dest)
	}

	for _, roomItem := range p.CurrentRoom.Items {
		if strings.Contains(strings.ToLower(roomItem.Name), strings.ToLower(target)) {
			return fmt.Sprintf("вы применили %s к %s — (действие пока не реализовано)", item.Name, roomItem.Name)
		}
	}

	return "не к чему применить"
}

func (p Player) ExamineItem(itemName string) string {
	for item := range p.Invertory {
		if itemName == item.Name {
			return item.Description
		}
	}
	return fmt.Sprintf("нет предмета %s в инвертаре", itemName)
}
