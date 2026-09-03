package models

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
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
