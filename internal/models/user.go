package models

import (
	"github.com/Doris-Mwito5/ginja-ai/internal/custom_types"
)

type User struct {
	custom_types.SequentialIdentifier
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"`
	IsActive     bool   `json:"is_active"`
	custom_types.Timestamps
}