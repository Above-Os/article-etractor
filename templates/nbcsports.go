package templates

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func (t *Template) NBCSportsScrapContent(document *goquery.Document) string {

	contents := ""

	document.Find("div.ArticlePage-articleBody").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}

type NBCSportsMetaData struct {
	Context       string    `json:"@context"`
	Type          string    `json:"@type"`
	URL           string    `json:"url"`
	DateCreated   time.Time `json:"dateCreated"`
	DateModified  time.Time `json:"dateModified"`
	DatePublished time.Time `json:"datePublished"`
	Description   string    `json:"description"`
	Identifier    string    `json:"identifier"`
	Image         []struct {
		Context         string `json:"@context"`
		Type            string `json:"@type"`
		CopyrightNotice string `json:"copyrightNotice"`
		Credit          string `json:"credit"`
		Height          int    `json:"height"`
		URL             string `json:"url"`
		Width           int    `json:"width"`
	} `json:"image"`
	Speakable struct {
		Type        string   `json:"@type"`
		CSSSelector []string `json:"cssSelector"`
	} `json:"speakable"`
	MainEntityOfPage struct {
		Type string `json:"@type"`
		ID   string `json:"@id"`
	} `json:"mainEntityOfPage"`
	Author []struct {
		Context     string `json:"@context"`
		Type        string `json:"@type"`
		Description string `json:"description"`
		Email       string `json:"email"`
		Name        string `json:"name"`
		URL         string `json:"url"`
	} `json:"author"`
	Publisher struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		Logo struct {
			Type   string `json:"@type"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"logo"`
		URL string `json:"url"`
	} `json:"publisher"`
	Name     string `json:"name"`
	Headline string `json:"headline"`
}

func (t *Template) NBCSPortScrapMetaData(document *goquery.Document) (string, string) {

	author := ""
	published_at := ""
	scriptSelectorFirst := "head > script[type=\"application/ld+json\"]"
	scriptSelectorSecond := "body > script[type=\"application/ld+json\"]"
	scriptSelectorThird := "script[type=\"application/ld+json\"]"

	scriptSelectorList := make([]string, 100)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorFirst)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorSecond)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorThird)

	for _, scriptSelector := range scriptSelectorList {

		document.Find(scriptSelector).Each(func(i int, s *goquery.Selection) {
			if author != "" {
				return
			}
			scriptContent := strings.TrimSpace(s.Text())
			var firstTypeMetaData NBCSportsMetaData
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)
				return
			}

			for _, currentAuthor := range firstTypeMetaData.Author {
				if len(currentAuthor.Name) != 0 {
					if len(author) != 0 {
						author = author + " & " + currentAuthor.Name
					} else {
						author = currentAuthor.Name
					}

				}
			}
		})
		if author != "" {
			break
		}
	}
	log.Printf("author last: %s", author)
	return author, published_at
}

func (t *Template) NBCSportsPublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

	var publishedAt int64 = 0

	scriptSelectorFirst := "head > script[type=\"application/ld+json\"]"
	scriptSelectorSecond := "body > script[type=\"application/ld+json\"]"
	scriptSelectorThird := "script[type=\"application/ld+json\"]"

	scriptSelectorList := make([]string, 100)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorFirst)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorSecond)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorThird)

	for _, scriptSelector := range scriptSelectorList {

		document.Find(scriptSelector).Each(func(i int, s *goquery.Selection) {
			if publishedAt != 0 {
				return
			}
			scriptContent := strings.TrimSpace(s.Text())
			var firstTypeMetaData NBCSportsMetaData
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)
				return

			}
			publishedAt = firstTypeMetaData.DateCreated.Unix()
		})

	}
	return publishedAt
}
