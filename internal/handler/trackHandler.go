package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"uuid"

	"velo-server/internal/model"
	"velo-server/internal/service"
)

type Handler struct {
	music *service.MusicService
}

func NewHandler(music *service.MusicService) *Handler {
	return &Handler{
		music: music,
	}
}

// API response structure
type response struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

// POST /api/v1/tracks
func (h *Handler) CreateTrack(
	w http.ResponseWriter,
	r *http.Request,
) {
	var track model.Track

	r.Body = http.MaxBytesReader(w, r.Body, 64*1024)

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&track); err != nil {
		writeJSON(w, http.StatusBadRequest, response{
			Error: "Invalid request body",
		})
		return
	}

	created, err := h.music.RegisterTrack(r.Context(), &track)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, response{
			Error: err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusCreated, response{
		Message: "Track created successfully",
		Data:    created,
	})
}

// GET /api/v1/tracks
func (h *Handler) GetTracks(
	w http.ResponseWriter,
	r *http.Request,
) {
	limit := 20
	offset := 0

	if value := r.URL.Query().Get("limit"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed <= 0 {
			writeJSON(w, http.StatusBadRequest, response{
				Error: "Invalid limit",
			})
			return
		}
		limit = parsed
	}

	if value := r.URL.Query().Get("offset"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 0 {
			writeJSON(w, http.StatusBadRequest, response{
				Error: "Invalid offset",
			})
			return
		}
		offset = parsed
	}

	tracks, err := h.music.GetTracks(
		r.Context(),
		limit,
		offset,
	)

	if err != nil {
		log.Printf("list tracks failed: %v", err)
		writeJSON(w, http.StatusInternalServerError, response{
			Error: "Failed to retrieve tracks",
		})
		return
	}

	writeJSON(w, http.StatusOK, response{
		Data: tracks,
	})
}

// GET /api/v1/tracks/{id}
func (h *Handler) GetTrack(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, response{
			Error: "Invalid track ID",
		})
		return
	}

	track, err := h.music.GetTrack(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTrackNotFound) {
			writeJSON(w, http.StatusNotFound, response{
				Error: "Track not found",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, response{
			Error: "Failed to retrieve track",
		})
		return
	}

	writeJSON(w, http.StatusOK, response{
		Data: track,
	})
}

// GET /api/v1/tracks/{id}/stream
func (h *Handler) StreamTrack(
	w http.ResponseWriter,
	r *http.Request,
) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, response{
			Error: "Invalid track ID",
		})
		return
	}

	streamURL, err := h.music.GetStreamURL(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrTrackNotFound) {
			writeJSON(w, http.StatusNotFound, response{
				Error: "Track not found",
			})
			return
		}

		writeJSON(w, http.StatusInternalServerError, response{
			Error: "Failed to generate streaming URL",
		})
		return
	}

	writeJSON(w, http.StatusOK, response{
		Data: map[string]string{
			"stream_url": streamURL,
		},
	})
}

// DELETE /api/v1/tracks/{id}
// func (h *Handler) DeleteTrack(
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) {
// 	id, err := uuid.Parse(r.PathValue("id"))
// 	if err != nil {
// 		writeJSON(w, http.StatusBadRequest, response{
// 			Error: "Invalid track ID",
// 		})
// 		return
// 	}
//
// 	if err := h.music.DeleteTrack(r.Context(), id); err != nil {
// 		if errors.Is(err, service.ErrTrackNotFound) {
// 			writeJSON(w, http.StatusNotFound, response{
// 				Error: "Track not found",
// 			})
// 			return
// 		}
//
// 		writeJSON(w, http.StatusInternalServerError, response{
// 			Error: "Failed to delete track",
// 		})
// 		return
// 	}
//
// 	writeJSON(w, http.StatusOK, response{
// 		Message: "Track deleted successfully",
// 	})
// }
