package song

import (
	"SongsLibrary/internal/models"
	songservice "SongsLibrary/internal/services/song_service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	Service *songservice.SongStore
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/api/v1/songs", h.GetSongs)
	router.HandlerFunc(http.MethodPost, "/api/v1/songs", h.CreateSong)
	router.HandlerFunc(http.MethodDelete, "/api/v1/songs", h.DeleteSong)
	router.HandlerFunc(http.MethodGet, "/api/v1/songs/:id", h.GetText)
	router.HandlerFunc(http.MethodPut, "/api/v1/songs", h.UpdateSong)
}

// GetSongs godoc
// @Summary Get songs
// @Description Get songs
// @Produce  json
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Param group query string false "group"
// @Param song query string false "song"
// @Param text query string false "text"
// @Param release_date query string false "release date"
// @Param link query string false "link"
// @Success 200 {array} models.Song
// @Failure 400 "Bad request"
// @Failure 404 "Not found"
// @Failure 500 "Internal"
// @Router /api/v1/songs [get]
func (h *Handler) GetSongs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	filter := make(map[string]string, 5)
	filter["group"] = r.URL.Query().Get("group")
	filter["song"] = r.URL.Query().Get("song")
	filter["text"] = r.URL.Query().Get("text")
	filter["link"] = r.URL.Query().Get("link")
	filter["release_date"] = r.URL.Query().Get("release_date")

	songs, status, err := h.Service.GetSongs(r.Context(), r.URL.Query().Get("limit"), r.URL.Query().Get("offset"), filter)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	bytes, err := json.Marshal(songs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(bytes)
}

// CreateSong godoc
// @Summary Create song
// @Description Create song
// @Accept json
// @Produce  json
// @Param request body models.CreateSongReq true "Request"
// @Success 201 "Song created"
// @Failure 409 "Already exists"
// @Failure 500 "Internal"
// @Router /api/v1/songs [post]
func (h *Handler) CreateSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	song := &models.CreateSongReq{}

	if err := json.NewDecoder(r.Body).Decode(song); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	status, err := h.Service.CreateSong(r.Context(), song)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	w.WriteHeader(status)
}

// DeleteSong godoc
// @Summary Delete song
// @Description Delete song
// @Produce  json
// @Param id query int true "Request"
// @Success 204 "OK"
// @Failure 404 "Not found"
// @Failure 500 "Internal"
// @Router /api/v1/songs [delete]
func (h *Handler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	status, err := h.Service.DelSong(r.Context(), int64(id))
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	w.WriteHeader(status)
}

// UpdateSong godoc
// @Summary Update song
// @Description Update song
// @Accept json
// @Produce  json
// @Param request body models.Song true "Request"
// @Success 204 "Updated"
// @Failure 400 "Bad request"
// @Failure 404 "Not found"
// @Failure 500 "Internal"
// @Router /api/v1/songs [put]
func (h *Handler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	song := &models.Song{}

	if err := json.NewDecoder(r.Body).Decode(song); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	status, err := h.Service.UpdateSong(r.Context(), song)
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	w.WriteHeader(status)
}

// GetText godoc
// @Summary Get text
// @Description Get text
// @Produce  json
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Param id path int true "ID"
// @Success 200 {string} song text
// @Failure 404 "Not found"
// @Failure 500 "Internal"
// @Router /api/v1/songs/{id} [get]
func (h *Handler) GetText(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := httprouter.ParamsFromContext(r.Context())
	id, _ := strconv.Atoi(params.ByName("id"))

	defer r.Body.Close()

	text, status, err := h.Service.GetText(r.Context(), int64(id), r.URL.Query().Get("limit"), r.URL.Query().Get("offset"))
	if err != nil {
		http.Error(w, err.Error(), status)
		return
	}

	bytes, err := json.Marshal(text)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(bytes)
}
