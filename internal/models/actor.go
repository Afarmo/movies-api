package models

type Actor struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	BirthDate string `json:"birthdate"`
}

type ActorPatch struct {
	Name      *string `json:"name"`
	BirthDate *string `json:"birthdate"`
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
