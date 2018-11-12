package media

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/anhthii/go-echo/utils"
	"github.com/gin-gonic/gin"
)

// GetSong return data needed for a song to be played at client site
func GetSong(c *gin.Context) {
	name := c.Query("name")
	id := c.Query("id")

	pageURL := fmt.Sprintf("https://mp3.zing.vn/bai-hat/%s/%s.html", name, id)
	r := regexp.MustCompile(`key=\w+`)
	s, err := utils.GetStringDataFromHTTPGet(pageURL)
	if err != nil {
		utils.InternalErrorJSON(c, err)
	}

	key := r.FindStringSubmatch(s)[0]
	songResourceEndpoint := fmt.Sprintf("media/get-source?type=audio&%s", key)
	c.JSON(http.StatusOK, songResourceEndpoint)
}
