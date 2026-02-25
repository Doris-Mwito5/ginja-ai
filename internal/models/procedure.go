package models

import "github.com/Doris-Mwito5/ginja-ai/internal/custom_types"

type Procedure struct {
	custom_types.SequentialIdentifier
	Code        string  `json:"code"`
	Description string  `json:"description"`
	AverageCost float64 `json:"average_cost"`
	custom_types.Timestamps
}
