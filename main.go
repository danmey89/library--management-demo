package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/go-yaml/yaml"
	_ "github.com/lib/pq"
)

var(
	db *sql.DB
	templates = template.Must(template.ParseFiles("templates/index.html"))
	books []Book
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	books = nil
	p := Page{
		Title: "index",
		Data: books,
	}
	if err := templates.ExecuteTemplate(w, "index.html", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	req := "%" + query.Get("author") + "%"

	if err := querryAuthor(req); err != nil {
		log.Fatal(err)
	}
	p := Page{
		Title: "index",
		Data: books,
	}
	fmt.Println(req)
	if err := templates.ExecuteTemplate(w, "index.html", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func serve() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/request", requestHandler)
	
	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func connectDB() (error) {
	var config dbParams

	rf, err := os.ReadFile("config.yaml")
	if err != nil {
		return fmt.Errorf("Error reading config file: %s", err)
	}
	if err := yaml.Unmarshal(rf, &config); err != nil {
		return fmt.Errorf("Error unmarshalling config file: %s", err)
	}
	conn := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=%s",
		config.Host, config.DbName, config.User, config.Password, config.Sslmode)

	db, err = sql.Open("postgres", conn)
	if err != nil {
		return fmt.Errorf("unable to use configuration: %s", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to open db connection: %s", err)
	}
	return nil

}

func parseRows(rows *sql.Rows) error {
	books = nil
	for rows.Next() {
		var bookItem bookEntry
		var authorList []string
		var genreList []string

		if err := rows.Scan(&bookItem.ID, &bookItem.Title, &bookItem.Author, &bookItem.ISBN, &bookItem.ISBN13,
			&bookItem.Publication_date, &bookItem.Publisher, &bookItem.Genres); err != nil {
			return fmt.Errorf("Error processing query: %s", err)
		}
		authorList = strings.Split(bookItem.Author, "/")
		genreList = strings.Split(bookItem.Genres, "/")

		var book = Book {bookItem.Title, authorList,bookItem.ISBN,bookItem.ISBN13,bookItem.Publication_date,
			bookItem.Publisher, genreList}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	return nil
}

func querryAuthor(author string) (error) {
	rows, err := db.Query("SELECT * FROM books WHERE lower(author) LIKE lower($1)", author)
	if err != nil {
		return fmt.Errorf("Query error: %s", err)
	}
	defer rows.Close()
	parseRows(rows)
	return nil
}

func querryTitle(title string) (error) {
	rows, err := db.Query("SELECT * FROM books WHERE title LIKE $1", title)
	if err != nil {
		return fmt.Errorf("Query error: %s", err)
	}
	defer rows.Close()
	parseRows(rows)
	return nil
}

func main() {
	if err := connectDB(); err != nil {
		log.Fatal(err)
	}
	/*
	name := "%Douglas Adams%"
	
	if err := querryAuthor(name); err != nil {
		log.Fatal(err)
	}
	*/
	serve()
}
