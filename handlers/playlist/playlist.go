package playlist

import (
	"net/http"

	"github.com/anhthii/go-echo/db/models"
	Model "github.com/anhthii/go-echo/db/models"
	. "github.com/anhthii/go-echo/utils"
	"github.com/gin-gonic/gin"
	validator "gopkg.in/go-playground/validator.v8"
)

type PlaylistBody struct {
	Title string `json:"title" binding:"required,min=3,max=50"`
}

func CreatePlaylist(c *gin.Context) {
	var json PlaylistBody
	username := c.Param("username")

	if err := c.ShouldBindJSON(&json); err != nil {
		validateResponse := Validate(err.(validator.ValidationErrors))
		c.JSON(http.StatusBadRequest, validateResponse)
		return
	}

	playlistData, errorStr := Model.CreatePlaylist(username, json.Title)
	if errorStr != "" {
		c.JSON(http.StatusOK, errorStr)
	}
	c.JSON(http.StatusOK, playlistData)
}

func AddSongToPlaylist(c *gin.Context) {
	username := c.Param("username")
	playlistTitle := c.Param("playlistTitle")

	var song models.Song
	if err := c.ShouldBindJSON(&song); err != nil {
		c.JSON(http.StatusBadRequest, "Invalid song data")
		return
	}

	ok, errorStr := song.AddToPlaylist(username, playlistTitle)
	if !ok {
		c.JSON(400, errorStr)
	}
	c.JSON(http.StatusOK, "ok")
}

func GetPlaylists(c *gin.Context) {
	username := c.Param("username")
	var playlists []models.Playlist
	playlists = models.GetPlaylists(username)
	c.JSON(http.StatusOK, playlists)
}

func DeletePlaylist(c *gin.Context) {
	username := c.Param("username")
	playlistTitle := c.Param("playlistTitle")
	ok, errorStr := models.DeletePlaylist(username, playlistTitle)
	if !ok {
		c.JSON(400, errorStr)
	}
	c.JSON(http.StatusOK, "ok")
}
