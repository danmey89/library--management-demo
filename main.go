package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-yaml/yaml"
	_ "github.com/lib/pq"
)

type dbParams struct {
	DbName   string `yaml:"dbName"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Sslmode  string `yaml:"sslmode"`
}

type bookEntry struct {
	ID			int
	Title			string
	Author			string
	ISBN			string
	ISBN13			int
	Publication_date	time.Time
	Publisher		string
	Genres			string
}

type Book struct {
	Title			string
	Author			[]string
	ISBN			string
	ISBN13			int
	Publication_date	time.Time
	Publisher		string
	Genres			[]string

}

var db *sql.DB

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
	//defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to open db connection: %s", err)
	}
	return nil

}

func getBooks(author string) ([]Book, error) {
	rows, err := db.Query("SELECT * FROM books WHERE author LIKE $1", author)
	if err != nil {
		return nil, fmt.Errorf("Query error: %s", err)
	}
	defer rows.Close()

	var books []Book

	for rows.Next() {
		var bookItem bookEntry
		var authorList []string
		var genreList []string

		if err := rows.Scan(&bookItem.ID, &bookItem.Title, &bookItem.Author, &bookItem.ISBN, &bookItem.ISBN13,
			&bookItem.Publication_date, &bookItem.Publisher, &bookItem.Genres); err != nil {
			return nil, fmt.Errorf("Error processing query: %s", err)
		}
		authorList = strings.Split(bookItem.Author, "/")
		genreList = strings.Split(bookItem.Genres, "/")

		var book = Book {bookItem.Title, authorList,bookItem.ISBN,bookItem.ISBN13,bookItem.Publication_date,
			bookItem.Publisher, genreList}
		books = append(books, book)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

func main() {
	if err := connectDB(); err != nil {
		log.Fatal(err)
	}
	name := "%Douglas Adams%"
	
	books, err := getBooks(name)
	if err != nil {
		log.Fatal(err)
	}
}
