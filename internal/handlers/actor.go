package handlers

import (
	"encoding/json"
	"movies-api/internal/apperrors"
	"movies-api/internal/models"
	"movies-api/internal/service"
	"movies-api/internal/validation"
	"net/http"
	"strconv"
)

type ActorHandler struct {
	actorService *service.ActorService
}

func NewActorHandler(actorService *service.ActorService) *ActorHandler {
	return &ActorHandler{actorService: actorService}
}

func (h *ActorHandler) GetAllActorsHandler(w http.ResponseWriter, req *http.Request) {
	actors, err := h.actorService.ListAllActors(req.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actors)
}

func (h *ActorHandler) GetOneActorHandler(w http.ResponseWriter, req *http.Request) {
	actorID, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid actor id", http.StatusBadRequest)
		return
	}

	actor, err := h.actorService.ListOneActor(req.Context(), actorID)
	if err != nil {
		apperrors.WriteError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actor)
}

func (h *ActorHandler) CreateActorHandler(w http.ResponseWriter, req *http.Request) {
	var actor models.Actor

	if err := json.NewDecoder(req.Body).Decode(&actor); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if err := validation.ValidateActor(actor); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := h.actorService.InsertActor(req.Context(), &actor); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(actor)
}

func (h *ActorHandler) DeleteActorHandler(w http.ResponseWriter, req *http.Request) {
	idString := req.PathValue("id")

	actorID, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return
	}

	force := req.URL.Query().Get("force") == "true"

	if err = h.actorService.DeleteActor(req.Context(), actorID, force); err != nil {
		apperrors.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ActorHandler) UpdateActorHandler(w http.ResponseWriter, req *http.Request) {
	actorID, err := strconv.ParseInt(req.PathValue("id"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid actor ID", http.StatusBadRequest)
		return
	}

	var actor models.ActorPatch

	if err := json.NewDecoder(req.Body).Decode(&actor); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.actorService.UpdateActor(req.Context(), actorID, &actor); err != nil {
		apperrors.WriteError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ActorHandler) SearchActorsHandler(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")

	switch {
	case name != "":
		actors, err := h.actorService.SearchActors(req.Context(), name)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(actors)
	default:
		actors, err := h.actorService.ListAllActors(req.Context())
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(actors)
	}
}

func (h *ActorHandler) GetActorsByMovie(w http.ResponseWriter, req *http.Request) {
	movieId, err := strconv.ParseInt(req.PathValue("movieId"), 10, 64)
	actors, err := h.actorService.ActorsByMovie(req.Context(), movieId)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(actors)
}
