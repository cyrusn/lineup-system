// Package schedule manage database, the code below is written for sqlite
package schedule

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// Schedule record the information of a schedule in parents day
type Schedule struct {
	ClassCode  string    `json:"classcode"`
	ClassNo    int       `json:"classno"`
	ArrivedAt  time.Time `json:"arrivedAt"`
	Priority   int       `json:"priority"`
	IsNotified bool      `json:"isNotified"`
	IsMeeting  bool      `json:"isMeeting"`
	IsComplete bool      `json:"isComplete"`
}

// DB is *sql.DB that store information of schedule
type DB struct {
	*sql.DB
}

// Query store information for classcode
type Query struct {
	Classcodes []string
	IsComplete string
	Priority   string
}

// SelectedBy find all schedules by query
func (db *DB) SelectedBy(q *Query) ([]*Schedule, error) {
	var schedules []*Schedule

	classcodes := strings.Join(q.Classcodes, "\" or classcode = \"")
	query := fmt.Sprintf(`SELECT * FROM Schedule WHERE (classcode = "%s")`, classcodes)

	if q.IsComplete != "" {
		query += fmt.Sprintf(` and is_complete%s`, q.IsComplete)
	}

	if q.Priority != "" {
		query += fmt.Sprintf(` and priority%s`, q.Priority)
	}

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		s := new(Schedule)
		if err := rows.Scan(
			&s.ClassCode,
			&s.ClassNo,
			&s.ArrivedAt,
			&s.Priority,
			&s.IsNotified,
			&s.IsMeeting,
			&s.IsComplete,
		); err != nil {
			return nil, err
		}

		schedules = append(schedules, s)
	}
	return schedules, nil
}

// Insert Schedule by given classCode and classNo
func (db *DB) Insert(classCode string, classNo int) error {
	_, err := db.Exec(`INSERT INTO Schedule (
      classcode,
      classno,
      arrived_at,
      priority,
      is_notified,
      is_meeting,
      is_complete
    ) values (?, ?, ?, ?, ?, ?, ?)`,
		classCode, classNo, time.Now(), 0, false, false, false,
	)

	return err
}

// Delete delete schedule
func (db *DB) Delete(classcode string, classno int) error {
	_, err := db.Exec(
		`DELETE FROM Schedule WHERE (
      classcode = ? and classno = ?
    )`,
		classcode, classno,
	)

	return err
}

// UpdatePriority update schedule's priority
func (db *DB) UpdatePriority(classcode string, classno int, priority int) error {
	_, err := db.Exec(`UPDATE Schedule SET priority = ? WHERE (
      classcode = ? and classno = ?
    )`,
		priority, classcode, classno,
	)

	return err
}

func (db *DB) toggleFactory(key string) func(string, int) error {
	return func(classCode string, classNo int) error {

		exec := fmt.Sprintf(`
    UPDATE Schedule SET %s = NOT %s WHERE (
      classcode = ? and classno = ?
    )`, key, key)

		_, err := db.Exec(exec,
			classCode, classNo,
		)
		return err

	}
}

// ToggleIsNotified toggle IsNotified
func (db *DB) ToggleIsNotified(classCode string, classNo int) error {
	return db.toggleFactory("is_notified")(classCode, classNo)
}

// ToggleIsMeeting toggle IsMeeting
func (db *DB) ToggleIsMeeting(classCode string, classNo int) error {
	return db.toggleFactory("is_meeting")(classCode, classNo)
}

// ToggleIsComplete toggle IsComplete
func (db *DB) ToggleIsComplete(classCode string, classNo int) error {
	return db.toggleFactory("is_complete")(classCode, classNo)
}
