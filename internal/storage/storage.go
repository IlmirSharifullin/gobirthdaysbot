package storage

import (
	"errors"
	"time"
)

var (
	ErrUserExists = errors.New("User with this ID already exists")
	ErrNotExists  = errors.New("User with this ID not exists")
)

type Storage interface {
	GetUser(ID int64) (*User, error)
	InsertUser(ID int64, username string) error

	GetBirthday(ID int64) (*Birthday, error)
	GetBirthdays(UserID int64) ([]*Birthday, error)
	InsertBirthday(ID int64, name string, date time.Time, additional string, userId int64) error
}

type User struct {
	ID        int64
	Username  string
	Birthdays []Birthday
}

type Birthday struct {
	ID         int64
	Name       string
	Date       time.Time
	Additional string
	UserID     int64
}
