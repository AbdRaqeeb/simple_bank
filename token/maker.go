package token

import "time"

// Maker is an interface to manage tokens
type Maker interface {
	// CreateToken generates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration)

	// VerifyToken checks for token validity
	VerifyToken(token string) (*Payload, error)
}
