package main

import (
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Username     string
	CurrentGame  string
	CurrentLevel int
}

func NewUser(userName, currentGame string, currentLevel int) *User {
	return &User{
		ID:           uuid.New(),
		Username:     userName,
		CurrentGame:  currentGame,
		CurrentLevel: currentLevel,
	}
}
