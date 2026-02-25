package models

import "github.com/Doris-Mwito5/ginja-ai/internal/custom_types"

type Provider struct {
	custom_types.SequentialIdentifier
	Name     string `json:"name"`
	Location string `json:"location"`
	custom_types.Timestamps
}
