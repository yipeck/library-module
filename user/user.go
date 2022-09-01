package user

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/yipeck/library-module/toolkit"
)

type User struct {
	Id       int
	Email    string
	PassWord string
	NickName string
	Slogan   string
	Avatar   string
	Created  int
}

type Response struct {
	Success bool
	Data    User
	Message string
}

type SignInResponse struct {
	Success bool
	Token   string
	Message string
}

// Fetch user by id
func FetchUserById(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response
		var user User

		userId, _ := toolkit.GetToken(r)
		row := db.QueryRow(`
			SELECT Id, Email, NickName, Slogan, Avatar, Created
			FROM user
			WHERE Id = ?
		`, userId)
		err := row.Scan(
			&user.Id,
			&user.Email,
			&user.NickName,
			&user.Slogan,
			&user.Avatar,
			&user.Created,
		)

		response = Response{
			Success: true,
			Data:    user,
		}

		if err != nil {
			response = Response{
				Success: false,
				Message: err.Error(),
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

// Is exist
func IsExist(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response
		var count int

		r.ParseMultipartForm(1024 * 1024 * 10)
		email := r.MultipartForm.Value["email"][0]

		err := db.QueryRow(`
			SELECT COUNT(*) FROM user WHERE Email = ?
		`, email).Scan(&count)

		if count == 0 {
			response = Response{
				Success: false,
			}
		} else {
			response = Response{
				Success: true,
			}
		}

		if err != nil {
			response = Response{
				Success: false,
				Message: err.Error(),
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

// Sign in
func SignIn(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response SignInResponse

		r.ParseMultipartForm(1024 * 1024 * 10)
		email := r.MultipartForm.Value["email"][0]
		password := r.MultipartForm.Value["password"][0]
		encryptedPwd := toolkit.EncryptByMD5(password)

		var user User

		row := db.QueryRow(`
			SELECT Id, Email, NickName, Slogan, Avatar, Created
			FROM user
			WHERE Email = ? AND Password = ?
		`, email, encryptedPwd)
		err := row.Scan(
			&user.Id,
			&user.Email,
			&user.NickName,
			&user.Slogan,
			&user.Avatar,
			&user.Created,
		)

		userId := strconv.Itoa(user.Id)
		timestamp := strconv.Itoa(int(time.Now().Unix()))

		token := toolkit.Encrypt(userId) + "@" + toolkit.Encrypt(timestamp)
		response = SignInResponse{
			Success: true,
			Token:   token,
		}

		if err != nil {
			response = SignInResponse{
				Success: false,
				Message: err.Error(),
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

// Sign up
func SignUp(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		r.ParseMultipartForm(1024 * 1024 * 10)
		email := r.MultipartForm.Value["email"][0]
		password := r.MultipartForm.Value["password"][0]
		encryptedPwd := toolkit.EncryptByMD5(password)

		_, err := db.Exec(`
			INSERT INTO user (Email, Password, Created) VALUES (?, ?, ?)
		`, email, encryptedPwd, time.Now().Unix())

		response = Response{
			Success: true,
		}

		if err != nil {
			response = Response{
				Success: false,
				Message: err.Error(),
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

// Update user
func UpdateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		userId, _ := toolkit.GetToken(r)
		r.ParseMultipartForm(1024 * 1024 * 10)
		name := r.MultipartForm.Value["name"][0]
		value := r.MultipartForm.Value["value"][0]
		_, err := db.Exec(
			"UPDATE user SET "+name+" = ? WHERE Id = ?",
			value, userId,
		)

		response = Response{
			Success: true,
		}

		if err != nil {
			response = Response{
				Success: false,
				Message: err.Error(),
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

// Delete a user
func DeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		vars := r.URL.Query()
		id := vars["id"]
		_, err := db.Exec(`DELETE FROM user WHERE Id = ?`, id[0])

		response = Response{
			Success: true,
		}

		if err != nil {
			response = Response{
				Success: false,
				Message: err.Error(),
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}
