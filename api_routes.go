package main

import (
	"net/http"
	"time"
	"fmt"
	"log"
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

func ApiAuthenticationChangePasswordRoute(w http.ResponseWriter, r *http.Request) {
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

func makeUserJsonString(userId string, username string, email string, status string, admin string) string {
	return fmt.Sprintf("{\"type\": \"user_record\", \"userId\": \"%s\", \"username\": \"%s\", \"email\": \"%s\", \"status\": \"%s\", \"admin\": \"%s\"}", userId, username, email, status, admin)
}

func ApiInternalFlushRedisCache(w http.ResponseWriter, r *http.Request) {
	FlushRedisDB()
	log.Println("Flushed Redis Cache")
	w.Write([]byte("{}"))
}

func ApiListUsersSafelyRoute(w http.ResponseWriter, r *http.Request) {
	Logger("v1/Api/User/List", r, time.Now())
	start := time.Now()
	val, exists := AccountListGetCache()
	if exists {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(val))
	} else {
		db := MakeDatabase()
		rows := db.Query("SELECT id, username, email, status, admin FROM users WHERE status = ? OR status = ?", "new", "active")
		defer rows.Close()
		var userId string
		var username string
		var email string
		var status string
		var admin string
		tempJson := ""
		for rows.Next() {
			rows.Scan(&userId, &username, &email, &status, &admin)
			log.Printf("id=%s,user=%s,email=%s,status=%s,admin=%s", userId, username, email, status, admin)
			tempJson = tempJson + makeUserJsonString(userId, username, "[hidden]", status, admin)
			tempJson = tempJson + ","
		}
		accountJson := fmt.Sprintf("[%s", tempJson)
		accountJson = accountJson + fmt.Sprintf("{\"type\": \"last_record\", \"start_time\": \"%s\", \"time_since_start\": \"%s\"}]", start, time.Since(start))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(accountJson))
		AccountListSetCache(accountJson)
	}
}

func ApiListUsersWhoAreActive(w http.ResponseWriter, r *http.Request) {
	Logger("v1/Api/User/ListActive", r, time.Now())
	start := time.Now()
	val, exists := AccountListActiveGetCache()
	if exists {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(val))
	} else {
		db := MakeDatabase()
		rows := db.Query("SELECT id, username, email, status, admin FROM users WHERE status = ?", "active")
		defer rows.Close()
		var userId string
		var username string
		var email string
		var status string
		var admin string
		tempJson := ""
		for rows.Next() {
			rows.Scan(&userId, &username, &email, &status, &admin)
			log.Printf("id=%s,user=%s,email=%s,status=%s,admin=%s", userId, username, email, status, admin)
			tempJson = tempJson + makeUserJsonString(userId, username, "[hidden]", status, admin)
			tempJson = tempJson + ","
		}
		accountJson := fmt.Sprintf("[%s", tempJson)
		accountJson = accountJson + fmt.Sprintf("{\"type\": \"last_record\", \"start_time\": \"%s\", \"time_since_start\": \"%s\"}]", start, time.Since(start))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(accountJson))
		AccountListActiveSetCache(accountJson)
	}
}

func ApiListUsersWhoAreNew(w http.ResponseWriter, r *http.Request) {
	Logger("v1/Api/User/ListActive", r, time.Now())
	start := time.Now()
	val, exists := AccountListNewGetCache()
	if exists {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(val))
	} else {
		db := MakeDatabase()
		rows := db.Query("SELECT id, username, email, status, admin FROM users WHERE status = ?", "new")
		defer rows.Close()
		var userId string
		var username string
		var email string
		var status string
		var admin string
		tempJson := ""
		for rows.Next() {
			rows.Scan(&userId, &username, &email, &status, &admin)
			log.Printf("id=%s,user=%s,email=%s,status=%s,admin=%s", userId, username, email, status, admin)
			tempJson = tempJson + makeUserJsonString(userId, username, "[hidden]", status, admin)
			tempJson = tempJson + ","
		}
		accountJson := fmt.Sprintf("[%s", tempJson)
		accountJson = accountJson + fmt.Sprintf("{\"type\": \"last_record\", \"start_time\": \"%s\", \"time_since_start\": \"%s\"}]", start, time.Since(start))
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(accountJson))
		AccountListNewSetCache(accountJson)
	}
}