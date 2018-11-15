package scraper

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const URL = "https://mp3.zing.vn/the-loai-album/Pop/IWZ9Z097.html"

func TestScrape(t *testing.T) {
	scraper := Init()
	result := scraper.Scrape(URL).
		SetNoun("albums").
		FindList(".row.fn-list .album-item").
		Extract("img", Extractor{Attr: "src", To: "cover"}).
		ExtractMultiple(".fn-name.fn-link", []Extractor{
			Extractor{Attr: "text", To: "title", Filter: func(s string) string {
				return strings.TrimSpace(s)
			}},
			Extractor{Attr: "href", To: "id", Filter: func(s string) string {
				r := regexp.MustCompile(`(\w+).html`)
				return r.FindStringSubmatch(s)[1]
			}},
			Extractor{Attr: "href", To: "alias", Filter: func(s string) string {
				r := regexp.MustCompile(`\/([0-9A-Za-z_-]+)\/\w+\.html$`)
				return r.FindStringSubmatch(s)[1]
			}},
		}).
		ExtractMultipleIntoList(".fn-artist  a", "artists", []Extractor{
			Extractor{Attr: "href", To: "alias", Filter: func(s string) string {
				regex := regexp.MustCompile(`(\/nghe-si\/|http:\/\/mp3\.zing.vn\/nghe-si\/)`)
				return regex.ReplaceAllString(s, "")
			}},
			Extractor{Attr: "text", To: "name"},
		}).
		GetResult()

	albums := result["albums"].([]map[string]interface{})
	assert.Equal(t, len(albums), 20, "albums should have 20 items")
	assert.Equal(t, strings.HasSuffix(albums[0]["cover"].(string), ".jpg"), true, "album should have `cover` field")
	_, haveAlias := albums[0]["alias"]
	_, haveID := albums[0]["id"]
	_, haveTitle := albums[0]["title"]
	assert.Equal(t, haveAlias, true, "album should have alias")
	assert.Equal(t, haveID, true, "album should have id")
	assert.Equal(t, haveTitle, true, "album should have title")
}
