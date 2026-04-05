package domain

import (
	"errors"
	"net/mail"
)

// User represents a platform user who can initiate charging sessions.
type User struct {
	ID    string
	Name  string
	Email string
}

// Validate checks the user fields for consistency.
func (u *User) Validate() error {
	if u.ID == "" {
		return errors.New("user id is required")
	}
	if u.Name == "" {
		return errors.New("user name is required")
	}
	if u.Email == "" {
		return errors.New("user email is required")
	}
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("user email is invalid")
	}
	return nil
}
