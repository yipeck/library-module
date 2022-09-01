package author

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/yipeck/library-module/toolkit"

	"github.com/gorilla/mux"
)

type Author struct {
	Id        int
	Name      string
	CountryId int
	UserId    int
	Avatar    string
	Letter    string
	Created   int
}

type Response struct {
	Success bool
	Data    []Author
	Message string
}

type CountResponse struct {
	Success bool
	Data    int
	Message string
}

// Fetch authors
func FetchAuthors(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response
		var sqlStrings = []string{}

		sqlStrings = append(sqlStrings, "SELECT * FROM author")

		query := r.URL.Query()
		search := query.Get("search")
		userId, _ := toolkit.GetToken(r)

		sqlStrings = append(sqlStrings, "WHERE UserId = "+userId)

		if len(search) > 0 {
			sqlStrings = append(sqlStrings, "AND Name LIKE '%"+search+"%'")
		}

		sqlStrings = append(sqlStrings, "ORDER BY Letter ASC")

		rows, err := db.Query(strings.Join(sqlStrings, " "))

		if err != nil {
			response = Response{
				Success: false,
				Message: err.Error(),
			}
		}

		defer rows.Close()

		var authors []Author

		for rows.Next() {
			var author Author

			err := rows.Scan(
				&author.Id,
				&author.Name,
				&author.CountryId,
				&author.UserId,
				&author.Avatar,
				&author.Letter,
				&author.Created,
			)

			if err != nil {
				response = Response{
					Success: false,
					Message: err.Error(),
				}
			}

			authors = append(authors, author)
		}

		response = Response{
			Success: true,
			Data:    authors,
		}

		json.NewEncoder(w).Encode(response)
	}
}

// Add author
func AddAuthor(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		r.ParseMultipartForm(1024 * 1024 * 10)
		name := r.MultipartForm.Value["name"][0]
		countryId := r.MultipartForm.Value["countryId"][0]
		avatar := r.MultipartForm.Value["avatar"][0]
		letter := r.MultipartForm.Value["letter"][0]
		userId, _ := toolkit.GetToken(r)

		_, err := db.Exec(`
			INSERT INTO author (Name, CountryId, UserId, Avatar, Letter, Created) VALUES (?, ?, ?, ?, ?, ?)
		`, name, countryId, userId, avatar, letter, time.Now().Unix())

		if err != nil {
			response = Response{
				Success: false,
				Message: err.Error(),
			}
		} else {
			response = Response{
				Success: true,
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

// Update author
func UpdateAuthor(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		userId, _ := toolkit.GetToken(r)
		r.ParseMultipartForm(1024 * 1024 * 10)
		id := r.MultipartForm.Value["id"][0]
		name := r.MultipartForm.Value["name"][0]
		countryId := r.MultipartForm.Value["countryId"][0]
		avatar := r.MultipartForm.Value["avatar"][0]
		letter := r.MultipartForm.Value["letter"][0]

		_, err := db.Exec(`
			UPDATE author
			SET Name = ?, CountryId = ?, Avatar = ?, Letter = ?
			WHERE Id = ? AND UserId = ?
		`, name, countryId, avatar, letter, id, userId)

		if err != nil {
			response = Response{
				Success: false,
				Message: err.Error(),
			}
		} else {
			response = Response{
				Success: true,
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

// Delete a author
func DeleteAuthor(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		vars := r.URL.Query()
		id := vars["id"]
		userId, _ := toolkit.GetToken(r)
		_, err := db.Exec(`
			DELETE FROM author WHERE Id = ? AND UserId = ?
		`, id[0], userId)

		if err != nil {
			response = Response{
				Success: false,
				Message: err.Error(),
			}
		} else {
			response = Response{
				Success: true,
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}

// Count by country
func CountByCountry(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response CountResponse
		var count int
		userId, _ := toolkit.GetToken(r)
		vars := mux.Vars(r)
		countryId := vars["cid"]

		err := db.QueryRow(`
			SELECT COUNT(*) FROM author WHERE CountryId = ? AND UserId = ?
		`, countryId, userId).Scan(&count)

		response = CountResponse{
			Success: true,
			Data:    count,
		}

		if err != nil {
			response = CountResponse{
				Success: false,
				Data:    0,
				Message: err.Error(),
			}
		}

		json.NewEncoder(w).Encode(response)
	}
}
