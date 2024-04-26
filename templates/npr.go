package templates

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type NprMetadata struct {
	Type      string `json:"@type"`
	Publisher struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		Logo struct {
			Type string `json:"@type"`
			URL  string `json:"url"`
		} `json:"logo"`
	} `json:"publisher"`
	Headline         string `json:"headline"`
	MainEntityOfPage struct {
		Type string `json:"@type"`
		ID   string `json:"@id"`
	} `json:"mainEntityOfPage"`
	DatePublished time.Time `json:"datePublished"`
	DateModified  time.Time `json:"dateModified"`
	Author        struct {
		Type string   `json:"@type"`
		Name []string `json:"name"`
	} `json:"author"`
	Description string `json:"description"`
	SubjectOf   []struct {
		Type         string `json:"@type"`
		Name         string `json:"name"`
		Description  string `json:"description"`
		ThumbnailURL string `json:"thumbnailUrl"`
		UploadDate   string `json:"uploadDate"`
		EmbedURL     string `json:"embedUrl"`
	} `json:"subjectOf"`
	Context string `json:"@context"`
}

func (t *Template) NprScrapMetaData(document *goquery.Document) (string, string) {

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
			var firstTypeMetaData NprMetadata
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)
				return
			}
			for _, currentName := range firstTypeMetaData.Author.Name {
				if len(author) != 0 {
					author = " & "
				}
				author = author + currentName
			}
		})
		if author != "" {
			break
		}
	}
	log.Printf("author last: %s", author)
	return author, published_at
}

func (t *Template) NprPublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

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
			var firstTypeMetaData NprMetadata
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)
				return

			}
			fmt.Println(firstTypeMetaData.DatePublished)
			publishedAt = firstTypeMetaData.DatePublished.Unix()
		})

	}
	return publishedAt
}

func (t *Template) NprScrapContent(document *goquery.Document) string {
	contents := ""

	document.Find("div.internallink,div.enlarge-options,div.enlarge_measure,div.enlarge_html").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)

	})
	document.Find("div#storytext").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
