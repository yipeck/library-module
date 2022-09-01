package country

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/yipeck/library-module/toolkit"
)

type Country struct {
	Id     int
	UserId int
	Name   string
}

type Response struct {
	Success bool
	Data    []Country
	Message string
}

// Fetch countries
func FetchCountries(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		userId, _ := toolkit.GetToken(r)
		rows, err := db.Query(`
			SELECT * FROM country WHERE UserId = ?
		`, userId)

		if err != nil {
			response = Response{
				Success: false,
				Message: err.Error(),
			}
		}

		defer rows.Close()

		var countries []Country

		for rows.Next() {
			var country Country

			err := rows.Scan(&country.Id, &country.UserId, &country.Name)

			if err != nil {
				response = Response{
					Success: false,
					Message: err.Error(),
				}
			}

			countries = append(countries, country)
		}

		response = Response{
			Success: true,
			Data:    countries,
		}

		json.NewEncoder(w).Encode(response)
	}
}

// Add a country
func AddCountry(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		r.ParseMultipartForm(1024 * 1024 * 10)
		name := r.MultipartForm.Value["name"][0]
		userId, _ := toolkit.GetToken(r)
		_, err := db.Exec(`
			INSERT INTO country (UserId, Name) VALUES (?, ?)
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

// Update country
func UpdateCountry(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		userId, _ := toolkit.GetToken(r)
		r.ParseMultipartForm(1024 * 1024 * 10)
		id := r.MultipartForm.Value["id"][0]
		name := r.MultipartForm.Value["name"][0]
		_, err := db.Exec(`
			UPDATE country SET Name = ? WHERE Id = ? AND UserId = ?
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

// Delete countries by ids
func DeleteCountry(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		userId, _ := toolkit.GetToken(r)
		vars := r.URL.Query()
		id := vars["id"]
		_, err := db.Exec(`
			DELETE FROM country WHERE Id = ? AND UserId = ?
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
