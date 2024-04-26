package templates

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type SkyNewsMetaData struct {
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
	DatePublished time.Time `json:"datePublished"`
	DateModified  time.Time `json:"dateModified"`
	DateCreated   time.Time `json:"dateCreated"`
	Image         struct {
		Type   string `json:"@type"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"image"`
	URL string `json:"url"`
}

func (t *Template) SkyNewsScrapMetaData(document *goquery.Document) (string, string) {

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
			var firstTypeMetaData SkyNewsMetaData
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

func (t *Template) SkyNewsPublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

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
			var firstTypeMetaData SkyNewsMetaData
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)
				return

			}
			publishedAt = firstTypeMetaData.DatePublished.Unix()
		})

	}
	return publishedAt
}

func (t *Template) SkyNewsScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("div.sdc-article-related-stories,div.sdc-site-video,a,span[data-label-text=Advertisement],div.sdc-article-related-stories,div.sdc-article-strapline").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("p").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "Read more:") {
			RemoveNodes(s)
		}
	})
	document.Find("div.sdc-site-component-top__media,div.sdc-article-body").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
