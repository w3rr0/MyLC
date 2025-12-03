package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

func CreateEvent(db *sql.DB, start time.Time, end time.Time) error {
	if end.Before(start.Add(15 * time.Minute)) {
		return fmt.Errorf("end time must be greater than start time")
	}

	var newId int

	err := db.QueryRow(
		`INSERT INTO event_manager DEFAULT VALUES RETURNING id`,
	).Scan(&newId)
	if err != nil {
		return err
	}

	newName := tableFrom(newId)
	columns := CreateColumnsFromTime(start, end)
	var colDefs []string
	for _, colName := range columns {
		colDefs = append(colDefs, fmt.Sprintf("%q INT[]", colName))
	}
	colsString := strings.Join(colDefs, ", ")

	if len(colDefs) == 0 {
		return fmt.Errorf("cannot create table with no columns")
	}

	query := fmt.Sprintf("CREATE TABLE %q (%s)", newName, colsString)

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func DeleteEvent(db *sql.DB, eventId int) error {
	err := CheckTable(db, eventId)
	if err != nil {
		return err
	}

	toDelete := tableFrom(eventId)

	query1 := fmt.Sprintf("DROP TABLE %q", toDelete)
	query2 := fmt.Sprintf("DELETE FROM event_manager WHERE id = %d", eventId)

	_, err = db.Exec(query1)
	if err != nil {
		return err
	}
	_, err = db.Exec(query2)
	if err != nil {
		return err
	}

	return nil
}

func ChangeAvailability(db *sql.DB, eventId int, userId int, availability map[string]string) error {
	err := CheckUser(db, userId)
	if err != nil {
		return err
	}

	colNames := getMapKeys(availability)
	err = CheckTableColumns(db, eventId, colNames)
	if err != nil {
		return err
	}

	colValues := getMapValues(availability)
	err = CheckValues(colValues)
	if err != nil {
		return err
	}

	for timeSlot, attendance := range availability {
		err := SetAvailabilityFor(db, tableFrom(eventId), timeSlot, attendance, userId)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetAllCurrentEvents(db *sql.DB, userId int) ([]models.EventSummary, error) {
	var eventsList []models.EventSummary

	getEventFromManagerQuery := `
		SELECT id, name
		FROM event_manager`

	maxFilledQuery := `
		SELECT COUNT(*)
		FROM information_schema.columns
		WHERE table_name = $1`

	rows, err := db.Query(getEventFromManagerQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var event models.EventSummary
		if err := rows.Scan(&event.ID, &event.Name); err != nil {
			return nil, err
		}
		tableName := tableFrom(event.ID)
		tableColumns, err := GetTableColumns(db, tableName)
		if err != nil {
			return nil, err
		}

		var timeSlots []string
		for _, col := range tableColumns {
			if col != "name" {
				timeSlots = append(timeSlots, col)
			}
		}

		var parts []string
		for _, slot := range timeSlots {
			parts = append(parts, fmt.Sprintf("SELECT unnest(%q) AS val FROM %q", slot, tableName))
		}

		queryBody := strings.Join(parts, " UNION ALL ")

		query := fmt.Sprintf(`
			SELECT COUNT(*)
			FROM (%s) AS combined
			WHERE val = $1
			`, queryBody)

		var filledCount int
		err = db.QueryRow(query, userId).Scan(&filledCount)
		if err != nil {
			return nil, err
		}

		var maxFilled int
		err = db.QueryRow(maxFilledQuery, tableName).Scan(&maxFilled)
		if err != nil {
			return nil, err
		}

		filled := filledCount == maxFilled
		event.Filled = filled

		eventsList = append(eventsList, event)
	}

	return eventsList, nil
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

func tableFrom(eventId int) string {
	return fmt.Sprintf("table_%d", eventId)
}

func IfTableExist(db *sql.DB, id int) (bool, error) {
	tableName := tableFrom(id)

	var exist bool
	query1 := `
		SELECT EXISTS (
			SELECT 1
			FROM pg_tables
			WHERE tablename = $1
		)
	`
	err := db.QueryRow(query1, tableName).Scan(&exist)
	if err != nil {
		return false, err
	}

	var inManager bool
	query2 := `
		SELECT EXISTS (
			SELECT 1
			FROM event_manager
			WHERE id = $1
		)
	`
	err = db.QueryRow(query2, id).Scan(&inManager)
	if err != nil {
		return false, err
	}

	if exist != inManager {
		if exist {
			log.Println("Table exist, but no corresponding entry in event_manager")
		}
		if !inManager {
			log.Println("Table doesn't exist, but entry exist in event_manager")
		}
	}

	if !exist {
		return false, nil
	}
	if !inManager {
		return false, nil
	}
	return true, nil
}

func IfUserExists(db *sql.DB, userId int) (bool, error) {
	var exist bool

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE id = $1
		)
	`
	err := db.QueryRow(query, userId).Scan(&exist)
	if err != nil {
		return false, err
	}

	if !exist {
		return false, nil
	}
	return true, nil
}

func CheckTable(db *sql.DB, eventId int) error {
	check, err := IfTableExist(db, eventId)
	if err != nil {
		return err
	}
	if !check {
		return errors.New("no corresponding event in database")
	}
	return nil
}

func CheckUser(db *sql.DB, userId int) error {
	check, err := IfUserExists(db, userId)
	if err != nil {
		return err
	}
	if !check {
		return errors.New("no corresponding user in database")
	}
	return nil
}

func CheckTableColumns(db *sql.DB, eventId int, colNames []string) error {
	err := CheckTable(db, eventId)
	if err != nil {
		return err
	}

	tableName := tableFrom(eventId)
	tableColumns, err := GetTableColumns(db, tableName)
	if err != nil {
		return err
	}

	if !isSubset(colNames, tableColumns) {
		return errors.New("column names not in table columns")
	}

	return nil
}

func CheckValues(colValues []string) error {
	allowed := []string{"yes", "maybe", "no"}
	if !isSubset(colValues, allowed) {
		return errors.New("not allowed values")
	}
	return nil
}

func GetTableColumns(db *sql.DB, tableName string) ([]string, error) {
	query := `
		SELECT column_name
		FROM information_schema.columns
		WHERE table_name = $1
	`
	rows, err := db.Query(query, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []string
	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {
			return nil, err
		}
		cols = append(cols, col)
	}
	return cols, nil
}

func isSubset(smaller []string, bigger []string) bool {
	set := make(map[string]struct{}, len(bigger))
	for _, s := range bigger {
		set[s] = struct{}{}
	}

	for _, s := range smaller {
		if _, ok := set[s]; !ok {
			return false
		}
	}
	return true
}

func getMapKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func getMapValues(m map[string]string) []string {
	values := make([]string, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}

func SetAvailabilityFor(db *sql.DB, tableName string, timeSlot string, attendance string, userId int) error {
	addQuery := fmt.Sprintf(`
		UPDATE %q
		SET %q = array_append(%q, $1)
		WHERE name = $2
		AND NOT (%q @> ARRAY[$1::integer])`, tableName, timeSlot, timeSlot, timeSlot)
	res, err := db.Exec(addQuery, userId, attendance)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	switch rowsAffected {
	case 0:
		log.Println("User already present in the availability list")
	case 1:
	default:
		return errors.New("multiple rows affected during availability update for one user")
	}

	opposites, err := oppositeAttendance(attendance)
	if err != nil {
		return err
	}

	removeQuery := fmt.Sprintf(`
		UPDATE %q
		SET %q = array_remove(%q, $1)
		WHERE name IN ($2, $3)
		AND %q @> ARRAY[$1::integer]`, tableName, timeSlot, timeSlot, timeSlot)
	res, err = db.Exec(removeQuery, userId, opposites[0], opposites[1])
	if err != nil {
		return err
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}

func oppositeAttendance(attendance string) ([2]string, error) {
	var opposites [2]string

	switch attendance {
	case "yes":
		opposites = [2]string{"no", "maybe"}
	case "no":
		opposites = [2]string{"yes", "maybe"}
	case "maybe":
		opposites = [2]string{"yes", "no"}
	default:
		return [2]string{}, errors.New("attendance must be 'yes', 'maybe' or 'no'")
	}

	return opposites, nil
}
