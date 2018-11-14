package main

import (
	"github.com/anhthii/go-echo/db"
	"github.com/anhthii/go-echo/db/models"

	"github.com/anhthii/go-echo/handlers/media"
	playlistHandlers "github.com/anhthii/go-echo/handlers/playlist"
	userHandlers "github.com/anhthii/go-echo/handlers/user"

	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	router := gin.Default()
	db.Init()
	db.GetDB().AutoMigrate(&models.User{}, &models.Playlist{}, &models.Song{}, &models.Artist{})
	defer db.Close()

	api := router.Group("/api")
	{
		mediaRoutes := api.Group("/media")
		{
			mediaRoutes.GET("/song", media.GetSong)
			mediaRoutes.GET("/albums/default", media.GetDefaultAlbums)
			mediaRoutes.GET("/albums", media.GetAlbums)
			mediaRoutes.GET("/artists/default", media.GetDefaultArtists)
			mediaRoutes.GET("/artists", media.GetArtists)
			mediaRoutes.GET("/artist/:name/:type", media.GetSingleArtist)
			mediaRoutes.GET("/chart/:popType", media.GetChart)
			mediaRoutes.GET("/suggested-song", media.GetSuggestedSongs)
			mediaRoutes.GET("/album_playlist", media.GetAlbumPlaylist)
			mediaRoutes.GET("/top100/:typeID", media.GetTop100)
		}

		userRoutes := api.Group("/user")
		{
			userRoutes.POST("/signup", userHandlers.CreateNewUser)
			userRoutes.POST("/login", userHandlers.Login)

		}

		playlistRoutes := api.Group("/playlist")
		{
			// get all playlists of a user
			playlistRoutes.GET("/:username", playlistHandlers.GetPlaylists)

			// get a specific playlist with title
			// not implemented yet
			// playlistRoutes.GET("/:username/:title", getPlaylist)
			playlistRoutes.POST("/:username", playlistHandlers.CreatePlaylist)

			// // delete a playlist
			playlistRoutes.DELETE("/:username/:playlistTitle", playlistHandlers.DeletePlaylist)

			// add a song to a playlist
			playlistRoutes.PUT("/:username/:playlistTitle", playlistHandlers.AddSongToPlaylist)

			// delete a song from a playlist
			// playlistRoutes.DELETE("/:username/:playlistTitle/:songId", deleteSongFromPlaylist)
		}

	}

	router.Run(":3000")
}
