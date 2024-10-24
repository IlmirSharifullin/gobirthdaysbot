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
func (s *Storage) DeleteBirthday(birthdayID int64, userID int64) error {
	const fn = "storage.sqlite.DeleteBirthday"

	stmt, err := s.db.Prepare("DELETE FROM birthdays WHERE id=? AND user_id=?")
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(birthdayID, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", fn, err)
	}

	if count == 0 {
		return storage.ErrNoDeleted
	}

	return nil
}

func (s *Storage) GetUser(ID int64) (*storage.User, error) {
	const fn = "storage.sqlite.GetUser"

	stmt, err := s.db.Prepare("SELECT id, username FROM users WHERE id=?")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer stmt.Close()

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
	defer stmt.Close()

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
	defer stmt.Close()

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

func (s *Storage) GetNextBirthdays(UserID int64) ([]*storage.Birthday, error) {
	const fn = "storage.sqlite.GetNextBirthdays"

	query := `WITH ranked_birthdays AS (SELECT *,
                                 dense_rank() OVER (PARTITION BY user_id ORDER BY
                                     CASE
                                         WHEN strftime('%m-%d', date) < strftime('%m-%d', 'now')
                                             THEN '2' || strftime('%m-%d', date)
                                         ELSE strftime('%m-%d', date)
                                         END) AS drank
                          FROM birthdays)

SELECT id, name, date, additional, user_id
FROM ranked_birthdays
WHERE user_id = ? AND drank = 1;`

	stmt, err := s.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", fn, err)
	}
	defer stmt.Close()

	var birthdays []*storage.Birthday
	rows, err := stmt.Query(UserID)
	for rows.Next() {
		var birthday storage.Birthday
		err := rows.Scan(&birthday.ID, &birthday.Name, &birthday.Date, &birthday.Additional, &birthday.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, storage.ErrBirthdayNotExists
			}
			return nil, fmt.Errorf("%s: %w", fn, err)
		}
		birthdays = append(birthdays, &birthday)
	}

	return birthdays, nil
}

func (s *Storage) GetFilteredBirthdays(nd storage.NotificationDays) ([]*storage.Birthday, error) {
	const fn = "storage.sqlite.GetFilteredBirthdays"

	var birthdays []*storage.Birthday

	var whereClause = make([]string, 0, 4)
	if nd.WeekBefore {
		whereClause = append(whereClause, `next_birthday = DATE('now', '+7 days')`)
	}
	if nd.ThreeDaysBefore {
		whereClause = append(whereClause, `next_birthday = DATE('now', '+3 days')`)
	}
	if nd.DayBefore {
		whereClause = append(whereClause, `next_birthday = DATE('now', '+1 days')`)
	}
	if nd.AtBirthDay {
		whereClause = append(whereClause, `next_birthday = DATE('now')`)
	}
	whereString := strings.Join(whereClause, " OR ")

	query := `
	SELECT id, name, date, additional, user_id,
		CASE 
			WHEN strftime('%m-%d', date) >= strftime('%m-%d')
			    THEN strftime('%Y', 'now') || '-' || strftime('%m-%d', date)
			ELSE strftime('%Y', 'now', '+1 year') || '-' || strftime('%m-%d', date)
		END as next_birthday
	FROM birthdays
	WHERE ` + whereString

	rows, err := s.db.Query(query)

	if err != nil {
		return []*storage.Birthday{}, fmt.Errorf("%s: %w", fn, err)
	}

	var unused interface{}

	for rows.Next() {
		var bd storage.Birthday
		err = rows.Scan(&bd.ID, &bd.Name, &bd.Date, &bd.Additional, &bd.UserID, &unused)
		if err != nil {
			return []*storage.Birthday{}, fmt.Errorf("%s: %w", fn, err)
		}
		birthdays = append(birthdays, &bd)
	}
	return birthdays, nil
}
