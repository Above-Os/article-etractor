package templates

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type SkySportsMetaData struct {
	Context             string `json:"@context"`
	Type                string `json:"@type"`
	AlternativeHeadline string `json:"alternativeHeadline"`
	ArticleBody         string `json:"articleBody"`
	MainEntityOfPage    struct {
		Type string `json:"@type"`
		URL  string `json:"url"`
	} `json:"mainEntityOfPage"`
	WordCount  string `json:"wordCount"`
	InLanguage string `json:"inLanguage"`
	Genre      string `json:"genre"`
	Publisher  struct {
		Type string `json:"@type"`
		ID   string `json:"@id"`
		Name string `json:"name"`
		Logo struct {
			Type   string `json:"@type"`
			ID     string `json:"@id"`
			URL    string `json:"url"`
			Width  string `json:"width"`
			Height string `json:"height"`
		} `json:"logo"`
	} `json:"publisher"`
	Headline        string `json:"headline"`
	Description     string `json:"description"`
	Dateline        string `json:"dateline"`
	CopyrightHolder struct {
		ID string `json:"@id"`
	} `json:"copyrightHolder"`
	Author struct {
		Type string `json:"@type"`
		Name string `json:"name"`
	} `json:"author"`
	DatePublished string `json:"datePublished"`
	DateModified  string `json:"dateModified"`
	DateCreated   string `json:"dateCreated"`
	URL           string `json:"url"`
}

func (t *Template) SkySportsScrapMetaData(document *goquery.Document) (string, string) {

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
			var firstTypeMetaData SkySportsMetaData
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)
				return

			}
			author = firstTypeMetaData.Author.Name
		})
		if author != "" {
			break
		}
	}
	log.Printf("author last: %s", author)
	return author, published_at
}

func (t *Template) SkySportsPublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

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
			var firstTypeMetaData SkySportsMetaData
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)
				return

			}
			fmt.Println(firstTypeMetaData.DatePublished)
			parsedPublisedAt, publishedAtparseErr := parseToTimestampSkysports(firstTypeMetaData.DatePublished)
			if publishedAtparseErr != nil {
				log.Printf("parse skysport time error %v", publishedAtparseErr)
				return
			}
			publishedAt = parsedPublisedAt
		})

	}
	return publishedAt
}

func parseToTimestampSkysports(timeStr string) (int64, error) {
	const layout = "2006-01-02T15:04:05-0700"
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}

func (t *Template) SkySportsScrapContent(document *goquery.Document) string {

	contents := ""
	document.Find("div.sdc-article-related-stories,div.sdc-article-factbox,div.sdc-article-strapline,div.site-share-wrapper,div[data-format=floated-mpu]").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	document.Find("div.sdc-article-image,div[data-type=article]").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents

}
