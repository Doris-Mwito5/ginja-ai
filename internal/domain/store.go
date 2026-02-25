package domain

type Store struct {
	ClaimDomain     ClaimDomain
	MemberDomain    MemberDomain
	ProcedureDomain ProcedureDomain
	ProviderDomain  ProviderDomain
	UserDomain      UserDomain
}

func NewStore() *Store {
	return &Store{
		ClaimDomain:     NewClaimDomain(),
		MemberDomain:    NewMemberDomain(),
		ProcedureDomain: NewProcedureDomain(),
		ProviderDomain:  NewProviderDomain(),
		UserDomain:      NewUserDomain(),
	}
}
