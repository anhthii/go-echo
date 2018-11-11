package media

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/anhthii/go-echo/utils"
	"github.com/gin-gonic/gin"
)

// GetTop100 return the top 100 songs in the home page
func GetTop100(c *gin.Context) {
	// const (
	// 	PopID  = "ZWZB96AB"
	// 	KPopID = "ZWZB96DC"
	// 	VPopID = "ZWZB969E"
	// )
	typeID := c.Param("typeID")

	pageNo := c.DefaultQuery("page", "1")
	num, _ := strconv.Atoi(pageNo)
	start := strconv.Itoa((num - 1) * 20)
	URL := fmt.Sprintf("https://mp3.zing.vn/xhr/media/get-list?op=top100&start=%s&length=20&id=%s", start, typeID)
	result, err := utils.GetMapDataFromHTTPGet(URL)
	if err != nil {
		utils.InternalErrorJSON(c, err)
	}
	data := result["data"].(map[string]interface{})
	data["items"] = data["items"].([]interface{})[:20]
	c.JSON(http.StatusOK, result)
}
