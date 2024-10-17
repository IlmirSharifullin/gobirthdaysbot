package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"strings"
	"telegram-bot/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const fn = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) InsertUser(ID int64, username string) error {
	const fn = "storage.sqlite.InsertUser"

	stmt, err := s.db.Prepare("INSERT INTO users(id, username) VALUES(?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	defer stmt.Close()
	_, err = stmt.Exec(ID, username)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintPrimaryKey) {
			return storage.ErrUserExists
		}
		return fmt.Errorf("%s: %w", fn, err)
	}

	return nil
}

func (s *Storage) InsertBirthday(birthday *storage.Birthday) error {
	const fn = "storage.sqlite.InsertBirthday"

	stmt, err := s.db.Prepare("INSERT INTO birthdays(name, date, additional, user_id) VALUES(?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(birthday.Name, birthday.Date, birthday.Additional, birthday.UserID)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	birthday.ID = id

	return nil
}

func (s *Storage) GetUser(ID int64) (*storage.User, error) {
	const fn = "storage.sqlite.GetUser"

	stmt, err := s.db.Prepare("SELECT id, username FROM users WHERE id=?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var user storage.User
	err = stmt.QueryRow(ID).Scan(&user.ID, &user.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotExists
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &user, nil
}

func (s *Storage) GetBirthdays(UserID int64) ([]*storage.Birthday, error) {
	const fn = "storage.sqlite.GetBirthdays"

	stmt, err := s.db.Prepare("SELECT id, name, date, additional, user_id FROM birthdays WHERE user_id=?")
	if err != nil {
		return []*storage.Birthday{}, fmt.Errorf("%s: %w", fn, err)
	}

	var birthdays []*storage.Birthday
	rows, err := stmt.Query(UserID)
	for rows.Next() {
		var bd storage.Birthday
		err = rows.Scan(&bd.ID, &bd.Name, &bd.Date, &bd.Additional, &bd.UserID)
		if err != nil {
			return []*storage.Birthday{}, err
		}
		birthdays = append(birthdays, &bd)
	}
	return birthdays, nil
}

func (s *Storage) GetBirthday(ID int64) (*storage.Birthday, error) {
	const fn = "storage.sqlite.GetBirthday"

	stmt, err := s.db.Prepare("SELECT id, name, date, additional, user_id FROM birthdays WHERE id=?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	var birthday storage.Birthday
	err = stmt.QueryRow(ID).Scan(&birthday.ID, &birthday.Name, &birthday.Date, &birthday.Additional, &birthday.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrBirthdayNotExists
		}
		return nil, fmt.Errorf("%s: %w", fn, err)
	}

	return &birthday, nil
}

func (s *Storage) GetFilteredBirthdays(nd storage.NotificationDays) ([]*storage.Birthday, error) {
	const fn = "storage.sqlite.GetBirthdays"

	var birthdays []*storage.Birthday

	var whereClause = make([]string, 0, 4)
	if nd.WeekBefore {
		whereClause = append(whereClause, `date = DATE("now", "-7days")`)
	}
	if nd.ThreeDaysBefore {
		whereClause = append(whereClause, `date = DATE("now", "-3days")`)
	}
	if nd.DayBefore {
		whereClause = append(whereClause, `date = DATE("now", "-1days")`)
	}
	if nd.AtBirthDay {
		whereClause = append(whereClause, `date = DATE("now")`)
	}
	whereString := strings.Join(whereClause, " OR ")

	rows, err := s.db.Query(fmt.Sprintf("SELECT id, name, date, additional, user_id FROM birthdays WHERE %s", whereString))

	if err != nil {
		return []*storage.Birthday{}, fmt.Errorf("%s: %w", fn, err)
	}

	for rows.Next() {
		var bd *storage.Birthday
		err = rows.Scan(&bd.ID, &bd.Name, &bd.Date, &bd.Additional, &bd.UserID)
		if err != nil {
			return []*storage.Birthday{}, err
		}
		birthdays = append(birthdays, bd)
	}
	return birthdays, nil
}
