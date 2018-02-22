package main

import "golang.org/x/crypto/pbkdf2"

import (
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"strings"
	"crypto/subtle"
	"log"
	"time"
)

/**
 * Make a secure password
 */
func authMakePassword(password string) string {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	checkErr(err)
	passwordHash := pbkdf2.Key([]byte(password), salt, 8192, 64, sha512.New)
	return hex.EncodeToString(salt) + ":" + hex.EncodeToString(passwordHash)
}

func authCreate(username string, password string, email string) (int, error) {
	userId, err := UserCreate(username, password, email)
	if err != nil {
		return 0, err
	}
	return userId, nil
}

func authCheckPassword(password string, actualPasswordCombined string) bool {
	passwordParts := strings.Split(actualPasswordCombined, ":")
	salt, _ := hex.DecodeString(passwordParts[0])
	actualPasswordHash, _ := hex.DecodeString(passwordParts[1])
	providedPasswordHash := pbkdf2.Key([]byte(password), salt, 8192, 64, sha512.New)
	return subtle.ConstantTimeCompare(actualPasswordHash, providedPasswordHash) == 1
}

func authLogin(username string, password string) (int, error) {
	db := MakeDatabase()
	if len(password) > 255 {
		return 0, fmt.Errorf("password is too long, max length is %s", "255")
	}

	rows := db.Query("SELECT id, password FROM users WHERE username = ?", username)
	if !rows.Next() {
		log.Printf("Authentication failure on user=%s: bad username (%s)", username, username)
		return 0, fmt.Errorf("invalid_username=%s", username)
	}
	var userId int
	var actualPasswordCombined string
	rows.Scan(&userId, &actualPasswordCombined)
	rows.Close()

	if authCheckPassword(password, actualPasswordCombined) {
		log.Printf("Authentication successful for user=%s", username)
		return userId, nil
	} else {
		log.Printf("Authentication failure on user=%s: bad username (%s)", username, username)
		return 0, fmt.Errorf("invalid_username=%s", username)
	}
}

func authGetUserInfo(userId string) (string, error) {
	db := MakeDatabase()

	rows := db.Query("SELECT username, email FROM users WHERE id = ?", userId);
	if !rows.Next() {
		GenLogger(fmt.Sprintf("[Auth_Service::GetUserInfo] User %d is not found", userId), time.Now())
	}
	var username string
	var email string
	rows.Scan(&username, &email)
	rows.Close()

	return fmt.Sprintf("id=%s,user=%s,email=%s", userId, username, email), nil
}

func authChangePassword(userId int, oldPassword string, newPassword string) error {
	db := MakeDatabase()
	rows := db.Query("SELECT password FROM users WHERE id = ?", userId)
	if !rows.Next() {
		log.Printf("Error changing password: bad user ID (%d)", userId)
		return fmt.Errorf("invalid_user_id=%d", userId)
	}
	var actualPasswordCombined string
	rows.Scan(&actualPasswordCombined)
	rows.Close()

	if authCheckPassword(oldPassword, actualPasswordCombined) {
		db.Exec("UPDATE users SET password = ? WHERE id = ?", authMakePassword(newPassword), userId)
		log.Printf("Successful password change for user_id=%d", userId)
		return nil
	} else {
		log.Printf("Change password authentication failure for user_id=%d", userId)
		return fmt.Errorf("change_password_failure_for_user_id=%d", userId)
	}
}

func authForceChangePassword(userId int, password string) {
	db := MakeDatabase()
	db.Exec("UPDATE users SET password = ? WHERE id = ?", authMakePassword(password), userId)
}