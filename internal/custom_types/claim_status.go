package custom_types

import "database/sql/driver"

type ClaimStatus string

const (
	ClaimStatusApproved ClaimStatus = "approved"
	ClaimStatusPartial  ClaimStatus = "partial"
	ClaimStatusRejected ClaimStatus = "rejected"
)

func (c *ClaimStatus) Scan(value interface{}) error {
	*c = ClaimStatus(string(value.([]uint8)))
	return nil
}

func (c *ClaimStatus) Value() (driver.Value, error) {
	return c.String(), nil
}

func (c ClaimStatus) String() string {
	return string(c)
}
