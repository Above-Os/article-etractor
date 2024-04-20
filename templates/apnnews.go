package templates

import (
	"encoding/json"
	"log"

	//"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ApNewsMetaData []struct {
	Context       string    `json:"@context"`
	Type          string    `json:"@type"`
	URL           string    `json:"url"`
	DateModified  time.Time `json:"dateModified"`
	DatePublished time.Time `json:"datePublished"`
	Description   string    `json:"description"`
	Image         []struct {
		Context      string `json:"@context"`
		Type         string `json:"@type"`
		Height       int    `json:"height"`
		ThumbnailURL string `json:"thumbnailUrl"`
		URL          string `json:"url"`
		Width        int    `json:"width"`
	} `json:"image"`
	MainEntityOfPage struct {
		Type string `json:"@type"`
		ID   string `json:"@id"`
	} `json:"mainEntityOfPage"`
	Author []struct {
		Context string `json:"@context"`
		Type    string `json:"@type"`
		Name    string `json:"name"`
	} `json:"author"`
	Publisher struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		Logo struct {
			Type string `json:"@type"`
			URL  string `json:"url"`
		} `json:"logo"`
	} `json:"publisher"`
	ArticleSection []string `json:"articleSection"`
	Keywords       []string `json:"keywords"`
	ThumbnailURL   string   `json:"thumbnailUrl"`
	Name           string   `json:"name"`
	Headline       string   `json:"headline"`
}

func (t *Template) ApNewsCommonGetPublishedAtTimestamp(document *goquery.Document) int64 {
	var publishedAt int64 = 0
	publishedAt = t.CommonGetPublishedAtTimestampSingleJson(document)
	if publishedAt == 0 {
		publishedAt = t.ApNewsPublishedAtTimeFromScriptMetadata(document)
	}
	return publishedAt

}

func (t *Template) ApNewsPublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {
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
			var firstTypeMetaData ApNewsMetaData
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert ApNewsMetaData unmarshalError %v", unmarshalErr)
				return

			}
			firstElement := firstTypeMetaData[0]
			publishedAt = firstElement.DateModified.Unix()

			// var jsonMap map[string]interface{}

			// logger.Info("script content %s  author length %d",scriptContent, len(currentDWMetadata.Author))
		})
		if publishedAt != 0 {
			break
		}
	}
	return publishedAt

}

func (t *Template) ApNewsAuthorExtractFromScriptMetadata(document *goquery.Document) (string, string) {

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
			var firstTypeMetaData ApNewsMetaData
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert ApNewsMetaData unmarshalError %v", unmarshalErr)
				return

			}
			firstElement := firstTypeMetaData[0]
			for currentIndex, currentAuthor := range firstElement.Author {

				if currentIndex != 0 {
					author += " & "
				}
				log.Printf("author Name: ", currentAuthor.Name)
				author += currentAuthor.Name
			}

			// var jsonMap map[string]interface{}

			// logger.Info("script content %s  author length %d",scriptContent, len(currentDWMetadata.Author))
		})
		if author != "" {
			break
		}
	}
	log.Printf("author last: %s", author)
	return author, published_at
}

func (t *Template) ApnNewsScrapMetaData(document *goquery.Document) (string, string) {
	author := ""
	published_at := ""
	author, published_at = t.AuthorExtractFromScriptMetadata(document)
	if author == "" {
		author, published_at = t.AuthorExtractFromScriptMetadata(document)
	}
	if author == "" {
		author, published_at = t.ApNewsAuthorExtractFromScriptMetadata(document)
	}
	if author == "" {
		author = "AP News"
	}

	return author, published_at
}
