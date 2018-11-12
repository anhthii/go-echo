package main

import (
	"net/http"

	"github.com/anhthii/go-echo/db"
	"github.com/anhthii/go-echo/db/models"

	"github.com/anhthii/go-echo/handlers/media"
	userHandlers "github.com/anhthii/go-echo/handlers/user"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func main() {
	router := gin.Default()
	db.Init()
	db.Tables(&models.User{})
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

		// userRoutes := api.Group("/user")
		// {
		// 	userRoutes.POST("/signup", user.CreateNewUser)
		// }

		playlist := api.Group("/playlist")
		{
			playlist.GET("/username", func(c *gin.Context) {

				c.JSON(http.StatusOK, gin.H{
					"message": "playlist",
				})
			})
		}

		user := api.Group("/user")
		{
			user.POST("/signup", userHandlers.CreateNewUser)
			user.POST("/login", userHandlers.Login)
		}
	}

	// router.NoRoute(func(c *gin.Context) {
	// 	c.File("./public/index.html")
	// })

	router.Run(":3000")
}
