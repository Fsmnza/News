// models.go

package models

import (
	"errors"
	"time"
)

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
	ErrUserNotFound       = errors.New("models: user not found")
)

const (
	RoleUser    = "user"
	RoleAdmin   = "admin"
	RoleTeacher = "teacher"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
	Role           string
}

type News struct {
	ID       int
	Title    string
	Content  string
	Date     time.Time
	Category string
}
