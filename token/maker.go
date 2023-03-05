package token

import "time"

// interface for managing tokens
type Maker interface {
	// creates a new token for a specific username and duration
	CreateToken(username string, duration time.Duration) (string, error)
	// checks if the token is valid or not, if valid - method returns the payload stored in the token
	VerifyToken(token string) (*Payload, error)
}
