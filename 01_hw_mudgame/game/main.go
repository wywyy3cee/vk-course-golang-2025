package main

import (
	"fmt"
)

func main() {
	initGame()

	for {
		cmd, err := parseCommand()
		if err != nil {
			fmt.Println(err)
			continue
		}

		answer, err := BuildCommand(cmd)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println(answer)
	}

}

func initGame() {
	world, player = initWorld()
}

func initWorld() (*World, *Player) {

	world := NewWorld()
	world.AddRoom(
		"комната",
		map[string][]string{CmdLook: {"биг бубси"}},
		map[string]bool{"коридор": true},
		[]*Item{
			NewItem("рюкзак", "твой старый школьный рюкзак", true, true),
			NewItem("конспекты", "записи по мат. логике", false, true),
			NewItem("ключи", "связка ключей от всех дверей в доме", false, true),
		},
		"",
	)
	world.AddRoom(
		"коридор",
		map[string][]string{CmdLook: {"тут пусто и пыльно"}},
		map[string]bool{"кухня": true, "комната": true, "улица": false},
		nil,
		"",
	)
	world.AddRoom(
		"кухня",
		map[string][]string{CmdLook: {"пахнет кофе"}},
		map[string]bool{"коридор": true},
		[]*Item{NewItem("чай", "тёплый с малиной", false, true)},
		"надо собрать рюкзак и идти в универ",
	)
	world.AddRoom(
		"улица",
		map[string][]string{CmdLook: {"чистый воздух ёпт"}},
		map[string]bool{"коридор": true},
		nil,
		"",
	)

	world.StartRoom = world.Rooms["кухня"]
	player := NewPlayer(
		"wywyy3cee",
		"Активен",
		[]string{CmdLook, CmdGo, CmdPutOn, CmdTake, CmdExamineTheItem},
		world.StartRoom,
	)

	return world, player
}

func handleCommand(command string) string {
	cmd, err := parseCommandLine(command)
	if err != nil {
		return fmt.Sprintf("ошибка: %v", err)
	}

	answer, err := BuildCommand(cmd)
	if err != nil {
		return fmt.Sprintf("ошибка: %v", err)
	}

	return answer
}
