package templates

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type CNCFMetaData struct {
	Context string `json:"@context"`
	Graph   []struct {
		Type            string `json:"@type"`
		ID              string `json:"@id"`
		URL             string `json:"url"`
		Name            string `json:"name"`
		InLanguage      string `json:"inLanguage"`
		PotentialAction struct {
			Type   string `json:"@type"`
			Target struct {
				Type        string `json:"@type"`
				URLTemplate string `json:"urlTemplate"`
			} `json:"target"`
			QueryInput string `json:"query-input"`
		} `json:"potentialAction,omitempty"`
		Publisher struct {
			Type   string   `json:"@type"`
			ID     string   `json:"@id"`
			Name   string   `json:"name"`
			URL    string   `json:"url"`
			SameAs []string `json:"sameAs"`
			Logo   struct {
				Type        string `json:"@type"`
				URL         string `json:"url"`
				ContentURL  string `json:"contentUrl"`
				Width       int    `json:"width"`
				Height      int    `json:"height"`
				ContentSize string `json:"contentSize"`
			} `json:"logo"`
		} `json:"publisher,omitempty"`
		Description string `json:"description,omitempty"`
		IsPartOf    struct {
			ID string `json:"@id"`
		} `json:"isPartOf,omitempty"`
		Breadcrumb struct {
			Type            string `json:"@type"`
			ID              string `json:"@id"`
			ItemListElement []struct {
				Type     string `json:"@type"`
				Position int    `json:"position"`
				Item     string `json:"item,omitempty"`
				Name     string `json:"name"`
			} `json:"itemListElement"`
		} `json:"breadcrumb,omitempty"`
		DatePublished time.Time `json:"datePublished,omitempty"`
		DateModified  time.Time `json:"dateModified,omitempty"`
		Author        struct {
			Type string `json:"@type"`
			ID   string `json:"@id"`
			Name string `json:"name"`
		} `json:"author,omitempty"`
	} `json:"@graph"`
}

func (t *Template) CNCFScrapContent(document *goquery.Document) string {

	contents := ""
	document.Find("div.social-share,figure.wp-block-embed-twitter,div.post-author").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	document.Find("article.container").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}

func (t *Template) CNCFScrapMetaData(document *goquery.Document) (string, string) {

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
			var firstTypeMetaData CNCFMetaData
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)
				return
			}
			for _, currentGraph := range firstTypeMetaData.Graph {
				if len(currentGraph.Author.Name) != 0 {
					author = author + " & " + currentGraph.Author.Name
				} else {
					author = currentGraph.Author.Name
				}
			}
		})
		if author != "" {
			break
		}
	}
	log.Printf("author last: %s", author)
	return author, published_at
}

func (t *Template) CNCFPublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

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
			var firstTypeMetaData CNCFMetaData
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)
				return

			}

			for _, currentGraph := range firstTypeMetaData.Graph {
				if currentGraph.DatePublished.Unix() > 0 {
					publishedAt = currentGraph.DatePublished.Unix()
					break
				}
			}
			//publishedAt = firstTypeMetaData[0].DatePublished.Unix()
		})

	}
	return publishedAt
}
