package templates

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ESPNMetaDataWithAuthor struct {
	Context          string `json:"@context"`
	Type             string `json:"@type"`
	MainEntityOfPage struct {
		Type string `json:"@type"`
		ID   string `json:"@id"`
	} `json:"mainEntityOfPage"`
	Headline      string `json:"headline"`
	Description   string `json:"description"`
	DatePublished string `json:"datePublished"`
	DateModified  string `json:"dateModified"`
	Image         struct {
		Type   string `json:"@type"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"image"`
	Author struct {
		Type string `json:"@type"`
		Name string `json:"name"`
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
	} `json:"publisher"`
}

type ESPNMetaDataWithoutAuthor struct {
	Context          string `json:"@context"`
	Type             string `json:"@type"`
	MainEntityOfPage struct {
		Type string `json:"@type"`
		ID   string `json:"@id"`
	} `json:"mainEntityOfPage"`
	Headline      string `json:"headline"`
	Description   string `json:"description"`
	DatePublished string `json:"datePublished"`
	DateModified  string `json:"dateModified"`
	Image         struct {
		Type   string `json:"@type"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"image"`
	Publisher struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		Logo struct {
			Type   string `json:"@type"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"logo"`
	} `json:"publisher"`
}

func (t *Template) EspnScrapContent(document *goquery.Document) string {

	contents := ""
	/*document.Find("header.article-header,aside.float-r,div.article-meta,div.content-reactions_reactions-wrapper").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})*/
	document.Find("h2").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("article#article").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}

func (t *Template) ESPNScrapMetaData(document *goquery.Document) (string, string) {

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
			var firstTypeMetaData ESPNMetaDataWithAuthor
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert espn unmarshalError %v", unmarshalErr)

			} else {
				author = firstTypeMetaData.Author.Name
			}

			if len(author) != 0 {
				return
			}

			var secondTypeMetaData ESPNMetaDataWithoutAuthor
			unmarshalErr = json.Unmarshal([]byte(scriptContent), &secondTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert espn unmarshalError %v", unmarshalErr)

			} else {
				author = secondTypeMetaData.Publisher.Name
			}

		})
		if author != "" {
			break
		}
	}
	return author, published_at
}

func (t *Template) ESPNPublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

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
			var firstTypeMetaData ESPNMetaDataWithAuthor
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert spn unmarshalError %v", unmarshalErr)
			} else {
				parsedPublishedAt, parseErr := StringToTimestampEspn(firstTypeMetaData.DatePublished)
				if parseErr != nil {
					log.Printf("espn convert timestamp fail %v", parsedPublishedAt)

				} else {
					publishedAt = parsedPublishedAt
				}
				// publishedAt = firstTypeMetaData.DatePublished.Unix()
			}

			if publishedAt != 0 {
				return
			}
			var secondTypeMetaData ESPNMetaDataWithoutAuthor
			unmarshalErr = json.Unmarshal([]byte(scriptContent), &secondTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert espn unmarshalError %v", unmarshalErr)

			} else {
				// fmt.Printf("------------ %s",secondTypeMetaData.DatePublished)

				parsedPublishedAt, parseErr := StringToTimestampEspn(secondTypeMetaData.DatePublished)
				if parseErr != nil {
					log.Printf("espn convert timestamp fail %v", parsedPublishedAt)

				} else {
					publishedAt = parsedPublishedAt
				}
			}
			//publishedAt = firstTypeMetaData[0].DatePublished.Unix()
		})

	}
	return publishedAt
}

func StringToTimestampEspn(datetimeStr string) (int64, error) {
	layout := time.RFC3339
	t, err := time.Parse(layout, datetimeStr)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}
