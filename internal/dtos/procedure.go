package dtos

type Procedure struct {
	Code        string  `json:"code"`
	Description string  `json:"description"`
	AverageCost float64 `json:"average_cost"`
}