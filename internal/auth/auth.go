package auth

import (
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int
	Username string
	Password string
	Role     string
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, plainTextPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainTextPassword)) == nil
}

func CreateUser(db *sql.DB, username, password, role string) (*User, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	result, err := db.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3)", username, hashedPassword, role)
	if err != nil {
		return nil, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &User{ID: int(userID), Username: username, Role: role}, nil
}

func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, username, password, role FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
