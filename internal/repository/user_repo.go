package repository

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"go_server/internal/config"
	"go_server/internal/models"
	"net/smtp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func GetAllUsers(db *sql.DB) ([]models.User, error) {
	rows, err := db.Query(`SELECT id, first_name, last_name, email, "group" FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []models.User{}

	for rows.Next() {
		var u models.User
		var first, last, group string

		if err := rows.Scan(
			&u.ID,
			&first,
			&last,
			&u.Email,
			&group,
		); err != nil {
			return nil, err
		}

		u.Name.First = first
		u.Name.Last = last
		u.Group = models.Group(group)

		users = append(users, u)
	}
	return users, nil
}

func CreateUser(db *sql.DB, email string, password string) error {
	exists, err := CheckAccount(db, email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	token, err := GenerateToken(32)
	if err != nil {
		return err
	}

	newAccount := models.Account{
		Email:             email,
		PasswordHash:      string(passwordHash),
		VerificationToken: token,
	}

	err = AppendAccount(db, newAccount)
	if err != nil {
		return err
	}

	err = sendVerificationEmail(email, token)
	if err != nil {
		return err
	}

	return nil
}

func VerifyUser(db *sql.DB, token string) error {
	var email string

	row := db.QueryRow(`
		SELECT email
		FROM accounts
		WHERE verification_token = $1`, token)
	err := row.Scan(&email)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE accounts
		SET verification_token='', is_verified=TRUE
		WHERE email = $1`, email)
	if err != nil {
		return err
	}

	first, last, err := extractNameFromEmail(email)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO users
			(first_name, last_name, email, "group")
		VALUES
		    ($1, $2, $3, 'IT')
	`, first, last, email)

	return nil
}

func extractNameFromEmail(email string) (string, string, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", "", errors.New("invalid email format")
	}

	domain := parts[1]
	if !strings.HasPrefix(domain, "iaeste.") {
		return "", "", errors.New("not an iaeste email")
	}

	nameParts := strings.Split(parts[0], ".")
	if len(nameParts) != 2 {
		return "", "", errors.New("expected format firstname.lastname")
	}

	first := capitalize(nameParts[0])
	last := capitalize(nameParts[1])

	return first, last, nil
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

func sendVerificationEmail(to string, token string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	subject := "Weryfikacja konta MyLC"
	body := fmt.Sprintf("Kliknij w link, aby zweryfikować konto: http://localhost:8080/verify_user?token=%s", token)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to, subject, body))

	auth := smtp.PlainAuth("", config.EmailUser, config.EmailPassword, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, config.EmailUser, []string{to}, msg)
}

func CheckAccount(db *sql.DB, email string) (bool, error) {
	var exists bool

	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE email = $1
		)
	`, email).Scan(&exists)

	return exists, err
}

func GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func AppendAccount(db *sql.DB, account models.Account) error {
	_, err := db.Exec(`
		INSERT INTO accounts
			(email, password_hash, verification_token)
		VALUES
		    ($1, $2, $3)
	`, account.Email, account.PasswordHash, account.VerificationToken)
	if err != nil {
		return err
	}

	return nil
}
