package handlers

import (
	"database/sql"
	"encoding/json"
	"movies-api/internal/models"
	"movies-api/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type ActorHandler struct {
	actorService *service.ActorService
}

func NewActorHandler(actorService *service.ActorService) *ActorHandler {
	return &ActorHandler{actorService: actorService}
}

func (h *ActorHandler) GetOneActorHandler(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()

	idString := req.PathValue("id")
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, "invalid actor id", http.StatusBadRequest)
		return
	}

	actor, err := h.actorService.ListOneActor(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			http.Error(w, "actor not found", http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actor)
}

func (h *ActorHandler) GetAllActorsHandler(w http.ResponseWriter, req *http.Request) {

	ctx := req.Context()
	actors, err := h.actorService.ListAllActors(ctx)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actors)

}

func (h *ActorHandler) CreateActorHandler(w http.ResponseWriter, req *http.Request) {
	actor := models.Actor{}
	err := json.NewDecoder(req.Body).Decode(&actor)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	ctx := req.Context()
	err = h.actorService.InsertActor(ctx, &actor)
	if err != nil {
		http.Error(w, "could not create actor", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(actor)
}

func (h *ActorHandler) GetActorByNameHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	name := req.PathValue("name")

	actors, err := h.actorService.ActorsByName(ctx, name)
	if err != nil {
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			http.Error(w, "No Actor Found With the name : "+name, http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actors)
}

func (h *ActorHandler) DeleteActor(w http.ResponseWriter, req *http.Request) {
	idString := req.PathValue("id")
	ctx := req.Context()

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	err = h.actorService.DeleteActor(ctx, id)
	if err != nil {
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			http.Error(w, "NO such Actor", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func (h *ActorHandler) UpdateActor(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	idString := req.PathValue("id")

	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	actor := models.Actor{}
	actor.ID = int(id)
	err = json.NewDecoder(req.Body).Decode(&actor)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	err = h.actorService.UpdateActor(ctx, &actor)
	if err != nil {
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			http.Error(w, "NO such Actor", http.StatusNotFound)
			return
		}

		http.Error(w, "Internal Server Error ", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusAccepted)

}

func (h *ActorHandler) ActorsByMovieHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	idString := req.PathValue("movieId")
	movieid, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, "Invalid id", http.StatusBadRequest)
		return
	}
	actors, err := h.actorService.ActorsByMovie(ctx, movieid)
	if err != nil {
		if strings.Contains(err.Error(), sql.ErrNoRows.Error()) {
			http.Error(w, "No Actor Found :", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actors)
}
