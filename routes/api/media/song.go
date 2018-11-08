package media

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/anhthii/zing.mp3/lrc_parser"
	"github.com/anhthii/zing.mp3/utils"
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
	resourceURL := fmt.Sprintf("https://mp3.zing.vn/xhr/media/get-source?type=audio&%s", key)
	result, err := utils.GetMapDataFromHTTPGet(resourceURL)
	if err != nil {
		utils.InternalErrorJSON(c, err)
	}

	data := result["data"].(map[string]interface{})
	lrcFileURL := data["lyric"].(string)
	lrcString, err := utils.GetStringDataFromHTTPGet(lrcFileURL)
	if err != nil {
		utils.InternalErrorJSON(c, err)
	}

	data["lyric"] = lrc_parser.Parse(lrcString)["scripts"].([]lrc_parser.Lyric)
	c.JSON(http.StatusOK, data)
}
