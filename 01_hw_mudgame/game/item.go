package main

type Item struct {
	Name         string
	Description  string
	Availability bool
	IsContainer  bool
}

func NewItem(name, description string, isContainer, availability bool) *Item {
	return &Item{
		Name:         name,
		Description:  description,
		IsContainer:  isContainer,
		Availability: availability,
	}
}
