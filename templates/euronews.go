package templates

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// site url  https://www.euronews.com/news/international
// rss  https://www.euronews.com/rss?format=mrss&level=theme&name=news
type EuronewsMetadata struct {
	Context string `json:"@context"`
	Graph   []struct {
		Type             string `json:"@type"`
		MainEntityOfPage struct {
			Type string `json:"@type"`
			URL  string `json:"url"`
		} `json:"mainEntityOfPage,omitempty"`
		Headline      string `json:"headline,omitempty"`
		Description   string `json:"description,omitempty"`
		ArticleBody   string `json:"articleBody,omitempty"`
		DateCreated   string `json:"dateCreated,omitempty"`
		DateModified  string `json:"dateModified,omitempty"`
		DatePublished string `json:"datePublished,omitempty"`
		Image         struct {
			Type      string `json:"@type"`
			URL       string `json:"url"`
			Width     string `json:"width"`
			Height    string `json:"height"`
			Caption   string `json:"caption"`
			Thumbnail string `json:"thumbnail"`
			Publisher struct {
				Type string `json:"@type"`
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"publisher"`
		} `json:"image,omitempty"`
		Author struct {
			Type   string `json:"@type"`
			Name   string `json:"name"`
			URL    string `json:"url"`
			SameAs string `json:"sameAs"`
		} `json:"author,omitempty"`
		Publisher struct {
			Type      string `json:"@type"`
			Name      string `json:"name"`
			LegalName string `json:"legalName"`
			URL       string `json:"url"`
			Logo      struct {
				Type   string `json:"@type"`
				URL    string `json:"url"`
				Width  string `json:"width"`
				Height string `json:"height"`
			} `json:"logo"`
			SameAs []string `json:"sameAs"`
		} `json:"publisher,omitempty"`
		Video struct {
			Type         string `json:"@type"`
			ContentURL   string `json:"contentUrl"`
			Description  string `json:"description"`
			Duration     string `json:"duration"`
			EmbedURL     string `json:"embedUrl"`
			Height       string `json:"height"`
			Name         string `json:"name"`
			ThumbnailURL string `json:"thumbnailUrl"`
			UploadDate   string `json:"uploadDate"`
			VideoQuality string `json:"videoQuality"`
			Width        string `json:"width"`
			InLanguage   struct {
				Name          string `json:"name"`
				AlternateName string `json:"alternateName"`
				Description   string `json:"description"`
				Identifier    string `json:"identifier"`
				URL           string `json:"url"`
				InLanguage    string `json:"inLanguage"`
			} `json:"inLanguage"`
			Publisher struct {
				Type      string `json:"@type"`
				Name      string `json:"name"`
				LegalName string `json:"legalName"`
				URL       string `json:"url"`
				Logo      struct {
					Type   string `json:"@type"`
					URL    string `json:"url"`
					Width  string `json:"width"`
					Height string `json:"height"`
				} `json:"logo"`
				SameAs []string `json:"sameAs"`
			} `json:"publisher"`
		} `json:"video,omitempty"`
		Speakable struct {
			Type  string   `json:"@type"`
			XPath []string `json:"xPath"`
			URL   string   `json:"url"`
		} `json:"speakable,omitempty"`
		Name            string `json:"name,omitempty"`
		URL             string `json:"url,omitempty"`
		PotentialAction struct {
			Type       string `json:"@type"`
			Target     string `json:"target"`
			QueryInput string `json:"query-input"`
		} `json:"potentialAction,omitempty"`
		SameAs []string `json:"sameAs,omitempty"`
	} `json:"@graph"`
}

func ConvertInterfaceToArrayOfMaps(value interface{}) ([]map[string]interface{}, error) {
	if valueSlice, ok := value.([]interface{}); ok {
		currentMapSlice := []map[string]interface{}{}
		for _, currentMapInterface := range valueSlice {
			currentMap, currentMapErr := ConvertInterfaceToMap(currentMapInterface)
			if currentMapErr == nil {
				currentMapSlice = append(currentMapSlice, currentMap)
			} else {

			}
		}
		return currentMapSlice, nil
	} else if valueMap, ok := value.(map[string]interface{}); ok {
		valueSlice := make([]map[string]interface{}, 1)
		valueSlice[0] = valueMap
		return valueSlice, nil
	} else {
		return nil, fmt.Errorf("value is not a []map[string]interface{} or map[string]interface{}")
	}
}

func EuronewsConvertSecond(timeStr string) int64 {
	pattern := `(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}) ([+-])(\d{2}):(\d{2})`

	re, err := regexp.Compile(pattern)
	if err != nil {
		fmt.Println("Error compiling regex:", err)
		return 0
	}
	matches := re.FindStringSubmatch(timeStr)
	if matches != nil {
		pureTimeStr := matches[1]
		positiveNegative := matches[2]
		var prefixNumber int64 = 1
		if positiveNegative == "+" {
			prefixNumber = -1
		}

		timezone := matches[3]
		strconv.ParseInt(timezone, 10, 64)
		timeZoneValue, parseIntErr := strconv.ParseInt(timezone, 10, 64)
		if parseIntErr != nil {
			log.Printf("parseIntErr %v", parseIntErr)
			return 0
		}

		log.Printf("pureTimeStr %s positiveNegative %s timezone %s", pureTimeStr, positiveNegative, timezone)
		layout := time.DateTime

		t, err := time.Parse(layout, pureTimeStr)
		if err != nil {
			log.Printf("convert str time to golang time fail %s error %v", pureTimeStr, err)
			return 0
		}
		return t.Unix() + prefixNumber*timeZoneValue*3600

	} else {
		return 0
	}

}

func (t *Template) EuroNewsGetPublishedAtTimeStampStruct(document *goquery.Document) int64 {
	// 2024-02-18 11:33:12
	// 2024-02-19 13:26:27 +01:00
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
			var firstTypeMetaData EuronewsMetadata
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert ApNewsMetaData unmarshalError %v", unmarshalErr)
				return

			}
			for _, currentGraph := range firstTypeMetaData.Graph {

				publishedAt = ConvertStringTimeToTimestampForEuroNews(currentGraph.DatePublished)
				if publishedAt == 0 {
					publishedAt = EuronewsConvertSecond(currentGraph.DatePublished)

				}
				if publishedAt != 0 {
					break
				}
			}

		})
		if publishedAt != 0 {
			break
		}
	}
	return publishedAt
}

func (t *Template) EuroNewsGetAuthorStruct(document *goquery.Document) string {
	author := ""
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
			var firstTypeMetaData EuronewsMetadata
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert ApNewsMetaData unmarshalError %v", unmarshalErr)
				return

			}
			for _, currentGraph := range firstTypeMetaData.Graph {
				author = currentGraph.Author.Name
				if author != "" {
					break
				}

			}

			// var jsonMap map[string]interface{}

			// logger.Info("script content %s  author length %d",scriptContent, len(currentDWMetadata.Author))
		})
		if author != "" {
			break
		}
	}
	log.Printf("author last: %s", author)
	return author
}

func ConvertInterfaceToMap(value interface{}) (map[string]interface{}, error) {
	if valueMap, ok := value.(map[string]interface{}); ok {
		return valueMap, nil
	} else {
		return nil, fmt.Errorf("value is not a map[string]interface{}")
	}
}
func (t *Template) EuroNewsScrapMetaData(document *goquery.Document) (string, string) {

	author := ""
	published_at := ""
	author, published_at = t.AuthorExtractFromScriptMetadata(document)
	if author == "" {
		author = t.EuroNewsGetAuthorStruct(document)
		log.Printf("author********************** [%s]", author)
	}

	return author, published_at
}

func (t *Template) EuroNewsScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("nav,h1.c-article-redesign-title,div.c-article-you-might-also-like,div.c-article-contributors,time.c-article-publication-date,div.c-ad__placeholder,a.c-article-partage-commentaire__links,div.c-article-caption,div.c-article-partage-commentaire-popup-overlay").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("div.o-article-newsy__main").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
