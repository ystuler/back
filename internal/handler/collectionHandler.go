package handler

import (
	"back/internal/exceptions"
	"back/internal/middleware"
	"back/internal/schemas"
	"back/internal/util"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) createCollection(w http.ResponseWriter, r *http.Request) {
	var collectionSchemaReq schemas.CreateCollectionReq

	if err := util.DecodeJSON(w, r, &collectionSchemaReq); err != nil {
		http.Error(w, exceptions.ErrInvalidJSONFormat, http.StatusBadRequest)
		return
	}

	if err := h.validator.ValidateWithDetailedErrors(&collectionSchemaReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		http.Error(w, exceptions.ErrInvalidToken, http.StatusUnauthorized)
		return
	}

	createdCollection, err := h.services.CreateCollection(&collectionSchemaReq, userID)
	if err != nil {
		http.Error(w, exceptions.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createdCollection)
	if err != nil {
		http.Error(w, exceptions.ErrInternalServer, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) editCollection(w http.ResponseWriter, r *http.Request) {
	var updatedCollectionSchema schemas.UpdateCollectionReq

	if err := util.DecodeJSON(w, r, &updatedCollectionSchema); err != nil {
		http.Error(w, exceptions.ErrInvalidJSONFormat, http.StatusBadRequest)
		return
	}

	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
	if err != nil {
		http.Error(w, exceptions.ErrInvalidCollectionID, http.StatusBadRequest)
	}
	updatedCollectionSchema.ID = collectionID

	if err := h.validator.ValidateWithDetailedErrors(&updatedCollectionSchema); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedCollection, err := h.services.UpdateCollection(&updatedCollectionSchema)
	if err != nil {
		http.Error(w, exceptions.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(updatedCollection)
	if err != nil {
		http.Error(w, exceptions.ErrInternalServer, http.StatusInternalServerError)
		return
	}
}

func (h *Handler) removeCollection(w http.ResponseWriter, r *http.Request) {
	var removedCollectionSchema schemas.RemoveCollectionReq

	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
	if err != nil {
		http.Error(w, exceptions.ErrInvalidCollectionID, http.StatusBadRequest)
	}
	removedCollectionSchema.ID = collectionID

	if err := h.validator.ValidateWithDetailedErrors(&removedCollectionSchema); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.services.RemoveCollection(&removedCollectionSchema)
	if err != nil {
		http.Error(w, exceptions.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) getAllCollections(w http.ResponseWriter, r *http.Request) {
	userID, err := middleware.GetUserId(r.Context())
	if err != nil {
		http.Error(w, exceptions.ErrInvalidToken, http.StatusUnauthorized)
	}

	allCollections, err := h.services.GetAllCollections(userID)
	if err != nil {
		http.Error(w, exceptions.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(allCollections)
	if err != nil {
		http.Error(w, exceptions.ErrInternalServer, http.StatusInternalServerError)
		return
	}

}

func (h *Handler) startPractise(w http.ResponseWriter, r *http.Request) {
	var practiseSchemaReq schemas.TrainSchemaReq

	collectionID, err := strconv.Atoi(chi.URLParam(r, "collectionID"))
	if err != nil {
		http.Error(w, exceptions.ErrInvalidCollectionID, http.StatusBadRequest)
		return
	}
	practiseSchemaReq.ID = collectionID

	if err := h.validator.ValidateWithDetailedErrors(&practiseSchemaReq); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	randomCards, err := h.services.TrainCards(&practiseSchemaReq)
	if err != nil {
		http.Error(w, exceptions.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(randomCards)
	if err != nil {
		http.Error(w, exceptions.ErrInternalServer, http.StatusInternalServerError)
		return
	}
}
