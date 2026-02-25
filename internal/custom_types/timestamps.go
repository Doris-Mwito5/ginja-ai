package custom_types

import (
	"time"
)

type Timestamps struct {
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (t *Timestamps) Touch() {

	t.UpdatedAt = time.Now()

	if t.CreatedAt.IsZero() {
		t.CreatedAt = t.UpdatedAt
	}
}
