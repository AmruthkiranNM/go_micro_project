package services

import (
	"errors"
	"inventory-system/database"
	"inventory-system/models"

	"golang.org/x/crypto/bcrypt"
)

// AuthenticateUser verifies username+password and returns the user on success.
func AuthenticateUser(username, password string) (*models.User, error) {
	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, username, password_hash, email FROM users WHERE username = ?", username,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid username or password")
	}
	return &user, nil
}

// GetUserByID fetches a user record by primary key.
func GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := database.DB.QueryRow(
		"SELECT id, username, email FROM users WHERE id = ?", id,
	).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateProfile updates username, email and optionally password_hash for a user.
func UpdateProfile(id int, username, email, newPassword string) error {
	if username == "" || email == "" {
		return errors.New("username and email are required")
	}
	if newPassword != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		_, err = database.DB.Exec(
			"UPDATE users SET username=?, email=?, password_hash=? WHERE id=?",
			username, email, string(hash), id,
		)
		return err
	}
	_, err := database.DB.Exec(
		"UPDATE users SET username=?, email=? WHERE id=?",
		username, email, id,
	)
	return err
}
