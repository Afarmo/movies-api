package models

import "time"

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseYear int    `json:"releaseyear"`
	Duration    int    `json:"duration"`
	ActorIds    []int  `json:actorids`
	GenreIds    []int  `json:genreids`
}

type Actor struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	BirthDate time.Time `json:"birthdate"`
}

type MoviePatch struct {
	Title       *string `json:"title"`
	ReleaseYear *int    `json:"releaseyear"`
	Duration    *int    `json:"duration"`
}

type ActorPatch struct {
	Name      *string    `json:"name"`
	BirthDate *time.Time `json:"birthdate"`
}
