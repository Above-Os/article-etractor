package templates

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type ViceMetaData struct {
	Context string `json:"@context"`
	Graph   []struct {
		Type    string `json:"@type"`
		Context string `json:"@context,omitempty"`
		Author  struct {
			Type string `json:"@type"`
			Name string `json:"name"`
		} `json:"author,omitempty"`
		DateModified     string   `json:"dateModified,omitempty"`
		DatePublished    string   `json:"datePublished,omitempty"`
		Headline         string   `json:"headline,omitempty"`
		Image            []string `json:"image,omitempty"`
		MainEntityOfPage struct {
			ID   string `json:"@id"`
			Type string `json:"@type"`
		} `json:"mainEntityOfPage,omitempty"`
		Publisher struct {
			Type string `json:"@type"`
			Logo struct {
				Type string `json:"@type"`
				URL  string `json:"url"`
			} `json:"logo"`
			Name string `json:"name"`
		} `json:"publisher,omitempty"`
		Name string `json:"name,omitempty"`
	} `json:"@graph"`
}

func (t *Template) ViceScrapMetaData(document *goquery.Document) (string, string) {

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
			validScriptContent, modifyViceInvalidJsonErr := modifyViceInvalidJSON(scriptContent)
			if modifyViceInvalidJsonErr != nil {
				log.Printf("convert Vic invalid json to valid json error %v", modifyViceInvalidJsonErr)
				return
			}

			var firstTypeMetaData ViceMetaData
			unmarshalErr := json.Unmarshal([]byte(validScriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert viceMetadata unmarshalError %v", unmarshalErr)
				return
			}
			for _, currentGraph := range firstTypeMetaData.Graph {
				if len(currentGraph.Author.Name) != 0 {
					if len(author) != 0 {
						author = author + " & " + currentGraph.Author.Name
					} else {
						author = currentGraph.Author.Name
					}

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

func (t *Template) VicePublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

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
			validScriptContent, modifyViceInvalidJsonErr := modifyViceInvalidJSON(scriptContent)
			if modifyViceInvalidJsonErr != nil {
				log.Printf("convert Vic invalid json to valid json error %v", modifyViceInvalidJsonErr)
				return
			}
			var firstTypeMetaData ViceMetaData
			unmarshalErr := json.Unmarshal([]byte(validScriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)
				return

			}
			for _, currentGraph := range firstTypeMetaData.Graph {
				if publishedAt != 0 {
					break
				}
				convertPublieshedAt, parseViceTimeErr := parseViceTimestamp(currentGraph.DatePublished)
				if parseViceTimeErr != nil {
					log.Printf("error convert Vice time to timestamp")
					continue
				}
				publishedAt = convertPublieshedAt

				// publishedAt = currentGraph.DatePublished.Unix()

			}
			//publishedAt = firstTypeMetaData[0].DatePublished.Unix()
		})

	}
	return publishedAt
}

func parseViceTimestamp(timeStr string) (int64, error) {
	// Parse the time string with the specific layout
	t, err := time.Parse(time.RFC3339Nano, timeStr)
	if err != nil {
		return 0, err
	}

	// Return the Unix timestamp (seconds since January 1, 1970 UTC)
	return t.Unix(), nil
}

func modifyViceInvalidJSON(input string) (string, error) {
	var data map[string]interface{}

	if err := json.Unmarshal([]byte(input), &data); err != nil {
		return "", err
	}

	if graph, found := data["@graph"].([]interface{}); found {
		var validGraph []interface{}
		for _, item := range graph {
			if _, ok := item.([]interface{}); !ok {
				validGraph = append(validGraph, item)
			}
		}
		data["@graph"] = validGraph
	}

	modifiedJSON, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}

	return string(modifiedJSON), nil
}

func (t *Template) ViceScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("div.adph,div.abc__article_embed").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	//https://www.vice.com/en/article/wxjymy/i-rode-melbournes-free-taylor-swift-trams-for-12-hours-because-i-love-free-things-and-hate-myself
	/*
			 <picture class="responsive-image lazyloader--lazy lazyloader--lowres"><source media="(min-width: 1000px)" srcSet="https://video-images.vice.com/_uncategorized/1708388089098-img0972.jpeg?resize=20:*"/>
		          <source media="(min-width: 700px)" srcSet="https://video-images.vice.com/_uncategorized/1708388089098-img0972.jpeg?resize=20:*"/>
		          <source media="(min-width: 0px)" srcSet="https://video-images.vice.com/_uncategorized/1708388089098-img0972.jpeg?resize=20:*"/>
		          <img class="responsive-image__img" alt="man on tram" decoding="async" loading="eager" width="2316" height="1745"/></picture>
	*/
	var rxPicture = regexp.MustCompile(`resize=(\d+):`)
	document.Find("picture").Each(func(i int, pic *goquery.Selection) {
		firstChild := pic.Children().First()
		childNode := firstChild.Get(0)
		if childNode.Data == "source" {
			src, exists := firstChild.Attr("srcset")
			if exists {
				replaceVal := rxPicture.ReplaceAllString(src, "resize=650:")
				firstChild.SetAttr("srcset", replaceVal)
			}

		}
	})

	document.Find("div.short-form__body__article-lede-image,div.article__body-components").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
