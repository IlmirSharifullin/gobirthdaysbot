package storage

import (
	"errors"
	"time"
)

var (
	ErrUserExists        = errors.New("user with this ID already exists")
	ErrUserNotExists     = errors.New("user with this ID not exists")
	ErrBirthdayNotExists = errors.New("birthday with this ID not exists")
)

type Storage interface {
	GetUser(ID int64) (*User, error)
	InsertUser(ID int64, username string) error

	GetBirthday(ID int64) (*Birthday, error)
	GetBirthdays(UserID int64) ([]*Birthday, error)
	GetNextBirthdays(UserID int64) ([]*Birthday, error)
	GetFilteredBirthdays(nd NotificationDays) ([]*Birthday, error)
	InsertBirthday(birthday *Birthday) error
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

type NotificationDays struct {
	WeekBefore      bool
	ThreeDaysBefore bool
	DayBefore       bool
	AtBirthDay      bool
}
