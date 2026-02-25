package dtos

type Member struct {
	FullName     string  `json:"full_name"`
	IsActive     bool    `json:"is_active"`
	BenefitLimit float64 `json:"benefit_limit"`
	UsedAmount   float64 `json:"used_amount"`
}