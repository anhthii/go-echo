package media

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/anhthii/go-echo/utils"
	"github.com/gin-gonic/gin"
)

// GetSuggestedSongs return suggested data for a playing song
func GetSuggestedSongs(c *gin.Context) {
	songID := c.Query("songId")
	artistID := c.Query("artistId")

	URL := fmt.Sprintf(`https://mp3.zing.vn/xhr/recommend?target=%%23block-recommend&count=20&start=0&artistid=%s&type=audio&id=%s`, artistID, songID)
	response, err := http.Get(URL)

	if err != nil {
		utils.InternalErrorJSON(c, err)
	}
	defer response.Body.Close()
	var result map[string]interface{}
	bytes, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(bytes, &result)
	c.JSON(http.StatusOK, result)
}
