package dtos

type ClaimSubmissionForm struct {
	MemberID        int64  `json:"member_id"         binding:"required"`
	ProviderID      int64  `json:"provider_id"       binding:"required"`
	ProcedureCode   string  `json:"procedure_code"    binding:"required"`
	DiagnosisCode   string  `json:"diagnosis_code"    binding:"required"`
	RequestedAmount float64 `json:"requested_amount"  binding:"required,gt=0"`
}

type ClaimSubmissionResponse struct {
	ClaimID         int64   `json:"claim_id"`
	Status          string  `json:"status"`
	FraudFlag       bool    `json:"fraud_flag"`
	ApprovedAmount  float64 `json:"approved_amount"`
	RejectionReason string  `json:"rejection_reason,omitempty"`
}