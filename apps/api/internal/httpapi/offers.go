package httpapi

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/eshche-est/eshche-est/apps/api/internal/store"
	"github.com/go-chi/chi/v5"
)

type offerInput struct {
	Title                string    `json:"title"`
	Merchant             string    `json:"merchant"`
	Category             string    `json:"category"`
	Description          string    `json:"description"`
	Contents             string    `json:"contents"`
	ImageURL             string    `json:"imageUrl"`
	PriceKopecks         int64     `json:"priceKopecks"`
	OriginalPriceKopecks int64     `json:"originalPriceKopecks"`
	PickupStart          time.Time `json:"pickupStart"`
	PickupEnd            time.Time `json:"pickupEnd"`
	Quantity             int       `json:"quantity"`
	Address              string    `json:"address"`
	District             string    `json:"district"`
	Latitude             *float64  `json:"latitude"`
	Longitude            *float64  `json:"longitude"`
	Delivery             bool      `json:"delivery"`
	Status               string    `json:"status"`
}

func (s *Server) listOffers(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	offers, err := s.store.ListOffers(r.Context(), store.OfferFilter{
		Query: r.URL.Query().Get("q"), Category: r.URL.Query().Get("category"),
		Limit: limit, Offset: offset,
	})
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": offers})
}

func (s *Server) getOffer(w http.ResponseWriter, r *http.Request) {
	offerID := chi.URLParam(r, "offerID")
	if !idPattern.MatchString(offerID) {
		writeError(w, http.StatusBadRequest, "invalid_offer_id", "Некорректный идентификатор предложения.")
		return
	}
	offer, err := s.store.OfferByID(r.Context(), offerID)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "offer_not_found", "Предложение не найдено.")
		return
	}
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"offer": offer})
}

func (s *Server) createOffer(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	if user.Role != "admin" && user.Role != "merchant" {
		writeError(w, http.StatusForbidden, "insufficient_role", "Для публикации нужна роль продавца или администратора.")
		return
	}
	var input offerInput
	if !decodeJSON(w, r, &input) || !validateOffer(w, input) {
		return
	}
	id, err := store.NewID()
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	offer := input.toOffer()
	offer.ID = id
	offer.CreatedBy = user.ID
	created, err := s.store.CreateOffer(r.Context(), offer)
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"offer": created})
}

func (s *Server) updateOffer(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	offerID := chi.URLParam(r, "offerID")
	if !idPattern.MatchString(offerID) {
		writeError(w, http.StatusBadRequest, "invalid_offer_id", "Некорректный идентификатор предложения.")
		return
	}
	existing, err := s.store.OfferByID(r.Context(), offerID)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "offer_not_found", "Предложение не найдено.")
		return
	}
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	if user.Role != "admin" && (user.Role != "merchant" || existing.CreatedBy != user.ID) {
		writeError(w, http.StatusForbidden, "offer_forbidden", "У вас нет доступа к этому предложению.")
		return
	}
	var input offerInput
	if !decodeJSON(w, r, &input) || !validateOffer(w, input) {
		return
	}
	offer := input.toOffer()
	offer.ID = existing.ID
	updated, err := s.store.UpdateOffer(r.Context(), offer)
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"offer": updated})
}

func (s *Server) deleteOffer(w http.ResponseWriter, r *http.Request) {
	user := currentUser(r)
	offerID := chi.URLParam(r, "offerID")
	if !idPattern.MatchString(offerID) {
		writeError(w, http.StatusBadRequest, "invalid_offer_id", "Некорректный идентификатор предложения.")
		return
	}
	existing, err := s.store.OfferByID(r.Context(), offerID)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusNotFound, "offer_not_found", "Предложение не найдено.")
		return
	}
	if err != nil {
		s.internalError(w, r, err)
		return
	}
	if user.Role != "admin" && (user.Role != "merchant" || existing.CreatedBy != user.ID) {
		writeError(w, http.StatusForbidden, "offer_forbidden", "У вас нет доступа к этому предложению.")
		return
	}
	if err := s.store.DeleteOffer(r.Context(), existing.ID); err != nil {
		s.internalError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func validateOffer(w http.ResponseWriter, input offerInput) bool {
	statusAllowed := input.Status == "draft" || input.Status == "active" || input.Status == "paused"
	validCoordinates := (input.Latitude == nil || (*input.Latitude >= -90 && *input.Latitude <= 90)) &&
		(input.Longitude == nil || (*input.Longitude >= -180 && *input.Longitude <= 180))
	if len([]rune(strings.TrimSpace(input.Title))) < 2 || len([]rune(input.Title)) > 120 ||
		len([]rune(strings.TrimSpace(input.Merchant))) < 2 || len([]rune(input.Merchant)) > 120 ||
		len([]rune(strings.TrimSpace(input.Category))) < 2 || len([]rune(input.Category)) > 80 ||
		input.PriceKopecks <= 0 || input.OriginalPriceKopecks < input.PriceKopecks ||
		!input.PickupEnd.After(input.PickupStart) || input.Quantity < 0 ||
		len([]rune(strings.TrimSpace(input.Address))) < 2 || !statusAllowed || !validCoordinates {
		writeError(w, http.StatusUnprocessableEntity, "validation_error", "Проверьте поля предложения, цену, количество и окно получения.")
		return false
	}
	return true
}

func (input offerInput) toOffer() store.Offer {
	return store.Offer{
		Title: strings.TrimSpace(input.Title), Merchant: strings.TrimSpace(input.Merchant),
		Category: strings.TrimSpace(input.Category), Description: strings.TrimSpace(input.Description),
		Contents: strings.TrimSpace(input.Contents), ImageURL: strings.TrimSpace(input.ImageURL),
		PriceKopecks: input.PriceKopecks, OriginalPriceKopecks: input.OriginalPriceKopecks,
		PickupStart: input.PickupStart, PickupEnd: input.PickupEnd, Quantity: input.Quantity,
		Address: strings.TrimSpace(input.Address), District: strings.TrimSpace(input.District),
		Latitude: input.Latitude, Longitude: input.Longitude, Delivery: input.Delivery,
		Status: input.Status,
	}
}
