package main

import "time"

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

type Page struct {
	Title string
	Data []Book
}
