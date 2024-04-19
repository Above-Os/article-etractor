package templates

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"github.com/PuerkitoBio/goquery"
	"recommend.common/logger"
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
	Image       struct {
		Type string `json:"@type"`
		URL  string `json:"url"`
	} `json:"image"`
	SubjectOf []struct {
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
			var firstTypeMetaData NprMetadata;
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				logger.Info("convert SkyNewsScrap unmarshalError %v",unmarshalErr) 
				return
			}
			for _,currentName := range firstTypeMetaData.Author.Name {
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
    logger.Info("author last: %s",author)
	return author, published_at
}

func (t* Template) NprPublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

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
			var firstTypeMetaData NprMetadata;
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				logger.Info("convert SkyNewsScrap unmarshalError %v",unmarshalErr) 
				return

			}
			fmt.Println(firstTypeMetaData.DatePublished)
			publishedAt = firstTypeMetaData.DatePublished.Unix()
		})

	}
	return publishedAt
}

