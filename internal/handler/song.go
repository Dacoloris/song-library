package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	"net/http"
	"song-library/internal/models"
	"song-library/internal/service"
	"strconv"
)

type SongHandler struct {
	service service.SongService
	log     *zap.Logger
	apiUrl  string
}

func NewSongHandler(r *gin.Engine, service service.SongService, log *zap.Logger, apiUrl string) *SongHandler {
	songHandler := &SongHandler{
		service: service,
		log:     log,
		apiUrl:  apiUrl,
	}

	r.GET("/songs", songHandler.GetAllSongs)
	r.GET("/songs/:id/lyrics", songHandler.GetSongLyrics)
	r.POST("/songs", songHandler.CreateSong)
	r.PUT("/songs/:id", songHandler.UpdateSong)
	r.DELETE("/songs/:id", songHandler.DeleteSong)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return songHandler
}

// GetAllSongs godoc
// @Summary Get all songs
// @Description Get a list of all songs
// @Tags songs
// @Produce json
// @Success 200 {array} models.Song
// @Failure 500 {object} models.ErrorResponse
// @Router /songs [get]
func (h *SongHandler) GetAllSongs(c *gin.Context) {
	group := c.Query("group")
	title := c.Query("song")

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	h.log.Info("Received request to get all songs", zap.String("group", group), zap.String("song", title))

	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}

	songs, err := h.service.GetAllSongs(group, title, page, limit)
	if err != nil {
		h.log.Error("couldn't get songs", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, songs)
}

// GetSongLyrics godoc
// @Summary Get lyrics of a song
// @Description Get the lyrics of a song by its ID
// @Tags songs
// @Produce json
// @Param id path string true "Song ID"
// @Success 200 {string} string "Lyrics of the song"
// @Failure 404 {object} models.ErrorResponse
// @Router /songs/{id}/lyrics [get]
func (h *SongHandler) GetSongLyrics(c *gin.Context) {
	id := c.Param("id")
	page := c.Query("page")
	limit := c.Query("limit")

	songID, err := strconv.ParseUint(id, 10, 32)

	h.log.Info("Received request to get song lyrics", zap.String("id", id), zap.String("page", page), zap.String("limit", limit))

	lyrics, err := h.service.GetSongLyrics(uint(songID), page, limit)
	if err != nil {
		h.log.Error("Error retrieving song lyrics", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Lyrics not found"})
		return
	}

	h.log.Info("Successfully retrieved song lyrics", zap.String("id", id), zap.Any("lyrics", lyrics))
	c.JSON(http.StatusOK, lyrics)
}

// CreateSong godoc
// @Summary Create a new song
// @Description Create a new song with the provided details
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.Song true "Song data"
// @Success 201 {object} models.Song
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /songs [post]
func (h *SongHandler) CreateSong(c *gin.Context) {
	h.log.Info("received request to create a new song")

	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		h.log.Error("error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	params := "?group=" + song.Group + "&song=" + song.Title
	resp, err := http.Get(h.apiUrl + params)
	if err != nil {
		h.log.Error("failed to call external API", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to call external API"})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		h.log.Error("external API returned an error", zap.Error(err))
		c.JSON(resp.StatusCode, models.ErrorResponse{Error: "External API returned an error"})
		return
	}

	var songDetail models.SongDetail
	if err = json.NewDecoder(resp.Body).Decode(&songDetail); err != nil {
		h.log.Error("failed to decode external API response", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to decode external API response"})
		return
	}

	song.ReleaseDate = songDetail.ReleaseDate
	song.Link = songDetail.Link
	song.Text = songDetail.Text

	if err = h.service.CreateSong(song); err != nil {
		h.log.Error("failed to create song", zap.Error(err))
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	h.log.Info("Successfully created song", zap.Any("song", song))
	c.JSON(http.StatusCreated, song)
}

// UpdateSong godoc
// @Summary Update a song
// @Description Update a song by ID
// @Tags songs
// @Accept json
// @Produce json
// @Param id path string true "Song ID"
// @Param song body models.Song true "Song data"
// @Success 200 {object} models.Song
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /songs/{id} [put]
func (h *SongHandler) UpdateSong(c *gin.Context) {
	var song models.Song
	id := c.Param("id")

	h.log.Info("Received request to update song", zap.String("id", id))

	if err := c.ShouldBindJSON(&song); err != nil {
		h.log.Error("error binding JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	songID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		h.log.Error("wrong song id", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	song.ID = uint(songID)
	if err = h.service.UpdateSong(song); err != nil {
		h.log.Error("failed to update song", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Successfully updated song", zap.String("id", id))
	c.JSON(http.StatusOK, song)
}

// DeleteSong godoc
// @Summary Delete a song
// @Description Delete a song by ID
// @Tags songs
// @Param id path string true "Song ID"
// @Success 204
// @Failure 500 {object} models.ErrorResponse
// @Router /songs/{id} [delete]
func (h *SongHandler) DeleteSong(c *gin.Context) {
	id := c.Param("id")
	h.log.Info("Received request to delete song", zap.String("id", id))

	songID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		h.log.Error("wrong song id", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = h.service.DeleteSong(uint(songID)); err != nil {
		h.log.Error("failed to delete song", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.log.Info("Successfully deleted song", zap.String("id", id))
	c.JSON(http.StatusNoContent, nil)
}
