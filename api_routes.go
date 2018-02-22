package main

import (
	"net/http"
	"time"
	"fmt"
)

func WriteError(w http.ResponseWriter, details string) {
	w.Write([]byte(fmt.Sprintf("{\"type\": \"error\", \"message\": \"%s\"}", details)))
}

func WriteSuccess(w http.ResponseWriter, details string) {
	w.Write([]byte(fmt.Sprintf("{\"type\": \"success\", \"message\": \"%s\"}", details)))
}

func WriteUserRegistrationSuccess(w http.ResponseWriter, userId int) {
	w.Write([]byte(fmt.Sprintf("{\"type\": \"success\", \"userId\": \"%d\"}", userId)))
}

func WriteUserInformationSuccess(w http.ResponseWriter, userInfo string) {
	w.Write([]byte(fmt.Sprintf(
		"{\"type\": \"success\", \"userInfo\": \"%s\"}",
		userInfo,
	)))
}

func ApiCreateUserRoute(w http.ResponseWriter, r *http.Request) {
	Logger("v1/Api/Create/User", r, time.Now())
	w.Header().Set("Content-Type", "application/json")
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	if len(username) < 5 {
		WriteError(w, "Username is too short")
		return
	}
	if len(password) < 5 {
		WriteError(w, "Password is too short")
		return
	}
	Logger(fmt.Sprintf("Registration API: username=%s,email=%s", username, email), r, time.Now())
	userId, err := authCreate(username, password, email)
	checkErr(err)
	WriteUserRegistrationSuccess(w, userId)
}

func ApiGetUserRoute(w http.ResponseWriter, r *http.Request) {
	Logger("v1/Api/User/Info", r, time.Now())
	w.Header().Set("Content-Type", "application/json")
	userId := r.FormValue("userId")
	if len(userId) <= 0 {
		WriteError(w, "User ID is not valid")
		return
	}
	Logger(fmt.Sprintf("Get User API: userId=%d", userId), r, time.Now())
	userInfo, err := authGetUserInfo(userId)
	checkErr(err)
	WriteUserInformationSuccess(w, userInfo)
}

func ApiAuthenticateUserRoute(w http.ResponseWriter, r *http.Request) {
	Logger("v1/Api/User/Login", r, time.Now())
	w.Header().Set("Content-Type", "application/json")
	username := r.FormValue("username")
	password := r.FormValue("password")
	if len(username) < 5 {
		WriteError(w, "Username is invalid and too short")
		return
	}
	if len(password) < 5 {
		WriteError(w, "Password is invalid and too short")
		return
	}
	Logger(fmt.Sprintf("Authenticate User API: username=%s", username), r, time.Now())
	userId, token, err := authLogin(username, password)
	checkErr(err)
	WriteSuccess(w, fmt.Sprintf("userId=%d,token=%s", userId, token)) // generic send back
}

func ApiAuthenticationChangePassword(w http.ResponseWriter, r *http.Request) {
	Logger("v1/Api/User/Changepassword", r, time.Now())
	userId := r.FormValue("userId")
	token := r.FormValue("token")
	newPassword := r.FormValue("newPassword")
	if authValidateToken(userId, token) {
		authForceChangePassword(userId, newPassword)
		WriteSuccess(w, "password changed")
	} else {
		WriteError(w, "invalid token/id combo")
	}
}