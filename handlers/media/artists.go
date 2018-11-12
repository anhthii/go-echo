package media

import (
	"fmt"
	"net/http"

	"github.com/anhthii/go-echo/scraper"
	"github.com/gin-gonic/gin"
)

// GetArtists return a list of artists
func GetArtists(c *gin.Context) {
	genre := c.Query("genre")
	id := c.Query("id")
	page := c.DefaultQuery("page", "1")

	URL := fmt.Sprintf(`http://mp3.zing.vn/the-loai-nghe-si/%s/%s.html?page=%s}`, genre, id, page)
	sc := scraper.NewZingMp3Scraper(nil)
	sc.Scrape(URL).
		SetNoun("artists").
		FindList(".pone-of-five .item").
		Extract("img", scraper.Extractor{Attr: "src", To: "thumb"}).
		Extract("a.txt-primary", scraper.Extractor{Attr: "text", To: "name"})
	sc.Paginate()
	result := sc.GetResult()
	c.JSON(http.StatusOK, result)
}
