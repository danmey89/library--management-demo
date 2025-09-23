CREATE TABLE IF NOT EXISTS books(
	id SERIAL PRIMARY KEY NOT NULL,
	title VARCHAR,
	author VARCHAR,
	isbn VARCHAR,
	isbn13 BIGINT,
	publication_date DATE,
	publisher VARCHAR,
	genres VARCHAR 
);
