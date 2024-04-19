package templates

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/mitchellh/mapstructure"
)

type Image struct {
	ImageType string `json:"@type"`
	Url       string `json:"url"`
	Height    int    `json:"height"`
	Width     int    `json:"width"`
}

type SingleAuthor struct {
	AuthorType string `json:"@type"`
	Name       string `json:"name"`
}

type ArticleScriptMetadata struct {
	Author []SingleAuthor `json:"author"`
}

type ArticleScriptMetadataSingleAuthor struct {
	Author SingleAuthor `json:"author"`
}

type ListArticleScriptMetadata struct {
}

func ExtractAuthorFromScriptMetaData(scriptContent string) string {
	author := ""
	var currentArticleScriptMetadata ArticleScriptMetadata

	unmarshalErr := json.Unmarshal([]byte(scriptContent), &currentArticleScriptMetadata)
	if unmarshalErr != nil {
		var currentArticleScriptMetadataSingleAuthor ArticleScriptMetadataSingleAuthor
		unmarshalErr = json.Unmarshal([]byte(scriptContent), &currentArticleScriptMetadataSingleAuthor)
		if unmarshalErr == nil {
			return currentArticleScriptMetadataSingleAuthor.Author.Name
		}

	} else {
		for index, currentAuthor := range currentArticleScriptMetadata.Author {
			if index != 0 {
				author = author + " & "
			}
			author = author + currentAuthor.Name
		}
	}
	return author
}

func (t *Template) CommonGetPublishedAtTimestamp(document *goquery.Document) int64 {
	var publishedAtTimestamp int64 = 0
	publishedAtTimestamp = t.CommonGetPublishedAtTimestampSingleJson(document)

	if publishedAtTimestamp == 0 {
		publishedAtTimestamp = t.CommonGetPublishedAtTimestampMultipleJson(document)
	}
	return publishedAtTimestamp

}

func ConvertStringTimeToTimestampForEuroNews(currentTime string) int64 {

	layout := time.DateTime

	t, err := time.Parse(layout, currentTime)
	if err != nil {
		log.Printf("convert str time to golang time fail %s error %v", currentTime, err)
		return 0
	}
	return t.Unix() - 3600
}

func ConvertStringTimeToTimestamp(currentTime string) int64 {
	layout := time.RFC3339Nano

	t, err := time.Parse(layout, currentTime)
	if err != nil {
		log.Printf("convert str time to golang time fail %s error %v", currentTime, err)
		return 0
	}
	return t.Unix()
}

func ConvertStringTimeToTimestampRFC33399(currentTime string) int64 {
	layout := time.RFC3339

	t, err := time.Parse(layout, currentTime)
	if err != nil {
		log.Printf("convert str time to golang time fail %s error %v", currentTime, err)
		return 0
	}
	return t.Unix()
}

func (t *Template) CommonGetPublishedAtTimestampSingleJson(document *goquery.Document) int64 {

	var publishedAtTimestamp int64 = 0
	scriptSelectorFirst := "head > script[type=\"application/ld+json\"]"
	scriptSelectorSecond := "body > script[type=\"application/ld+json\"]"
	scriptSelectorThird := "script[type=\"application/ld+json\"]"
	scriptSelectorList := make([]string, 100)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorFirst)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorSecond)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorThird)

	for _, scriptSelector := range scriptSelectorList {

		document.Find(scriptSelector).Each(func(i int, s *goquery.Selection) {
			if publishedAtTimestamp != 0 {
				return
			}
			scriptContent := strings.TrimSpace(s.Text())

			var jsonMap map[string]interface{}
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &jsonMap)
			if unmarshalErr != nil {
				log.Printf("unmarshal error")
				return
			}
			currentPublishedAt, ok := jsonMap["datePublished"]
			if !ok {

				return
			}
			currentPublishedAtStr := currentPublishedAt.(string)
			log.Printf("currentPublishedAtStr %s", currentPublishedAtStr)
			publishedAtTimestamp = ConvertStringTimeToTimestamp(currentPublishedAtStr)
			if publishedAtTimestamp == 0 {
				publishedAtTimestamp = ConvertStringTimeToTimestampRFC33399(currentPublishedAtStr)
			}

			// var jsonMap map[string]interface{}

			// logger.Info("script content %s  author length %d",scriptContent, len(currentDWMetadata.Author))
		})
		if publishedAtTimestamp != 0 {
			break
		}
	}

	return publishedAtTimestamp

}

