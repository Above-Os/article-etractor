package templates

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type BBCNewsMetadataSecond struct {
	Context   string `json:"@context"`
	Type      string `json:"@type"`
	URL       string `json:"url"`
	Publisher struct {
		Type                 string `json:"@type"`
		Name                 string `json:"name"`
		PublishingPrinciples string `json:"publishingPrinciples"`
		Logo                 struct {
			Type string `json:"@type"`
			URL  string `json:"url"`
		} `json:"logo"`
	} `json:"publisher"`
	DatePublished time.Time `json:"datePublished"`
	DateModified  time.Time `json:"dateModified"`
	Description   string    `json:"description"`
	Headline      string    `json:"headline"`
	Image         struct {
		Type   string `json:"@type"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
		URL    string `json:"url"`
	} `json:"image"`
	ThumbnailURL     string `json:"thumbnailUrl"`
	MainEntityOfPage string `json:"mainEntityOfPage"`
	Author           struct {
		Type            string `json:"@type"`
		Name            string `json:"name"`
		NoBylinesPolicy string `json:"noBylinesPolicy"`
		Logo            struct {
			Type string `json:"@type"`
			URL  string `json:"url"`
		} `json:"logo"`
	} `json:"author"`
}

type BBCNewsMetaData struct {
	Context   string `json:"@context"`
	Type      string `json:"@type"`
	URL       string `json:"url"`
	Publisher struct {
		Type                 string `json:"@type"`
		Name                 string `json:"name"`
		PublishingPrinciples string `json:"publishingPrinciples"`
		Logo                 struct {
			Type string `json:"@type"`
			URL  string `json:"url"`
		} `json:"logo"`
	} `json:"publisher"`
	DatePublished time.Time `json:"datePublished"`
	DateModified  time.Time `json:"dateModified"`
	Description   string    `json:"description"`
	Headline      string    `json:"headline"`
	Image         struct {
		Type   string `json:"@type"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
		URL    string `json:"url"`
	} `json:"image"`
	ThumbnailURL     string `json:"thumbnailUrl"`
	MainEntityOfPage string `json:"mainEntityOfPage"`
	Author           []struct {
		Type string `json:"@type"`
		Name string `json:"name"`
	} `json:"author"`
}

func (t *Template) BBCScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("header,div[data-testid=byline],div[data-component=topic-list],div[data-component=links-block],span.visually-hidden,h1#main-heading,div[data-component=byline-block],div[data-component=timestamp-block]").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("div.description").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents = content
	})
	if contents == "" {
		document.Find("article,#main-content").Each(func(i int, s *goquery.Selection) {
			var content string
			content, _ = goquery.OuterHtml(s)
			contents = content
		})
	}
	return contents
}

/**
func (t *Template) BBCScrapMetaData(document *goquery.Document) (string, string) {

	author := ""
	published_at := ""
	author, published_at = t.AuthorExtractFromScriptMetadata(document)
	if author != "" {
		byPrefix := "By "
		exist := strings.HasPrefix(author, byPrefix)
		if exist {
			author = author[len(byPrefix)-1:]
		}
	}

	return author, published_at
}
*/

func (t *Template) BBCNewsScrapMetaData(document *goquery.Document) (string, string) {

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
			var firstTypeMetaData BBCNewsMetaData
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)

			} else {
				for _, currentAuthor := range firstTypeMetaData.Author {
					if len(currentAuthor.Name) != 0 {
						if len(author) != 0 {
							author = author + " & " + currentAuthor.Name
						} else {
							author = currentAuthor.Name
						}
					}

				}
			}
			if len(author) != 0 {
				return
			}
			var secondTypeMetaData BBCNewsMetadataSecond
			unmarshalErr = json.Unmarshal([]byte(scriptContent), &secondTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)
				return
			}
			author = secondTypeMetaData.Author.Name
		})
		if author != "" {
			break
		}
	}
	return author, published_at
}

func (t *Template) BBCNewsPublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

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
			var firstTypeMetaData BBCNewsMetaData
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)

			} else {
				publishedAt = firstTypeMetaData.DatePublished.Unix()
			}

			if publishedAt != 0 {
				return
			}
			var secondTypeMetaData BBCNewsMetadataSecond
			unmarshalErr = json.Unmarshal([]byte(scriptContent), &secondTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)
				return
			}
			publishedAt = secondTypeMetaData.DatePublished.Unix()

		})

	}
	return publishedAt
}
