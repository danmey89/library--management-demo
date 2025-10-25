package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/go-yaml/yaml"
	_ "github.com/lib/pq"
)

var (
	db    *sql.DB
	books []map[string]string
)

func main() {

	if err := connectDB(); err != nil {
		log.Fatal(err)
	}
	serve()
}

func serve() {

	mux := http.NewServeMux()

	var fs = http.FileServer(http.Dir("./static"))
	var responseTemplate = template.Must(template.New("response.gohtml").Funcs(funcMap).ParseFiles("templates/response.gohtml"))

	mux.Handle("/static/", http.StripPrefix("/static", fs))

	mux.HandleFunc("/", serveTemplate)
	mux.HandleFunc("/request", requestHandler(responseTemplate))
	mux.HandleFunc("/inputBook", inputHandler)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func serveTemplate(w http.ResponseWriter, r *http.Request) {

	p := filepath.Clean(r.URL.Path)

	if p == "/" {
		p = "/index"
	}

	layoutPath := filepath.Join("templates", "layout.gohtml")
	templatePath := filepath.Join("templates", p) + ".html"

	info, err := os.Stat(templatePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}
	}

	if info.IsDir() {
		http.NotFound(w, r)
		return
	}
	
	tmpl, err := template.ParseFiles(layoutPath, templatePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}	

}

func requestHandler(temp *template.Template) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		arguments := ArgumentEvent{
			Selector1: r.Form.Get("selector1"),
			Input1:    "%" + r.Form.Get("input1") + "%",
			Selector2: r.Form.Get("selector2"),
			Input2:    "%" + r.Form.Get("input2") + "%",
		}

		if err := makeQuery(arguments); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := temp.Execute(w, books); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func inputHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	isbn, err := strconv.Atoi(r.Form.Get("ISBN13"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	year, err := strconv.Atoi(r.Form.Get("year"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	author := strings.ReplaceAll(r.Form.Get("author"), ", ", "/")
	genre := strings.ReplaceAll(r.Form.Get("genres"), ", ", "/")

	var newEntry = bookEntry{isbn, r.Form.Get("title"), author, year, r.Form.Get("publisher"), genre}

	if err := insertRow(newEntry); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, r.Header.Get("Referer"), 302)

}

func connectDB() error {

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
		var author string
		var genre string

		if err := rows.Scan( &bookItem.ISBN13, &bookItem.Title, &bookItem.Author, &bookItem.Publication_year, &bookItem.Publisher, &bookItem.Genres); err != nil {
			return fmt.Errorf("Error processing query: %s", err)
		}

		author = strings.ReplaceAll(bookItem.Author, "/", ", ")
		genre = strings.ReplaceAll(bookItem.Genres, "/", ", ")

		book := map[string]string{"title": bookItem.Title, "author": author, "ISBN13": strconv.Itoa(bookItem.ISBN13), "year": strconv.Itoa(bookItem.Publication_year),
			"publisher": bookItem.Publisher, "genre": genre}
		books = append(books, book)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

func makeQuery(arguments ArgumentEvent) error {

	if arguments.Selector2 == "" && arguments.Selector1 != "" {
		query := fmt.Sprintf(`SELECT * FROM books WHERE lower(%s) LIKE lower($1)`, arguments.Selector1)

		rows, err := db.Query(query, arguments.Input1)
		if err != nil {
			return fmt.Errorf("Query error: %s", err)
		}

		defer rows.Close()
		parseRows(rows)

	} else if arguments.Selector2 != "" && arguments.Input2 != "" {
		query := fmt.Sprintf(`SELECT * FROM books WHERE lower(%s) LIKE lower($1) AND lower(%s) LIKE lower($2)`, arguments.Selector1, arguments.Selector2)

		rows, err := db.Query(query, arguments.Input1, arguments.Input2)
		if err != nil {
			return fmt.Errorf("Query error: %s", err)
		}

		defer rows.Close()
		parseRows(rows)
	}

	return nil
}

func insertRow(entry bookEntry) error {

	sqlStatement := `
		INSERT INTO books (isbn13, title, author, publication_year, publisher, genres)
		VALUES ($1, $2, $3, $4, $5, $6)`

	if _, err := db.Exec(sqlStatement, entry.ISBN13, entry.Title, entry.Author, entry.Publication_year, entry.Publisher, entry.Genres); err != nil {
		return fmt.Errorf("SQL error: %s", err)
	}

	return nil
}
