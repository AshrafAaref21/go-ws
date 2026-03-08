package models

import (
	"errors"
	"time"

	"github.com/AshrafAaref21/go-ws/internal/db"
	"github.com/AshrafAaref21/go-ws/internal/middlewares"
)

type User struct {
	ID                   int64      `json:"id"`
	Name                 string     `json:"name"`
	Email                string     `json:"email"`
	Password             string     `json:"-"`
	RefreshTokenWeb      *string    `json:"-"`
	RefreshTokenWebAt    *time.Time `json:"-"`
	RefreshTokenMobile   *string    `json:"-"`
	RefreshTokenMobileAt *time.Time `json:"-"`
	CreatedAt            time.Time  `json:"created_at"`
}

func GetUserByEmail(email string) (*User, error) {

	db, err := db.GetDB()
	if err != nil {
		return nil, err
	}

	var user User
	err = db.QueryRow("SELECT id, name, email, password, refresh_token_web, refresh_token_web_at, refresh_token_mobile, refresh_token_mobile_at, created_at FROM users WHERE email = ?", email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.RefreshTokenWeb,
		&user.RefreshTokenWebAt,
		&user.RefreshTokenMobile,
		&user.RefreshTokenMobileAt,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateUserByEmail(email, name, password string) (*User, error) {

	db, err := db.GetDB()
	if err != nil {
		return nil, err
	}

	result, err := db.Exec("INSERT INTO users (email, name, password) VALUES (?, ?, ?)", email, name, password)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	createdAt := time.Now()

	return &User{
		ID:        id,
		Name:      name,
		Email:     email,
		CreatedAt: createdAt,
	}, nil
}

func UpdateUserRefreshToken(userId int64, platform string, refreshToken string) error {

	db, err := db.GetDB()
	if err != nil {
		return err
	}

	now := time.Now()

	switch platform {
	case middlewares.PlatformWeb:
		_, err = db.Exec("UPDATE users SET refresh_token_web = ?, refresh_token_web_at = ? WHERE id = ?", refreshToken, now, userId)
	case middlewares.PlatformMobile:
		_, err = db.Exec("UPDATE users SET refresh_token_mobile = ?, refresh_token_mobile_at = ? WHERE id = ?", refreshToken, now, userId)
	default:
		return errors.New("invalid platform")
	}

	return err
}

func DeleteUserRefreshToken(userId int64, platform string) error {

	db, err := db.GetDB()
	if err != nil {
		return err
	}

	switch platform {
	case middlewares.PlatformWeb:
		_, err = db.Exec("UPDATE users SET refresh_token_web = NULL, refresh_token_web_at = NULL WHERE id = ?", userId)
	case middlewares.PlatformMobile:
		_, err = db.Exec("UPDATE users SET refresh_token_mobile = NULL, refresh_token_mobile_at = NULL WHERE id = ?", userId)
	default:
		return errors.New("invalid platform")
	}

	return err
}

func GetUserByRefreshToken(refreshToken string, platform string) (*User, error) {

	db, err := db.GetDB()
	if err != nil {
		return nil, err
	}

	var user User
	var query string
	switch platform {
	case middlewares.PlatformWeb:
		query = "SELECT id, name, email, password, refresh_token_web, refresh_token_web_at, refresh_token_mobile, refresh_token_mobile_at, created_at FROM users WHERE refresh_token_web = ?"
	case middlewares.PlatformMobile:
		query = "SELECT id, name, email, password, refresh_token_web, refresh_token_web_at, refresh_token_mobile, refresh_token_mobile_at, created_at FROM users WHERE refresh_token_mobile = ?"
	default:
		return nil, errors.New("invalid platform")
	}

	err = db.QueryRow(query, refreshToken).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.RefreshTokenWeb,
		&user.RefreshTokenWebAt,
		&user.RefreshTokenMobile,
		&user.RefreshTokenMobileAt,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
