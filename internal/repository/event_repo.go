package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func CreateEvent(db *sql.DB, start time.Time, end time.Time) error {
	if end.Before(start.Add(15 * time.Minute)) {
		return fmt.Errorf("end time must be greater than start time")
	}

	var new_id int

	err := db.QueryRow(
		`INSERT INTO event_manager DEFAULT VALUES RETURNING id`,
	).Scan(&new_id)
	if err != nil {
		return err
	}

	new_name := fmt.Sprintf("table_%d", new_id)
	columns := CreateColumnsFromTime(start, end)
	var colDefs []string
	for _, colName := range columns {
		colDefs = append(colDefs, fmt.Sprintf("%q INT[]", colName))
	}
	colsString := strings.Join(colDefs, ", ")

	if len(colDefs) == 0 {
		return fmt.Errorf("cannot create table with no columns")
	}

	query := fmt.Sprintf("CREATE TABLE %q (%s)", new_name, colsString)

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func CreateColumnsFromTime(start time.Time, end time.Time) []string {
	roundedStart := RoundToHalfHour(start)
	roundedEnd := RoundToHalfHour(end)

	return GenerateTimeSlots(roundedStart, roundedEnd)
}

func RoundToHalfHour(t time.Time) time.Time {
	minutes := t.Minute()
	rounded := (minutes + 15) / 30 * 30

	roundedTime := time.Date(
		t.Year(), t.Month(), t.Day(),
		t.Hour(), rounded, 0, 0, t.Location(),
	)

	return roundedTime
}

func GenerateTimeSlots(start, stop time.Time) []string {
	var slots []string

	for current := start; current.Before(stop); current = current.Add(30 * time.Minute) {

		next := current.Add(30 * time.Minute)

		slot := fmt.Sprintf("%s-%s", current.Format("15:04"), next.Format("15:04"))

		slots = append(slots, slot)
	}

	return slots
}
