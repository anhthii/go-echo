package scraper

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

// Extractor specify what attribute a user need to extract from a html tag
type Extractor struct {
	Attr, To string              // `to` specify the corresponding name for the key in the `map` result
	Filter   func(string) string // filtering the extracted value
}

// SubFieldAsList ...
type SubFieldAsList struct {
	selector, fieldname string
}

// Scraper have all utilities for scraping a website easily
type Scraper struct {
	Doc                  *goquery.Selection
	Result               map[string]interface{}
	Noun                 string             // store name of the scraped list in the Result
	listElements         *goquery.Selection // store the list Nodes got from `FindList` method
	isListSelected       bool
	selectors            map[string][]Extractor         // every time Extract or ExtractMultiple method is called, a new selector is added to this
	resultSubFieldAsList map[SubFieldAsList][]Extractor // store extractors to scrape data for a field which have a list structure
	// ig: result {
	// 	field: {
	// 		subField: ...
	// 		subfieldAsList: [
	// 			...
	// 		]
	// 	}
}

// Init initialize a scraper
func Init() *Scraper {
	scraper := &Scraper{}
	scraper.Result = make(map[string]interface{})
	scraper.selectors = make(map[string][]Extractor)
	scraper.isListSelected = false
	return scraper
}

// Scrape specify URL used for extract HTML document from a website
func (scraper *Scraper) Scrape(url string) *Scraper {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	scraper.Doc = document.Find("body")
	return scraper
}

func (scraper *Scraper) SetRoot(selector string) *Scraper {
	scraper.Doc = scraper.Doc.Find(selector)
	return scraper
}

// FindList select a list of Nodes through a selector, this list then is used for making scraping job a lot easier
func (scraper *Scraper) FindList(selector string) *Scraper {
	scraper.listElements = scraper.Doc.Find(selector)
	scraper.isListSelected = true
	return scraper
}

func (scraper *Scraper) SetList(fn func(*goquery.Selection) *goquery.Selection) *Scraper {
	scraper.listElements = fn(scraper.Doc)
	scraper.isListSelected = true
	return scraper
}

// SetNoun specify the name of the list in the `map` result
// Example SetNoun("albums")
// {
// 	"albums": [
// 		{ "cover" : "...", "id": "... "}
// 	]
// }
func (scraper *Scraper) SetNoun(noun string) *Scraper {
	scraper.Noun = noun
	scraper.Result[noun] = make([]map[string]interface{}, 0)
	return scraper
}

func extract(extractor Extractor, element *goquery.Selection) string {
	var attr string
	if extractor.Attr == "text" {
		attr = element.Text()
	} else {
		_attr, _ := element.Attr(extractor.Attr)
		attr = _attr
	}

	// check if filter exists
	if extractor.Filter != nil {
		attr = extractor.Filter(attr)
	}
	return attr
}

// Extract attribute from a tag specified through a selector
func (scraper *Scraper) Extract(selector string, ext Extractor) *Scraper {
	if !scraper.isListSelected {
		scraper.Result[ext.To] = extract(ext, scraper.Doc.Find(selector))
		return scraper
	}
	_, ok := scraper.selectors[selector]
	if !ok {
		scraper.selectors[selector] = make([]Extractor, 0)
	}
	scraper.selectors[selector] = append(scraper.selectors[selector], ext)
	return scraper
}

// ExtractMultiple used for extracting multiple attributes from a single HTML tag
func (scraper *Scraper) ExtractMultiple(selector string, exts []Extractor) *Scraper {
	for _, ext := range exts {
		scraper.Extract(selector, ext)
	}
	return scraper
}

// ExtractMultipleIntoList scrape data for a field which have a list structure
func (scraper *Scraper) ExtractMultipleIntoList(selector string, fieldname string, exts []Extractor) *Scraper {
	if scraper.resultSubFieldAsList == nil {
		scraper.resultSubFieldAsList = make(map[SubFieldAsList][]Extractor)
	}
	subField := SubFieldAsList{selector, fieldname}
	_, ok := scraper.resultSubFieldAsList[subField]
	if !ok {
		scraper.resultSubFieldAsList[subField] = make([]Extractor, len(exts))
	}
	scraper.resultSubFieldAsList[subField] = exts
	return scraper
}

// GetResult return the scraped data
func (scraper *Scraper) GetResult() map[string]interface{} {
	if scraper.listElements == nil {
		return scraper.Result
	}

	scraper.listElements.Each(func(index int, selection *goquery.Selection) {
		item := make(map[string]interface{})
		for selector, extractors := range scraper.selectors {
			element := selection.Find(selector)
			for _, extractor := range extractors {
				item[extractor.To] = extract(extractor, element)
			}

		}

		for subFieldAsList, extractors := range scraper.resultSubFieldAsList {
			selector := subFieldAsList.selector
			fieldname := subFieldAsList.fieldname
			item[fieldname] = make([]map[string]string, 0)
			selection.Find(selector).Each(func(idx int, el *goquery.Selection) {
				_item := make(map[string]string)
				for _, extractor := range extractors {
					_item[extractor.To] = extract(extractor, el)
				}
				item[fieldname] = append(item[fieldname].([]map[string]string), _item)
			})

		}
		scraper.Result[scraper.Noun] = append(scraper.Result[scraper.Noun].([]map[string]interface{}), item)
	})
	return scraper.Result
}
