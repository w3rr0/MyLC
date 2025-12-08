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
		SET (verification_token, is_verified)
		VALUES ('', TRUE)
		WHERE email = $1`, email)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		INSERT INTO users
			()
		VALUES
		    ()`)

	return nil
}

func sendVerificationEmail(to string, token string) error {
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := fmt.Sprintf("Verification Link: http://localhost:8080/verify?token=%s", token)
	auth := smtp.PlainAuth("", config.EmailUser, config.EmailPassword, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, config.EmailUser, []string{to}, []byte(message))
}

func CheckAccount(db *sql.DB, email string) (bool, error) {
	var exists bool

	err := db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE email = $1
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
