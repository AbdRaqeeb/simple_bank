package token

import "github.com/google/uuid"

type Payload struct {
	ID uuid.UUID `json:"id,omitempty"`
}
