package models

type ClaimList struct {
	Claims []*Claim `json:"claims"`
	Pagination *Pagination `json:"pagination"`
}
