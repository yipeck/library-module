package book

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/yipeck/library-module/toolkit"
)

type Book struct {
	Id           int
	ISBN         string
	AuthorId     int
	CategoryId   int
	PublisherId  int
	UserId       int
	Title        string
	Picture      string
	PublishYear  int
	PublishMonth int
	ShopYear     int
	ShopMonth    int
	ShopDay      int
	ReadStatus   string
	Letter       string
	Created      int
}

type Response struct {
	Success bool
	Data    []Book
	Message string
}

type CountResponse struct {
	Success bool
	Data    int
	Message string
}

// Fetch books by sql
func FetchBooksBySQL(w http.ResponseWriter, r *http.Request, db *sql.DB, querySQL string) {
	var response Response

	userId, _ := toolkit.GetToken(r)
	rows, err := db.Query(querySQL + " AND UserId = " + userId)

	if err != nil {
		response = Response{
			Success: false,
			Message: err.Error(),
		}
	}

	defer rows.Close()

	var books []Book

	for rows.Next() {
		var book Book

		err := rows.Scan(
			&book.Id,
			&book.ISBN,
			&book.AuthorId,
			&book.CategoryId,
			&book.PublisherId,
			&book.UserId,
			&book.Title,
			&book.Picture,
			&book.PublishYear,
			&book.PublishMonth,
			&book.ShopYear,
			&book.ShopMonth,
			&book.ShopDay,
			&book.ReadStatus,
			&book.Letter,
			&book.Created,
		)

		if err != nil {
			response = Response{
				Success: false,
				Message: err.Error(),
			}
		}

		books = append(books, book)
	}

	response = Response{
		Success: true,
		Data:    books,
	}

	json.NewEncoder(w).Encode(response)
}

// Fetch books bought this month
func FetchBooksBoughtThisMonth(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		year := strconv.Itoa(now.Year())
		month := strconv.Itoa(int(now.Month()))

		querySQL := "SELECT * FROM book WHERE ShopYear = " + year + " And ShopMonth = " + month
		FetchBooksBySQL(w, r, db, querySQL)
	}
}

// Fetch books by type
func FetchBooksByType(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		typeVal := vars["type"]
		value := vars["value"]

		var attrName string

		if typeVal == "status" {
			attrName = "ReadStatus"
		} else if typeVal == "author" {
			attrName = "AuthorId"
		} else if typeVal == "category" {
			attrName = "CategoryId"
		} else if typeVal == "publisher" {
			attrName = "PublisherId"
		}

		querySQL := "SELECT * FROM book WHERE " + attrName + " = " + value
		FetchBooksBySQL(w, r, db, querySQL)
	}
}

// Fetch book by id
func FetchBookById(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response
		var books []Book
		var book Book

		userId, _ := toolkit.GetToken(r)
		vars := mux.Vars(r)
		row := db.QueryRow(`
			SELECT * FROM book WHERE Id = ? AND UserId = ?
		`, vars["id"], userId)
		err := row.Scan(
			&book.Id,
			&book.ISBN,
			&book.AuthorId,
			&book.CategoryId,
			&book.PublisherId,
			&book.UserId,
			&book.Title,
			&book.Picture,
			&book.PublishYear,
			&book.PublishMonth,
			&book.ShopYear,
			&book.ShopMonth,
			&book.ShopDay,
			&book.ReadStatus,
			&book.Letter,
			&book.Created,
		)

		if err != nil {
			response = Response{
				Success: false,
				Message: err.Error(),
			}
		}

		books = append(books, book)
		response = Response{
			Success: true,
			Data:    books,
		}

		json.NewEncoder(w).Encode(response)
	}
}

// Add a book
func AddBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		userId, _ := toolkit.GetToken(r)
		r.ParseMultipartForm(1024 * 1024 * 10)

		ISBN := r.MultipartForm.Value["ISBN"][0]
		authorId := r.MultipartForm.Value["authorId"][0]
		categoryId := r.MultipartForm.Value["categoryId"][0]
		publisherId := r.MultipartForm.Value["publisherId"][0]
		title := r.MultipartForm.Value["title"][0]
		picture := r.MultipartForm.Value["picture"][0]
		publishYear := r.MultipartForm.Value["publishYear"][0]
		publishMonth := r.MultipartForm.Value["publishMonth"][0]
		shopYear := r.MultipartForm.Value["shopYear"][0]
		shopMonth := r.MultipartForm.Value["shopMonth"][0]
		shopDay := r.MultipartForm.Value["shopDay"][0]
		letter := r.MultipartForm.Value["letter"][0]
		readStatus := r.MultipartForm.Value["readStatus"][0]

		_, err := db.Exec(`
			INSERT INTO book (
				ISBN,
				AuthorId,
				CategoryId,
				PublisherId,
				UserId,
				Title,
				Picture,
				PublishYear,
				PublishMonth,
				ShopYear,
				ShopMonth,
				ShopDay,
				Letter,
				ReadStatus,
				Created
			) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
			ISBN,
			authorId,
			categoryId,
			publisherId,
			userId,
			title,
			picture,
			publishYear,
			publishMonth,
			shopYear,
			shopMonth,
			shopDay,
			letter,
			readStatus,
			time.Now().Unix())

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

// Delete a book
func DeleteBook(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		userId, _ := toolkit.GetToken(r)
		vars := mux.Vars(r)
		_, err := db.Exec(`
			DELETE FROM book WHERE Id = ? AND UserId = ?
		`, vars["id"], userId)

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

// Update read status
func UpdateReadStatus(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var response Response

		userId, _ := toolkit.GetToken(r)
		r.ParseMultipartForm(1024 * 1024 * 10)
		id := r.MultipartForm.Value["id"][0]
		readStatus := r.MultipartForm.Value["readStatus"][0]

		_, err := db.Exec(`
			UPDATE book SET ReadStatus = ? WHERE Id = ? AND UserId = ?
		`, readStatus, id, userId)

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

// count books by sql
func CountBooksBySQL(w http.ResponseWriter, r *http.Request, db *sql.DB, querySQL string) {
	var response CountResponse
	var count int

	userId, _ := toolkit.GetToken(r)
	err := db.QueryRow(querySQL + " AND UserId = " + userId).Scan(&count)

	if err != nil {
		response = CountResponse{
			Success: false,
			Data:    0,
			Message: err.Error(),
		}
	}

	response = CountResponse{
		Success: true,
		Data:    count,
	}

	json.NewEncoder(w).Encode(response)
}

// Count by read status
func CountBooksByType(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		typeVal := vars["type"]
		value := vars["value"]

		var attrName string

		if typeVal == "status" {
			attrName = "ReadStatus"
		} else if typeVal == "author" {
			attrName = "AuthorId"
		} else if typeVal == "category" {
			attrName = "CategoryId"
		} else if typeVal == "publisher" {
			attrName = "PublisherId"
		}

		querySQL := "SELECT COUNT(*) FROM book WHERE " + attrName + " = " + value
		CountBooksBySQL(w, r, db, querySQL)
	}
}
