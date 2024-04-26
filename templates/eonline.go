package templates

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type EOnlineMetadata struct {
	Context          string `json:"@context"`
	Type             string `json:"@type"`
	MainEntityofPage struct {
		Type      string `json:"@type"`
		Speakable struct {
			Type        string   `json:"@type"`
			CSSSelector []string `json:"cssSelector"`
		} `json:"speakable"`
		ID string `json:"@id"`
	} `json:"mainEntityofPage"`
	Headline string `json:"headline"`
	Image    struct {
		Type    string `json:"@type"`
		URL     string `json:"url"`
		Width   int    `json:"width"`
		Height  int    `json:"height"`
		Caption string `json:"caption"`
	} `json:"image"`
	DatePublished string `json:"datePublished"`
	DateModified  string `json:"dateModified"`
	Author        []struct {
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
		Description string `json:"description"`
	} `json:"publisher"`
}

func (t *Template) EonlineScrapMetaData(document *goquery.Document) (string, string) {

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
			var firstTypeMetaData EOnlineMetadata
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)
				return
			}
			for _, currentAuthor := range firstTypeMetaData.Author {
				if len(author) != 0 {
					author = author + " & "
				}
				author = author + currentAuthor.Name

			}
		})
		if author != "" {
			break
		}
	}
	log.Printf("author last: %s", author)
	return author, published_at
}

func (t *Template) EonlinePublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

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
			var firstTypeMetaData EOnlineMetadata
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)
				return

			}
			fmt.Println(firstTypeMetaData.DatePublished)
			currentParsedPublishedAt, convertErr := ConvertISORFC3339MiliToTimestamp(firstTypeMetaData.DatePublished)
			if convertErr != nil {
				log.Printf("convert str timestamp error %v", convertErr)
				return
			}
			publishedAt = currentParsedPublishedAt
		})

	}
	return publishedAt
}

func ConvertISORFC3339MiliToTimestamp(dateStr string) (int64, error) {
	const layout = "2006-01-02T15:04:05.000Z"

	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return 0, err // 如果解析出错，返回错误
	}

	return t.Unix(), nil
}

func (t *Template) EOnlineScrapContent(document *goquery.Document) string {

	contents := ""

	document.Find("div.article-detail__main-content").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
