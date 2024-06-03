package main

import (
	"time"

	"github.com/google/uuid"
)

type CreateUserRequest struct {
	Username     string
	CurrentGame  string
	CurrentLevel int
}

type User struct {
	ID           int
	UUID         uuid.UUID
	Username     string
	CurrentGame  string
	CurrentLevel int
	CreatedAt    time.Time
}

func NewUser(userName, currentGame string, currentLevel int) *User {
	return &User{
		UUID:         uuid.New(),
		Username:     userName,
		CurrentGame:  currentGame,
		CurrentLevel: currentLevel,
		CreatedAt:    time.Now().UTC(),
	}
}
func UpdateUser(userName, currentGame string, currentLevel int, createdAt time.Time) *User {
	return &User{
		UUID:         uuid.New(),
		Username:     userName,
		CurrentGame:  currentGame,
		CurrentLevel: currentLevel,
		CreatedAt:    time.Now().UTC(),
	}
}
