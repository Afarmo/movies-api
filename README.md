# 🎬 Movies API
![Go](https://img.shields.io/badge/Go-1.26.5-00ADD8?logo=go&logoColor=white)
## Overview

A RESTful API built with Go and SQLite for managing movies, actors, and genres.

The API provides CRUD operations for all three entities, relationships between movies, actors, and genres, search and filtering functionality, and request validation.

The application is structured into separate handler, service, repository, validation, and database layers.

---

### Prerequisites & Dependencies

- Go 1.26.5
- SQLite driver `github.com/mattn/go-sqlite3`

---

## Setup

### Install & Run

Clone the repository:

```bash
git clone https://github.com/Afarmo/movies-api.git
cd movies-api
```

Install the Go dependencies:

```bash
go mod download
```

The SQLite database is created by the application and is not included in the repository.

Sample data can be loaded using the provided `seed.sql` file:

```bash
sqlite3 movieapi.db < seed.sql
```

Start the API:

```bash
go run ./cmd
```

---

## Usage

The API can be tested using an HTTP client such as Bruno, Postman, or `curl`.

### Get All Movies

```bash
curl http://localhost:8080/api/movies
```

### Get a Movie

```bash
curl http://localhost:8080/api/movies/1
```

### Create a Movie

```bash
curl -X POST http://localhost:8080/api/movies \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Example Movie",
    "releaseyear": 2026,
    "duration": 120,
    "actorids": [1, 2],
    "genreids": [1, 3]
  }'
```

### Update a Movie

Movie updates use `PATCH` and support updating movie information and modifying actor and genre relationships.

```bash
curl -X PATCH http://localhost:8080/api/movies/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Updated Movie",
    "addActorIds": [3],
    "removeActorIds": [1],
    "addGenreIds": [2],
    "removeGenreIds": [4]
  }'
```

### Delete a Movie

```bash
curl -X DELETE http://localhost:8080/api/movies/1
```

---

# Endpoints

## Movies

| Method | Route | Description | Success Status |
|--------|--------|-------------|----------------|
| GET | /api/movies | Get all movies | 200 |
| GET | /api/movies?genre={id} | Filter movies by genre | 200 |
| GET | /api/movies?year={year} | Filter movies by release year | 200 |
| GET | /api/movies?actor={id} | Filter movies by actor | 200 |
| GET | /api/movies/search?title={title} | Search movies by title | 200 |
| GET | /api/movies/{id} | Get a movie by ID | 200 |
| POST | /api/movies | Create a movie | 201 |
| PATCH | /api/movies/{id} | Update a movie | 200 |
| DELETE | /api/movies/{id} | Delete a movie | 204 |

---

## Actors

| Method | Route | Description | Success Status |
|--------|--------|-------------|----------------|
| GET | /api/actors | Get all actors | 200 |
| GET | /api/actors/search?name={name} | Search actors by name | 200 |
| GET | /api/actors/{id} | Get an actor by ID | 200 |
| GET | /api/movies/{movieId}/actors | Get actors associated with a movie | 200 |
| POST | /api/actors | Create an actor | 201 |
| PATCH | /api/actors/{id} | Update an actor | 200 |
| DELETE | /api/actors/{id} | Delete an actor | 204 |

---

## Genres

| Method | Route | Description | Success Status |
|--------|--------|-------------|----------------|
| GET | /api/genres | Get all genres | 200 |
| GET | /api/genres/{id} | Get a genre by ID | 200 |
| POST | /api/genres | Create a genre | 201 |
| PATCH | /api/genres/{id} | Update a genre | 200 |
| DELETE | /api/genres/{id} | Delete a genre | 204 |

---

## Request Formats

### Movie

```json
{
  "title": "Example Movie",
  "releaseyear": 1999,
  "duration": 136,
  "actorids": [1, 2],
  "genreids": [1, 3]
}
```

### Actor

```json
{
  "name": "Example Actor",
  "birthdate": "1970-01-01"
}
```

### Genre

```json
{
  "name": "Action"
}
```

---

## Movie Updates

Movie `PATCH` requests can modify both movie information and its relationships.

```json
{
  "title": "Updated Title",
  "releaseyear": 2026,
  "duration": 120,
  "addActorIds": [1, 2],
  "removeActorIds": [3],
  "addGenreIds": [1],
  "removeGenreIds": [4]
}
```

All fields are optional.

Relationship fields allow actors and genres to be added or removed without replacing the existing relationships.

---

## Features

- RESTful API
- CRUD operations for movies, actors, and genres
- Movie ↔ Actor relationships
- Movie ↔ Genre relationships
- Actor search
- Movie search
- Movie filtering by genre, year, and actor
- Request validation
- SQLite database
- Layered application architecture

## Project Structure

```text
.
├── cmd
│   └── main.go
├── internal
│   ├── database
│   │   └── db.go
│   ├── handlers
│   │   ├── actor.go
│   │   ├── genre.go
│   │   └── movie.go
│   ├── models
│   │   └── structs.go
│   ├── repository
│   │   ├── actor.go
│   │   ├── genre.go
│   │   └── movie.go
│   ├── router
│   │   └── router.go
│   ├── service
│   │   ├── actor.go
│   │   ├── genre.go
│   │   └── movie.go
│   └── validation
│       ├── actor.go
│       ├── genre.go
│       └── movie.go
├── seed.sql
├── go.mod
├── go.sum
└── README.md
```