package media

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/anhthii/go-echo/libs/scraper"
	"github.com/gin-gonic/gin"
)

// GetDefaultArtists return a list of albums
func GetDefaultArtists(c *gin.Context) {
	URL := "https://mp3.zing.vn/the-loai-nghe-si"
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
			return node.Find(".row.fn-list").Eq(index).Find(".artist-item")
		})
		_sc.SetNoun("artists").Extract("img", scraper.Extractor{Attr: "src", To: "thumb"})
		_sc.ExtractMany("a.txt-primary", []scraper.Extractor{
			scraper.Extractor{Attr: "href", To: "link"},
			scraper.Extractor{Attr: "text", To: "name"},
		})
		_result := _sc.GetResult()
		result["origins"].([]map[string]interface{})[index]["artists"] = _result["artists"]
		result["result"] = true
	}
	c.JSON(http.StatusOK, result)
}
