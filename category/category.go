package category

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/yipeck/library-module/toolkit"
)

type Category struct {
	Id      int
	UserId  int
	Title   string
	Created int
}

type Response struct {
	Success bool
	Data    []Category
	Message string
}

// Fetch categories
func FetchCategories(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		userId, _ := toolkit.GetToken(r)
		rows, err := db.Query(`SELECT * FROM category WHERE UserId = ? ORDER BY Id DESC`, userId)

		if err != nil {
			response = Response{
				Success: false,
				Message: err.Error(),
			}
		}

		defer rows.Close()

		var categories []Category

		for rows.Next() {
			var category Category

			err := rows.Scan(&category.Id, &category.UserId, &category.Title, &category.Created)

			if err != nil {
				response = Response{
					Success: false,
					Message: err.Error(),
				}
			}

			categories = append(categories, category)
		}

		response = Response{
			Success: true,
			Data:    categories,
		}

		json.NewEncoder(w).Encode(response)
	}
}

// Add a category
func AddCategory(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		userId, _ := toolkit.GetToken(r)
		r.ParseMultipartForm(1024 * 1024 * 10)
		title := r.MultipartForm.Value["title"][0]
		_, err := db.Exec(`
			INSERT INTO category (UserId, Title, Created) VALUES (?, ?, ?)
		`, userId, title, time.Now().Unix())

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

// update a category
func UpdateCategory(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		userId, _ := toolkit.GetToken(r)
		r.ParseMultipartForm(1024 * 1024 * 10)
		id := r.MultipartForm.Value["id"][0]
		title := r.MultipartForm.Value["title"][0]
		_, err := db.Exec(`
			UPDATE category SET Title = ? WHERE Id = ? AND UserId = ?
		`, title, id, userId)

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

// Delete a category
func DeleteCategory(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		userId, _ := toolkit.GetToken(r)
		vars := r.URL.Query()
		ids := vars["ids"]
		query, args, _ := sqlx.In(`
			DELETE FROM category WHERE Id IN (?) AND UserId = ?
		`, ids, userId)
		_, err := db.Exec(query, args...)

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
