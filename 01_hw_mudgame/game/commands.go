package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func parseCommandLine(input string) ([]string, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("нет ввода")
	}

	parts := strings.Fields(input)

	if len(parts) == 0 {
		return nil, fmt.Errorf("недостаточно аргументов")
	}
	if len(parts) > 4 {
		return nil, fmt.Errorf("слишком много аргументов")
	}

	return parts, nil
}

func parseCommand() ([]string, error) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")

	if !scanner.Scan() {
		return nil, fmt.Errorf("нет ввода")
	}

	input := strings.TrimSpace(scanner.Text())
	if input == "" {
		return nil, fmt.Errorf("пустая команда")
	}

	parts := strings.Fields(input)

	if len(parts) > 4 {
		return nil, fmt.Errorf("слишком много аргументов")
	}

	return parts, nil
}

func BuildCommand(parts []string) (string, error) {
	if len(parts) == 0 {
		return "", fmt.Errorf("команда не указана")
	}

	switch parts[0] {
	case CmdGo:
		return handleGo(parts)
	case CmdTake:
		return handleTake(parts)
	case CmdLook:
		return handleLook(parts)
	case CmdUse:
		return handleUse(parts)
	case CmdPutOn:
		return handlePutOn(parts)
	case CmdExamineTheItem:
		return handleExamine(parts)
	case CmdExit:
		return handleExit(parts)
	case CmdExits:
		return handleExits(parts)
	default:
		return "неизвестная команда", nil
	}
}

func handleGo(parts []string) (string, error) {
	if len(parts) < 2 {
		return "", fmt.Errorf("не указан путь")
	}
	if len(parts) > 2 {
		return "", fmt.Errorf("слишком много аргументов")
	}
	return player.MoveTo(parts[1], world), nil
}

func handleTake(parts []string) (string, error) {
	if len(parts) < 2 {
		return "", fmt.Errorf("не указан предмет")
	}
	if len(parts) > 2 {
		return "", fmt.Errorf("слишком много аргументов")
	}
	return player.AddToInvertory(parts[1]), nil
}

func handleLook(parts []string) (string, error) {
	if len(parts) > 1 {
		return "", fmt.Errorf("слишком много аргументов")
	}
	return player.DescribeCurrentRoom(), nil
}

func handleUse(parts []string) (string, error) {
	if len(parts) < 3 {
		return "", fmt.Errorf("мало аргументов")
	}
	if len(parts) > 4 {
		return "", fmt.Errorf("слишком много аргументов")
	}
	str := strings.Join(parts, " ")
	return player.ApplyItem(str), nil
}

func handlePutOn(parts []string) (string, error) {
	if len(parts) < 2 {
		return "", fmt.Errorf("не указан предмет")
	}
	if len(parts) > 2 {
		return "", fmt.Errorf("слишком много аргументов")
	}
	return player.PutOnContainer(), nil
}

func handleExamine(parts []string) (string, error) {
	if len(parts) < 2 {
		return "", fmt.Errorf("не указан предмет")
	}
	if len(parts) > 2 {
		return "", fmt.Errorf("слишком много аргументов")
	}
	return player.ExamineItem(parts[1]), nil
}

func handleExit(parts []string) (string, error) {
	if len(parts) > 1 {
		return "", fmt.Errorf("слишком много аргументов")
	}
	fmt.Println("выход из игры")
	os.Exit(0)
	return "", nil
}

func handleExits(parts []string) (string, error) {
	if len(parts) > 1 {
		return "", fmt.Errorf("слишком много аргументов")
	}
	return player.DescribeAllRoomExits(), nil
}
