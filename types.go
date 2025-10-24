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
	Author			string
	ISBN			string
	ISBN13			string	
	Publication_date	string
	Publisher		string
	Genres			string

}

type Page struct {
	Title string
	Data []map[string]string
}

type ArgumentEvent struct {
	Selector1	string
	Input1		string
	Selector2	string
	Input2		string
}
