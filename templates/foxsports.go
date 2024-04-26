package templates

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type FoxsportsMetaData struct {
	Context   string `json:"@context"`
	Type      string `json:"@type"`
	Headline  string `json:"headline"`
	Speakable struct {
		Type  string   `json:"@type"`
		Xpath []string `json:"xpath"`
		URL   string   `json:"url"`
	} `json:"speakable"`
	ArticleBody string `json:"articleBody"`
	Image       []struct {
		Type   string `json:"@type"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"image"`
	DatePublished    string `json:"datePublished"`
	DateModified     string `json:"dateModified"`
	MainEntityOfPage struct {
		Type string `json:"@type"`
		ID   string `json:"@id"`
	} `json:"mainEntityOfPage"`
	Author    []any `json:"author"`
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

func (t *Template) FoxSportsScrapContent(document *goquery.Document) string {

	contents := ""
	document.Find("h1.story-title,div.story-header-container,div.fwAdContainer,div.storyFavoriteContainer,div.story-social-group,div.story-topic-group,div.story-favorites-section-add").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	document.Find("div.story-content-main").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}

func (t *Template) FoxsportsScrapMetaData(document *goquery.Document) (string, string) {

	author := ""
	published_at := ""
	scriptSelectorFirst := "head > script[type=\"application/ld+json\"]"
	scriptSelectorSecond := "body > script[type=\"application/ld+json\"]"
	scriptSelectorThird := "script[type=\"application/ld+json\"]"

	scriptSelectorList := make([]string, 100)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorFirst)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorSecond)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorThird)
	fmt.Println("11111111111111111111111")
	for _, scriptSelector := range scriptSelectorList {

		document.Find(scriptSelector).Each(func(i int, s *goquery.Selection) {
			if author != "" {
				return
			}
			scriptContent := strings.TrimSpace(s.Text())
			var firstTypeMetaData FoxsportsMetaData
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)
				return
			}
			parseAuthor, parseAuthorErr := extractAuthorFoxSports(firstTypeMetaData.ArticleBody)
			if parseAuthorErr != nil {
				author = parseAuthor
			}

		})
		if author != "" {
			break
		}
	}
	log.Printf("author last: %s", author)
	return author, published_at
}

func (t *Template) FoxsportsPublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

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
			var firstTypeMetaData FoxsportsMetaData
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)
				return
			}
			fmt.Printf("------------------- %s", firstTypeMetaData.DatePublished)

			//publishedAt = firstTypeMetaData[0].DatePublished.Unix()
		})

	}
	return publishedAt
}

func extractAuthorFoxSports(text string) (string, error) {
	re, err := regexp.Compile(`^By (.*?) FOX Sports`)
	if err != nil {
		return "", fmt.Errorf("error compiling regex: %v", err)
	}

	match := re.FindStringSubmatch(text)
	if match != nil && len(match) > 1 {
		return match[1], nil // 返回捕获的作者名
	}
	return "", fmt.Errorf("no match found")
}
