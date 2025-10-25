CREATE TABLE IF NOT EXISTS books(
	isbn13 BIGINT PRIMARY KEY NOT NULL,
	title VARCHAR,
	author VARCHAR,
	publication_year INT,
	publisher VARCHAR,
	genres VARCHAR 
);
