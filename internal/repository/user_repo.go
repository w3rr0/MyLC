package repository

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"go_server/internal/models"

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

	return nil
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
