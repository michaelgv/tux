package main

import "fmt"

func UserCreate(username string, password string, email string) (int, error) {
	db := MakeDatabase()
	var userCount int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&userCount)
	if userCount > 0 {
		return 0, fmt.Errorf("username '%s' is already in use", username)
	}

	if email != "" {
		db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&userCount)
		if userCount > 0 {
			return 0, fmt.Errorf("email '%s' is already in use", email)
		}
	}

	result := db.Exec("INSERT INTO users (username, password, email, token) VALUES (?, ?, ?, ?)", username, authMakePassword(password), email, authMakePassword("temporary_string"))
	return result.LastInsertId(), nil
}

func UserValidateToken(username string, token string) (bool, error) {
	db := MakeDatabase()
	var userCount int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&userCount)
	if userCount > 0 {
		return false, fmt.Errorf("username '%s' does not exist", username)
	}
	result := db.Exec("SELECT token FROM users WHERE username = ? AND token = ?", username, token)
	if result.RowsAffected() < 0 {
		return false, fmt.Errorf("invalid token '%s'", token)
	} else {
		return true, nil
	}
}

func UpdateUserSessionToken(username string, token string) (string, error) {
	db := MakeDatabase()
	var userCount int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&userCount)
	if userCount > 0 {
		return "", fmt.Errorf("username '%s' does not exist", username)
	}
	db.Exec("UPDATE users SET token = ? WHERE username = ?", token, username)
	return token, nil
}