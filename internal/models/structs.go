package models

import "time"

type Genre struct {
	ID   int
	Name string
}
type Movie struct {
	ID          int
	Title       string
	ReleaseYear int
	Duration    int
}

type Actor struct {
	ID        int
	Name      string
	BirthDate time.Time
}
