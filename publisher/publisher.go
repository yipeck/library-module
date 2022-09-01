package publisher

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/yipeck/library-module/toolkit"
)

type Publisher struct {
	Id     int
	UserId int
	Name   string
}

type Response struct {
	Success bool
	Data    []Publisher
	Message string
}

// Fetch publisheres
func FetchPublisheres(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		userId, _ := toolkit.GetToken(r)
		rows, err := db.Query(`SELECT * FROM publisher WHERE UserId = ?`, userId)

		if err != nil {
			response = Response{
				Success: false,
				Message: err.Error(),
			}
		}

		defer rows.Close()

		var publisheres []Publisher

		for rows.Next() {
			var publisher Publisher

			err := rows.Scan(&publisher.Id, &publisher.UserId, &publisher.Name)

			if err != nil {
				response = Response{
					Success: false,
					Message: err.Error(),
				}
			}

			publisheres = append(publisheres, publisher)
		}

		response = Response{
			Success: true,
			Data:    publisheres,
		}

		json.NewEncoder(w).Encode(response)
	}
}

// Add a publisher
func AddPublisher(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		userId, _ := toolkit.GetToken(r)
		r.ParseMultipartForm(1024 * 1024 * 10)
		name := r.MultipartForm.Value["name"][0]
		_, err := db.Exec(`
			INSERT INTO publisher (UserId, Name) VALUES (?, ?)
		`, userId, name)

		if err != nil {
			response = Response{
				Success: false,
				Message: err.Error(),
			}
		}

		response = Response{
			Success: true,
		}

		json.NewEncoder(w).Encode(response)
	}
}

// Update publisher
func UpdatePublisher(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		userId, _ := toolkit.GetToken(r)
		r.ParseMultipartForm(1024 * 1024 * 10)
		id := r.MultipartForm.Value["id"][0]
		name := r.MultipartForm.Value["name"][0]
		_, err := db.Exec(`
			UPDATE publisher SET Name = ? WHERE Id = ? AND UserId = ?
		`, name, id, userId,
		)

		if err != nil {
			response = Response{
				Success: false,
				Message: err.Error(),
			}
		}

		response = Response{
			Success: true,
		}

		json.NewEncoder(w).Encode(response)
	}
}

// Delete publisheres by ids
func DeletePublisher(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		userId, _ := toolkit.GetToken(r)
		vars := r.URL.Query()
		id := vars["id"]
		_, err := db.Exec(`
			DELETE FROM publisher WHERE Id = ? AND UserId = ?
		`, id[0], userId)

		if err != nil {
			response = Response{
				Success: false,
				Message: err.Error(),
			}
		}

		response = Response{
			Success: true,
		}

		json.NewEncoder(w).Encode(response)
	}
}