func (t *Template) CommonGetPublishedAtTimestampMultipleJson(document *goquery.Document) int64 {

	var publishedAtTimestamp int64 = 0

	scriptSelectorFirst := "head > script[type=\"application/ld+json\"]"
	scriptSelectorSecond := "body > script[type=\"application/ld+json\"]"
	scriptSelectorThird := "script[type=\"application/ld+json\"]"

	scriptSelectorList := make([]string, 100)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorFirst)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorSecond)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorThird)

	for _, scriptSelector := range scriptSelectorList {
		document.Find(scriptSelector).Each(func(i int, s *goquery.Selection) {
			if publishedAtTimestamp != 0 {
				return
			}
			scriptContent := strings.TrimSpace(s.Text())
			jsonMap := []map[string]interface{}{}

			unmarshalErr := json.Unmarshal([]byte(scriptContent), &jsonMap)
			if unmarshalErr != nil {
				log.Printf("unmarshal error  %v", unmarshalErr)
				return
			}

			for _, currentJsonMap := range jsonMap {
				currentPublishedAt, ok := currentJsonMap["datePublished"]
				if ok {
					currentPublishedAtStr := currentPublishedAt.(string)
					log.Printf("currentPublishedAtStr %s", currentPublishedAtStr)
					publishedAtTimestamp = ConvertStringTimeToTimestamp(currentPublishedAtStr)
					if publishedAtTimestamp == 0 {
						publishedAtTimestamp = ConvertStringTimeToTimestampRFC33399(currentPublishedAtStr)
					}
					break
				}
			}
			// var jsonMap map[string]interface{}
			// author = ExtractAuthorFromScriptMetaData(scriptContent)

			// logger.Info("script content %s  author length %d",scriptContent, len(currentDWMetadata.Author))
		})
		if publishedAtTimestamp != 0 {
			break
		}
	}

	return publishedAtTimestamp
}

func (t *Template) AuthorExtractFromListScriptMetadata(document *goquery.Document) (string, string) {

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
			jsonMap := []map[string]interface{}{}

			unmarshalErr := json.Unmarshal([]byte(scriptContent), &jsonMap)
			if unmarshalErr != nil {
				log.Printf("unmarshal error  %v", unmarshalErr)
				return
			}

			for _, currentJsonMap := range jsonMap {
				_, ok := currentJsonMap["author"]
				if ok {
					var currentMetadata ArticleScriptMetadata
					decodeErr := mapstructure.Decode(currentJsonMap, &currentMetadata)
					if decodeErr != nil {
						log.Printf("decode content to ArticleScriptMetadata fail")
						continue
					}
					for index, currentAuthor := range currentMetadata.Author {
						if index != 0 {
							author = author + " & "
						}
						author = author + currentAuthor.Name
					}
					break
				}
			}
			// var jsonMap map[string]interface{}
			// author = ExtractAuthorFromScriptMetaData(scriptContent)

			// logger.Info("script content %s  author length %d",scriptContent, len(currentDWMetadata.Author))
		})
		if author != "" {
			break
		}
	}

	return author, published_at
}

func (t *Template) AuthorExtractFromScriptMetadata(document *goquery.Document) (string, string) {

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

			var jsonMap map[string]interface{}
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &jsonMap)
			if unmarshalErr != nil {
				log.Printf("unmarshal error")
				return
			}
			_, ok := jsonMap["author"]
			if !ok {

				return
			}
			// var jsonMap map[string]interface{}
			author = ExtractAuthorFromScriptMetaData(scriptContent)

			// logger.Info("script content %s  author length %d",scriptContent, len(currentDWMetadata.Author))
		})
		if author != "" {
			break
		}
	}

	return author, published_at
}
