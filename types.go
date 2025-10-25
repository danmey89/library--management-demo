package main

type dbParams struct {
	DbName   string `yaml:"dbName"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Sslmode  string `yaml:"sslmode"`
}

type bookEntry struct {
	ISBN13			int
	Title			string
	Author			string
	Publication_year	int
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
