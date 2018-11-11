package main

import (
	"net/http"

	"github.com/anhthii/go-echo/routes/api/media"
	"github.com/anhthii/go-echo/routes/api/user"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// router.Use(static.Serve("/", static.LocalFile("./public", true)))

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
			userRoutes.POST("/signup", user.CreateNewUser)
		}

		playlist := api.Group("/playlist")
		{
			playlist.GET("/username", func(c *gin.Context) {

				c.JSON(http.StatusOK, gin.H{
					"message": "playlist",
				})
			})
		}

		// user := api.Group("/user")
		// {
		// 	user.GET("/user1", func(c *gin.Context) {
		// 		c.JSON(http.StatusOK, gin.H{
		// 			"message": "user",
		// 		})
		// 	})
		// }
	}

	// router.NoRoute(func(c *gin.Context) {
	// 	c.File("./public/index.html")
	// })

	router.Run(":3000")
}
