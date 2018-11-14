package media

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/anhthii/go-echo/utils"
	"github.com/gin-gonic/gin"
)

// SearchSong search a song from the provided term
func SearchSong(c *gin.Context) {
	term := c.Query("term")
	URL := fmt.Sprintf(`https://ac.global.mp3.zing.vn/complete/desktop?type=artist,album,video,song&num=3&query=%s`, strings.Replace(term, " ", "+", -1))
	result, err := utils.GetMapDataFromHTTPGet(URL)
	data := result["data"].([]interface{})

	newResult := make(map[string]interface{})
	newResult["top"] = data[0].(map[string]interface{})["top"]
	newResult["result"] = true
	newResult["data"] = make(map[string]interface{})
	newResult["data"].(map[string]interface{})["song"] = data[1].(map[string]interface{})["song"]
	newResult["data"].(map[string]interface{})["artist"] = data[2].(map[string]interface{})["artist"]
	newResult["data"].(map[string]interface{})["album"] = data[3].(map[string]interface{})["album"]
	if err != nil {
		utils.InternalErrorJSON(c, err)
	}
	c.JSON(http.StatusOK, newResult)
}
