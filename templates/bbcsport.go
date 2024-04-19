package templates

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type BBCSportMetaDataFirst struct {
	Context   string `json:"@context"`
	Type      string `json:"@type"`
	URL       string `json:"url"`
	Publisher struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		Logo struct {
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

type BBCSportMetaDataSecond struct {
	Context   string `json:"@context"`
	Type      string `json:"@type"`
	URL       string `json:"url"`
	Publisher struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		Logo struct {
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
		Type string `json:"@type"`
		Name string `json:"name"`
		Logo struct {
			Type string `json:"@type"`
			URL  string `json:"url"`
		} `json:"logo"`
	} `json:"author"`
}

func (t *Template) BBCSportsScrapMetaData(document *goquery.Document) (string, string) {

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

			var firstTypeMetaData BBCSportMetaDataFirst
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert convert bbc news unmarshalError %v", unmarshalErr)
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

			var secondTypeMetaData BBCSportMetaDataSecond
			unmarshalErr = json.Unmarshal([]byte(scriptContent), &secondTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert BBCSport unmarshalErr %v \n %s", unmarshalErr, scriptContent)
				return
			}
			author = secondTypeMetaData.Author.Name
		})
		if author != "" {
			break
		}
	}

	if len(author) == 0 {
		document.Find("span.qa-contributor-name.gel-long-primer").Each(func(index int, item *goquery.Selection) {
			text := item.Text()
			trimmedText := strings.TrimPrefix(text, "By ")
			if len(trimmedText) != 0 {
				if len(author) == 0 {
					author = trimmedText
				} else {
					author = author + " & " + trimmedText
				}
			}
		})
	}
	log.Printf("author last: %s", author)
	return author, published_at
}

func (t *Template) BBCSportsPublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

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
			var firstTypeMetaData BBCSportMetaDataFirst
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)
			} else {
				publishedAt = firstTypeMetaData.DatePublished.Unix()
			}

			if publishedAt != 0 {
				return
			}

			var secondTypeMetaData BBCSportMetaDataSecond
			unmarshalErr = json.Unmarshal([]byte(scriptContent), &secondTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert BBCSport unmarshalErr %v", unmarshalErr)
				return
			}
			publishedAt = secondTypeMetaData.DatePublished.Unix()

		})

	}
	if publishedAt == 0 {
		timeTag := document.Find("time.gs-o-bullet__text.qa-status-date.gs-u-align-middle.gs-u-display-inline").First()
		if datetime, exists := timeTag.Attr("datetime"); exists {
			parsedTime, err := time.Parse(time.RFC3339, datetime)
			if err != nil {
				log.Printf("Error parsing datetime: %v", err)
			}
			timestamp := parsedTime.Unix()
			publishedAt = timestamp
		} else {
			fmt.Println("Datetime attribute not found")
		}
	}
	return publishedAt
}
