package database

import "database/sql"

const Query = `
CREATE TABLE IF NOT EXISTS Genre (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS Movie (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    releaseYear INTEGER NOT NULL,
	duration INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS Actor (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    birthDate TEXT NOT NULL
    
);

CREATE TABLE IF NOT EXISTS Movie_Genres (
	movie_id INTEGER NOT NULL,
	genre_id INTEGER NOT NULL,

	PRIMARY KEY (movie_id, genre_id),

    FOREIGN KEY (movie_id) REFERENCES Movie(id) ON DELETE CASCADE,
    FOREIGN KEY (genre_id) REFERENCES Genre(id)
); 

CREATE TABLE IF NOT EXISTS Movie_Actors (
	movie_id INTEGER NOT NULL,
	actor_id INTEGER NOT NULL,

	PRIMARY KEY (movie_id, actor_id),

    FOREIGN KEY (movie_id) REFERENCES Movie(id) ON DELETE CASCADE,
    FOREIGN KEY (actor_id) REFERENCES Actor(id)
); 
`

func NewDataBase(dbName string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func InitializeQuery(db *sql.DB) error {
	_, err := db.Exec(Query)
	return err
}
