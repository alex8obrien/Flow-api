package domain

import (
	"errors"
	"uuid"
)

type User struct {
	ID    uuid.UUID
	Name  string
	Email string
}

var ErrInvalidEmail = errors.New("invalid email")
