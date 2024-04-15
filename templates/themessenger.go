package templates

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// site url  http://www.dw.com/english/?maca=en-rss-en-world-4025-rdf
// rss  https://rss.dw.com/rdf/rss-en-world
func (t *Template) ThemessengerScrapMetaData(document *goquery.Document) (string, string) {

	author := ""
	published_at := ""

	if author == "" {
		author, published_at = t.AuthorExtractFromScriptMetadata(document)
	}

	return author, published_at
}

func (t *Template) TheMessengerGetPublishedAtTimestampSingleJson(document *goquery.Document) int64 {

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
			log.Printf("currentPublishedAtStr %s", currentPublishedAtStr+"Z")
			publishedAtTimestamp = ConvertStringTimeToTimestamp(currentPublishedAtStr + "Z")

			// var jsonMap map[string]interface{}

			// logger.Info("script content %s  author length %d",scriptContent, len(currentDWMetadata.Author))
		})
		if publishedAtTimestamp != 0 {
			break
		}
	}

	return publishedAtTimestamp

}
