package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pelicanch1k/EffectiveMobileTestTask/internal/dto"
)

type pagination struct {
	limit, offset int
}

// @Summary Get songs
// @Tags songs
// @Description Get a list of songs with optional filters and pagination
// @ID getSongs
// @Accept  json
// @Produce  json
// @Param  id header int false "Filter by ID"
// @Param  genre header string false "Filter by genre"
// @Param  song header string false "Filter by song name"
// @Param  releaseDate header string false "Filter by release date"
// @Param  text header string false "Filter by text content"
// @Param  link header string false "Filter by link"
// @Param  groupId header int false "Filter by group ID"
// @Param  group header string false "Filter by group name"
// @Param  limit query int false "Pagination limit"
// @Param  offset query int false "Pagination offset"
// @Success 200 {array} model.Song
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/songs [get]
func (h *Handler) GetSongs(c *gin.Context) {
	id, _ := strconv.Atoi(c.GetHeader("id"))
	genre := c.GetHeader("genre")
	song := c.GetHeader("song")
	releaseDate := c.GetHeader("releaseDate")
	text := c.GetHeader("text")
	link := c.GetHeader("link")
	groupId, _ := strconv.Atoi(c.GetHeader("groupId"))
	group := c.GetHeader("group")

	pag, err := h.initPagination(c)
	if err != nil {
		return
	}

	resp := dto.GetSongsRequest{
		Id:          id,
		Genre:       genre,
		Song:        song,
		ReleaseDate: releaseDate,
		Text:        text,
		Link:        link,
		GroupId:     groupId,
		Group:       group,
		Limit:       pag.limit,
		Offset:      pag.offset,
	}

	h.logger.Info("Fetching songs with filters: ", resp)

	songs, err := h.services.GetSongs(resp)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, songs)
}

// @Summary Get song by ID
// @Tags songs
// @Description Get a song by its ID
// @ID getSongById
// @Accept  json
// @Produce  json
// @Param  id path int true "Song ID"
// @Success 200 {object} model.Song
// @Failure 400 {object} errorResponse
// @Failure 404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/song/{id} [get]
func (h *Handler) GetSongById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "Invalid Song ID")
		return
	}

	song, err := h.services.GetSongById(id)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, song)
}

// @Summary Get song lyrics
// @Security ApiKeyAuth
// @Tags songs
// @Description Get lyrics of a song by its ID with optional pagination
// @ID getSongLyrics
// @Accept  json
// @Produce  json
// @Param  id path int true "Song ID"
// @Param  limit header int false "Pagination limit"
// @Param  offset header int false "Pagination offset"
// @Success 200 {object} map[string][]string "Lyrics of the song"
// @Failure 400 {object} errorResponse "Invalid Song ID"
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/song/{id}/lyrics [get]
func (h *Handler) GetSongLyrics(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "Invalid Song ID")
		return
	}

	pag, err := h.initPagination(c)
	if err != nil {
		return
	}

	resp := dto.GetSongLyricsRequest{
		Id:     id,
		Limit:  pag.limit,
		Offset: pag.offset,
	}

	verses, err := h.services.GetSongLyrics(resp)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"verses": verses,
	})
}

// @Summary Delete a song
// @Security ApiKeyAuth
// @Tags songs
// @Description Delete an existing song
// @ID deleteSong
// @Accept  json
// @Produce  json
// @Param  id path int true "Song ID"
// @Success 200 {object} map[string]interface{} "Song deleted successfully"
// @Failure 400 {object} errorResponse "Invalid ID"
// @Failure 404 {object} errorResponse "Song not found"
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/song/{id} [delete]
func (h *Handler) DeleteSong(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, "Invalid ID")
		return
	}

	if err = h.services.DeleteSong(id); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": "Song deleted successfully"})
}

// @Summary Update a song
// @Security ApiKeyAuth
// @Tags songs
// @Description Update an existing song
// @ID updateSong
// @Accept  json
// @Produce  json
// @Param  song body dto.UpdateSongRequest true "Updated song details"
// @Success 200 {object} map[string]interface{} "Song updated successfully"
// @Failure 400 {object} errorResponse "Invalid JSON"
// @Failure 404 {object} errorResponse "Song not found"
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/song [put]
func (h *Handler) UpdateSong(c *gin.Context) {
	var song dto.UpdateSongRequest

	if err := c.BindJSON(&song); err != nil {
		h.newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.UpdateSong(song); err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{"message": "Song updated successfully"})
}

// @Summary Add a new song
// @Security ApiKeyAuth
// @Tags songs
// @Description Add a new song to the catalog
// @ID addSong
// @Accept  json
// @Produce  json
// @Param  song body dto.AddSongRequest true "Details of the song to add"
// @Success 201 {object} map[string]interface{} "Song created successfully"
// @Failure 400 {object} errorResponse "Invalid JSON"
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/song [post]
func (h *Handler) AddSong(c *gin.Context) {
	var song dto.AddSongRequest

	if err := c.BindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	id, err := h.services.AddSong(song)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Song added successfully"})
}

// @Summary Search songs
// @Tags songs
// @Description Search songs by query string
// @ID searchSongs
// @Accept  json
// @Produce  json
// @Param  query query string true "Search query"
// @Success 200 {array} model.Song
// @Failure 400 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/v1/songs/search/{query} [get]
func (h *Handler) SearchSongs(c *gin.Context) {
	query := c.Query("query")

	if query == "" {
		h.newErrorResponse(c, http.StatusBadRequest, "Missing query parameter")
		return
	}

	songs, err := h.services.SearchSongs(query)
	if err != nil {
		h.newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, songs)
}
