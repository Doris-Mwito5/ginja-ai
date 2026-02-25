package models

import (
	"github.com/Doris-Mwito5/ginja-ai/internal/custom_types"
)

type Claim struct {
	custom_types.SequentialIdentifier
	MemberID        int64                   `json:"member_id"`
	ProviderID      int64                   `json:"provider_id"`
	ProcedureCode   string                   `json:"procedure_code"`
	DiagnosisCode   string                   `json:"diagnosis_code"`
	RequestedAmount float64                  `json:"requested_amount"`
	ApprovedAmount  float64                  `json:"approved_amount"`
	Status          custom_types.ClaimStatus `json:"status"`
	FraudFlag       bool                     `json:"fraud_flag"`
	RejectionReason string                   `json:"rejection_reason"`
	custom_types.Timestamps
}
