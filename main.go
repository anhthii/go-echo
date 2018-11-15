package main

import (
	"github.com/anhthii/go-echo/db"
	"github.com/anhthii/go-echo/db/models"

	"github.com/anhthii/go-echo/handlers/media"
	playlistHandlers "github.com/anhthii/go-echo/handlers/playlist"
	userHandlers "github.com/anhthii/go-echo/handlers/user"
	authMiddlewares "github.com/anhthii/go-echo/middlewares"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	// Init database
	db.Init()
	db.Tables(&models.User{}, &models.Playlist{}, &models.Song{}, &models.Artist{})
	defer db.Close()

	router := SetupRouter()
	router.Run(":3000")
}

// SetupRouter registers routes
func SetupRouter() *gin.Engine {
	router := gin.Default()
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
			mediaRoutes.GET("/search", media.SearchSong)
		}

		userRoutes := api.Group("/user")
		{
			userRoutes.POST("/signup", userHandlers.CreateNewUser)
			userRoutes.POST("/login", userHandlers.Login)

		}

		playlistRoutes := api.Group("/playlist", authMiddlewares.TokenAuthMiddleware())
		{
			// get all playlists of a user
			playlistRoutes.GET("/:username", authMiddlewares.IsValidUser, playlistHandlers.GetPlaylists)

			// get a specific playlist with title
			// not implemented yet
			// playlistRoutes.GET("/:username/:title", getPlaylist)
			playlistRoutes.POST("/:username", authMiddlewares.IsValidUser, playlistHandlers.CreatePlaylist)

			// // delete a playlist
			playlistRoutes.DELETE("/:username/:playlistTitle", authMiddlewares.IsValidUser, playlistHandlers.DeletePlaylist)

			// add a song to a playlist
			playlistRoutes.PUT("/:username/:playlistTitle", authMiddlewares.IsValidUser, playlistHandlers.AddSongToPlaylist)

			// delete a song from a playlist
			// playlistRoutes.DELETE("/:username/:playlistTitle/:songId", deleteSongFromPlaylist)
		}
	}
	// router.GET("/download/song/:songTitle/:id", handlers.DownloadSong)
	return router
}
