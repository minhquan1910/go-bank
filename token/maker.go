package token

import "time"

// Maker is a interface for managing token
type Maker interface {
	// Create token
	CreateToken(username string, duration time.Duration) (string, *Payload, error)
	// Verify token
	VerifyToken(token string) (*Payload, error)
}
