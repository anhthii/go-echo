package media

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/anhthii/go-echo/scraper"
	"github.com/gin-gonic/gin"
)

// GetDefaultAlbums return a list of default albums
func GetDefaultAlbums(c *gin.Context) {
	URL := "http://mp3.zing.vn/the-loai-album.html"
	sc := scraper.NewZingMp3Scraper(nil)
	sc.
		Scrape(URL).
		FindList(".zcontent .title-section").
		SetNoun("origins")
	sc.
		ExtractMany("a", []scraper.Extractor{
			scraper.Extractor{Attr: "href", To: "id"},
			scraper.Extractor{Attr: "text", To: "title"},
		})
	result := sc.GetResult()
	node := sc.Doc.Find(".zcontent")
	for index, _ := range result["origins"].([]map[string]interface{}) {
		_sc := scraper.NewZingMp3ScraperFromNode(node)
		_sc.SetList(func(node *goquery.Selection) *goquery.Selection {
			return node.Find(".row.fn-list").Eq(index).Find(".album-item.fn-item")
		})
		_sc.SetNoun("albums").Extract("img", scraper.Extractor{Attr: "src", To: "cover"})
		_sc.ExtractMany(".title-item a", []scraper.Extractor{
			scraper.Extractor{Attr: "text", To: "title"},
			scraper.Extractor{Attr: "href", To: "id"},
			scraper.Extractor{Attr: "href", To: "alias"},
		}).ExtractArtists(".fn-artist a")
		_result := _sc.GetResult()
		result["origins"].([]map[string]interface{})[index]["albums"] = _result["albums"]
		result["result"] = true
	}
	c.JSON(http.StatusOK, result)
}
