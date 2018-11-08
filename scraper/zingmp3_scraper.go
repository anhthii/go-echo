package scraper

import (
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ZingMp3Scraper struct {
	*Scraper
}

func NewZingMp3Scraper(r io.Reader) *ZingMp3Scraper {
	if r != nil {
		doc, err := goquery.NewDocumentFromReader(r)
		if err != nil {
			log.Fatal(err)
		}
		scraper := Init()
		scraper.Doc = doc.Find("body")
		return &ZingMp3Scraper{scraper}
	}
	return &ZingMp3Scraper{Init()}
}

func NewZingMp3ScraperFromNode(selection *goquery.Selection) *ZingMp3Scraper {
	scraper := Init()
	scraper.Doc = selection
	return &ZingMp3Scraper{scraper}
}

func (sc *ZingMp3Scraper) ExtractArtists(selector string) *ZingMp3Scraper {
	sc.ExtractMultipleIntoList(selector, "artists", []Extractor{
		Extractor{Attr: "href", To: "alias", Filter: func(s string) string {
			regex := regexp.MustCompile(`(\/nghe-si\/|http:\/\/mp3\.zing.vn\/nghe-si\/)`)
			return regex.ReplaceAllString(s, "")
		}},
		Extractor{Attr: "text", To: "name"},
	})
	return sc
}

func (sc *ZingMp3Scraper) ExtractMany(selector string, exts []Extractor) *ZingMp3Scraper {
	for i, _ := range exts {
		if exts[i].To == "id" {
			exts[i].Filter = func(s string) string {
				r := regexp.MustCompile(`(\w+).html`)
				return r.FindStringSubmatch(s)[1]
			}
		} else if exts[i].To == "alias" {
			exts[i].Filter = func(s string) string {
				r := regexp.MustCompile(`\/([0-9A-Za-z_-]+)\/\w+\.html`)
				return r.FindStringSubmatch(s)[1]
			}
		}
		if exts[i].Attr == "text" {
			exts[i].Filter = func(s string) string {
				return strings.TrimSpace(s)
			}
		}
	}
	sc.ExtractMultiple(selector, exts)
	return sc
}

func (sc *ZingMp3Scraper) Paginate() *ZingMp3Scraper {
	lastPageLink := sc.Doc.Find(".pagination ul li").Last().Children().Eq(0)
	if lastPageLink != nil {
		href, _ := lastPageLink.Attr("href")
		r := regexp.MustCompile(`\d+$`)
		i, _ := strconv.Atoi(r.FindString(href))
		sc.Result["numberOfPages"] = i
	} else {
		sc.Result["numberOfPages"] = 1
	}
	return sc
}
