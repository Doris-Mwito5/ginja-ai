package models

import "github.com/Doris-Mwito5/ginja-ai/internal/custom_types"

type Member struct {
	custom_types.SequentialIdentifier
	FullName     string  `json:"full_name"`
	IsActive     bool    `json:"is_active"`
	BenefitLimit float64 `json:"benefit_limit"`
	UsedAmount   float64 `json:"used_amount"`
	custom_types.Timestamps
}
