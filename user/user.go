package user

import "errors"

var ErrUserNotFound = errors.New("user is not found")

type User struct {
	ID       string
	Username string
	Email    string
}
