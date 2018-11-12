package media

import (
	"fmt"
	"net/http"

	"github.com/anhthii/go-echo/scraper"
	"github.com/gin-gonic/gin"
)

// GetAlbums return a list of albums
func GetAlbums(c *gin.Context) {
	genre := c.Query("genre")
	id := c.Query("id")
	page := c.DefaultQuery("page", "1")

	URL := fmt.Sprintf(`http://mp3.zing.vn/the-loai-album/%s/%s.html?page=%s}`, genre, id, page)
	sc := scraper.NewZingMp3Scraper(nil)
	sc.Scrape(URL).
		SetNoun("albums").
		FindList(".row.fn-list .album-item").
		Extract("img", scraper.Extractor{Attr: "src", To: "cover"})
	sc.ExtractMany(".fn-name.fn-link", []scraper.Extractor{
		scraper.Extractor{Attr: "text", To: "title"},
		scraper.Extractor{Attr: "href", To: "id"},
		scraper.Extractor{Attr: "href", To: "alias"},
	}).
		ExtractArtists(".fn-artist  a").
		Paginate()
	result := sc.GetResult()
	c.JSON(http.StatusOK, result)
}
