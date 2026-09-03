package models

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Movie struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	ReleaseYear int    `json:"releaseyear"`
	Duration    int    `json:"duration"`
	ActorIds    []int  `json:"actorids"`
	GenreIds    []int  `json:"genreids"`
}

type Actor struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birthdate"`
}

type MoviePatch struct {
	Title          *string `json:"title"`
	ReleaseYear    *int    `json:"releaseyear"`
	Duration       *int    `json:"duration"`
	AddActorIds    []int   `json:"addActorIds"`
	RemoveActorIds []int   `json:"removeActorIds"`
	AddGenreIds    []int   `json:"addGenreIds"`
	RemoveGenreIds []int   `json:"removeGenreIds"`
}

type ActorPatch struct {
	Name      *string `json:"name"`
	BirthDate *string `json:"birthdate"`
}

type MovieFilter struct {
	GenreID *int64
	ActorID *int64
	Year    *int
	Page    *int
	Size    *int
}

type ListMoviesResult struct {
	Movies []*Movie
	Total  int
}

type ListMoviesResponse struct {
	Movies     []*Movie `json:"movies"`
	Page       *int     `json:"page,omitempty"`
	Size       *int     `json:"size,omitempty"`
	Total      int      `json:"total"`
	TotalPages int      `json:"totalPages"`
}

type ActorFilter struct {
	Page *int
	Size *int
}

type ListActorsResult struct {
	Actors []*Actor
	Total  int
}

type ListActorsResponse struct {
	Actors     []*Actor `json:"actors"`
	Page       *int     `json:"page,omitempty"`
	Size       *int     `json:"size,omitempty"`
	Total      int      `json:"total"`
	TotalPages int      `json:"totalPages"`
}

type GenreFilter struct {
	Page *int
	Size *int
}

type ListGenresResult struct {
	Genres []*Genre
	Total  int
}

type ListGenresResponse struct {
	Genres     []*Genre `json:"genres"`
	Page       *int     `json:"page,omitempty"`
	Size       *int     `json:"size,omitempty"`
	Total      int      `json:"total"`
	TotalPages int      `json:"totalPages"`
}
