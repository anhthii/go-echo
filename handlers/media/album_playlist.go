package media

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/anhthii/go-echo/utils"

	"github.com/PuerkitoBio/goquery"
	"github.com/anhthii/go-echo/scraper"
	"github.com/gin-gonic/gin"
)

// GetAlbumPlaylist return a playlist of songs for an album
func GetAlbumPlaylist(c *gin.Context) {
	title := c.Query("title")
	id := c.Query("id")
	URL := fmt.Sprintf(`http://mp3.zing.vn/album/%s/%s.html`, title, id)
	rg := regexp.MustCompile(`media\/get-source\?type=album&key=\w+`)

	response, err := http.Get(URL)
	bytes, _ := ioutil.ReadAll(response.Body)
	data := string(bytes)
	album_resource_URL := rg.FindString(data)
	sc := scraper.NewZingMp3Scraper(nil)
	sc.
		Scrape(URL).
		Extract(".info-top-play img", scraper.Extractor{Attr: "src", To: "album_playlist_thumb"}).
		Extract(".ctn-inside > h1", scraper.Extractor{Attr: "text", To: "album_title"}).
		Extract(".info-song-top .inline", scraper.Extractor{Attr: "text", To: "release_year"}).
		Extract("img.thumb-art", scraper.Extractor{Attr: "src", To: "artist_thumb"}).
		Extract(".box-artist img", scraper.Extractor{Attr: "text", To: "artist"}).
		Extract(".artist-info-text > p", scraper.Extractor{Attr: "text", To: "artist_info"})

	result := sc.GetResult()
	url := "https://mp3.zing.vn/xhr/" + album_resource_URL
	_response, err := http.Get(url)

	if err != nil {
		utils.InternalErrorJSON(c, err)
	}
	defer response.Body.Close()
	var songData map[string]interface{}
	_bytes, _ := ioutil.ReadAll(_response.Body)
	json.Unmarshal(_bytes, &songData)
	result["songs"] = songData["data"].(map[string]interface{})["items"]
	// result["songs"].([]map[string]interface{}) = songData["data"].(map[string]interface{})["items"].([]map[string]interface{})
	result["genres"] = make([]string, 0)
	sc.SetRoot(".info-song-top").Doc.Find("a").Each(func(index int, element *goquery.Selection) {
		result["genres"] = append(result["genres"].([]string), element.Text())
	})
	c.JSON(http.StatusOK, result)
}
