package templates

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type SlashfilmMetadata struct {
	Context          string `json:"@context"`
	Type             string `json:"@type"`
	MainEntityOfPage struct {
		Type       string `json:"@type"`
		ID         string `json:"@id"`
		Breadcrumb struct {
			Type            string `json:"@type"`
			ItemListElement []struct {
				Type     string `json:"@type"`
				Position int    `json:"position"`
				Item     struct {
					ID   string `json:"@id"`
					Name string `json:"name"`
				} `json:"item"`
			} `json:"itemListElement"`
		} `json:"breadcrumb"`
	} `json:"mainEntityOfPage"`
	Headline string `json:"headline"`
	Image    struct {
		Type   string `json:"@type"`
		URL    string `json:"url"`
		Height int    `json:"height"`
		Width  int    `json:"width"`
	} `json:"image"`
	DatePublished time.Time `json:"datePublished"`
	DateModified  time.Time `json:"dateModified"`
	Author        []struct {
		Type       string   `json:"@type"`
		Name       string   `json:"name"`
		URL        string   `json:"url"`
		KnowsAbout []string `json:"knowsAbout"`
		SameAs     []string `json:"sameAs"`
	} `json:"author"`
	Publisher struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		URL  string `json:"url"`
		Logo struct {
			Type    string `json:"@type"`
			Caption string `json:"caption"`
			URL     string `json:"url"`
			Width   string `json:"width"`
			Height  string `json:"height"`
		} `json:"logo"`
		Description   string   `json:"description"`
		SameAs        []string `json:"sameAs"`
		AlternateName string   `json:"alternateName"`
	} `json:"publisher"`
	Description string `json:"description"`
}

func (t *Template) SlashfilmScrapMetaData(document *goquery.Document) (string, string) {

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
			var firstTypeMetaData SlashfilmMetadata
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

func (t *Template) SlashfilmNewsPublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

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
			var firstTypeMetaData SlashfilmMetadata
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
