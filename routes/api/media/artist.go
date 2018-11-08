package media

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/anhthii/zing.mp3/scraper"
	"github.com/gin-gonic/gin"
)

// GetSingleArtist return data of a single artist
func GetSingleArtist(c *gin.Context) {
	name := c.Param("name")
	_type := c.Param("type")
	page := c.Query("page")

	var _type_ string
	switch _type {
	case "songs":
		_type_ = "bai-hat"
	case "albums":
		_type_ = "album"
	case "biography":
		_type_ = "tieu-su"
	}

	var URL string
	if page == "" {
		URL = fmt.Sprintf(`http://mp3.zing.vn/nghe-si/%s/%s`, name, _type_)
	} else {
		URL = fmt.Sprintf(`http://mp3.zing.vn/nghe-si/%s/%s?&page=%s`, name, _type_, page)
	}

	sc := scraper.NewZingMp3Scraper(nil)
	sc.
		Scrape(URL).
		Extract(".box-info-artist img", scraper.Extractor{Attr: "src", To: "avatar"}).
		Extract(".container > img", scraper.Extractor{Attr: "src", To: "cover"}).
		Extract(".info-summary > h1", scraper.Extractor{Attr: "text", To: "artistName"})

	switch _type {
	case "biography":
	case "songs":
		sc.
			FindList(".group.fn-song").
			SetNoun("songs").
			Extract("._trackLink span", scraper.Extractor{Attr: "text", To: "artist_text"})
		sc.
			ExtractMany("._trackLink", []scraper.Extractor{
				scraper.Extractor{Attr: "text", To: "title", Filter: func(s string) string {
					regex := regexp.MustCompile(`(?m)\s*-\s*.+`)
					return strings.TrimSpace(regex.ReplaceAllString(s, ""))
				}},
				scraper.Extractor{Attr: "href", To: "id"},
				scraper.Extractor{Attr: "href", To: "alias"},
			}).
			Paginate()

	case "albums":
		sc.
			FindList(".album-item").
			SetNoun("albums").
			Extract("img", scraper.Extractor{Attr: "src", To: "thumb"}).
			Extract(".title-item .txt-primary", scraper.Extractor{Attr: "text", To: "title", Filter: func(s string) string {
				return strings.TrimSpace(s)
			}})
		sc.
			ExtractMany(".thumb._trackLink", []scraper.Extractor{
				scraper.Extractor{Attr: "href", To: "id"},
				scraper.Extractor{Attr: "href", To: "alias"},
			}).
			ExtractArtists(".fn-artist  a").
			Paginate()
	}

	result := sc.GetResult()
	c.JSON(http.StatusOK, result)
}
