package templates

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type BusinessInsiderMetaData struct {
	Context          string `json:"@context"`
	Type             string `json:"@type"`
	MainEntityOfPage struct {
		Type string `json:"@type"`
		ID   string `json:"@id"`
	} `json:"mainEntityOfPage"`
	Editor struct {
		Type string `json:"@type"`
		Name string `json:"name"`
	} `json:"editor"`
	Author struct {
		Type   string `json:"@type"`
		Name   string `json:"name"`
		SameAs string `json:"sameAs"`
	} `json:"author"`
	Publisher struct {
		Context      string   `json:"@context"`
		Type         string   `json:"@type"`
		Name         string   `json:"name"`
		LegalName    string   `json:"legalName"`
		FoundingDate string   `json:"foundingDate"`
		URL          string   `json:"url"`
		SameAs       []string `json:"sameAs"`
		Founder      struct {
			Type string `json:"@type"`
			Name string `json:"name"`
		} `json:"founder"`
		Logo struct {
			Type   string `json:"@type"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"logo"`
	} `json:"publisher"`
	Headline            string `json:"headline"`
	AlternativeHeadline string `json:"alternativeHeadline"`
	Image               struct {
		Type    string `json:"@type"`
		URL     string `json:"url"`
		Width   int    `json:"width"`
		Height  int    `json:"height"`
		Caption string `json:"caption"`
	} `json:"image"`
	Name                string    `json:"name"`
	DatePublished       time.Time `json:"datePublished"`
	DateModified        time.Time `json:"dateModified"`
	Description         string    `json:"description"`
	Keywords            string    `json:"keywords"`
	ArticleBody         string    `json:"articleBody"`
	ArticleSection      string    `json:"articleSection"`
	IsAccessibleForFree bool      `json:"isAccessibleForFree"`
	HasPart             []struct {
		Type                string `json:"@type"`
		IsAccessibleForFree bool   `json:"isAccessibleForFree"`
		CSSSelector         string `json:"cssSelector"`
	} `json:"hasPart"`
	IsPartOf struct {
		Type      []string `json:"@type"`
		Name      string   `json:"name"`
		ProductID string   `json:"productID"`
	} `json:"isPartOf"`
}

type BusinessInsiderMetaDataSecond struct {
	Context          string `json:"@context"`
	Type             string `json:"@type"`
	MainEntityOfPage struct {
		Type string `json:"@type"`
		ID   string `json:"@id"`
	} `json:"mainEntityOfPage"`
	Editor struct {
		Type string `json:"@type"`
		Name string `json:"name"`
	} `json:"editor"`
	Author []struct {
		Type   string `json:"@type"`
		Name   string `json:"name"`
		SameAs string `json:"sameAs"`
	} `json:"author"`
	Publisher struct {
		Context      string   `json:"@context"`
		Type         string   `json:"@type"`
		Name         string   `json:"name"`
		LegalName    string   `json:"legalName"`
		FoundingDate string   `json:"foundingDate"`
		URL          string   `json:"url"`
		SameAs       []string `json:"sameAs"`
		Founder      struct {
			Type string `json:"@type"`
			Name string `json:"name"`
		} `json:"founder"`
		Logo struct {
			Type   string `json:"@type"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"logo"`
	} `json:"publisher"`
	Headline            string `json:"headline"`
	AlternativeHeadline string `json:"alternativeHeadline"`
	Image               struct {
		Type    string `json:"@type"`
		URL     string `json:"url"`
		Width   int    `json:"width"`
		Height  int    `json:"height"`
		Caption string `json:"caption"`
	} `json:"image"`
	Name                string    `json:"name"`
	DatePublished       time.Time `json:"datePublished"`
	DateModified        time.Time `json:"dateModified"`
	Description         string    `json:"description"`
	Keywords            string    `json:"keywords"`
	ArticleBody         string    `json:"articleBody"`
	ArticleSection      string    `json:"articleSection"`
	IsAccessibleForFree bool      `json:"isAccessibleForFree"`
	HasPart             []struct {
		Type                string `json:"@type"`
		IsAccessibleForFree bool   `json:"isAccessibleForFree"`
		CSSSelector         string `json:"cssSelector"`
	} `json:"hasPart"`
	IsPartOf struct {
		Type      []string `json:"@type"`
		Name      string   `json:"name"`
		ProductID string   `json:"productID"`
	} `json:"isPartOf"`
}

func (t *Template) BusinessinsiderScrapMetaData(document *goquery.Document) (string, string) {

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
			var firstTypeMetaData BusinessInsiderMetaData
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)

			} else {
				author = firstTypeMetaData.Author.Name
			}
			if len(author) != 0 {
				return
			}
			var secondTypeMetaData BusinessInsiderMetaDataSecond
			unmarshalErr = json.Unmarshal([]byte(scriptContent), &secondTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)

			} else {
				for _, currentAuthor := range secondTypeMetaData.Author {
					if len(currentAuthor.Name) != 0 {
						if len(author) != 0 {
							author = author + " & " + currentAuthor.Name
						} else {
							author = currentAuthor.Name
						}
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

func (t *Template) BusinessinsiderPublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

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
			var firstTypeMetaData BusinessInsiderMetaData
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)

			} else {
				publishedAt = firstTypeMetaData.DatePublished.Unix()
			}
			if publishedAt != 0 {
				return
			}
			var secondTypeMetaData BusinessInsiderMetaDataSecond
			unmarshalErr = json.Unmarshal([]byte(scriptContent), &secondTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)

			} else {
				publishedAt = secondTypeMetaData.DatePublished.Unix()
			}
		})

	}
	return publishedAt
}

func (t *Template) BusinessInsiderScrapContent(document *goquery.Document) string {

	contents := ""

	document.Find("label.caption-drawer-label,div.piano-inline-content-wrapper,div.in-post-sticky,div.inline-newsletter-signup,div.ad-callout-wrapper,article.d-none,section.content-recommendations-component").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("figure.image-figure-image,section.post-content").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
