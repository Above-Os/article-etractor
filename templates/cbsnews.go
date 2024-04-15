package templates

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (t *Template) CbsNewsScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("aside,div.content-author,ul.content__tags").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("section.content__body").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}

func (t *Template) CbsNewsScrapMetaData(document *goquery.Document) (string, string) {
	author := ""
	published_at := ""
	author, published_at = t.AuthorExtractFromScriptMetadata(document)
	/**
		if author == ""{
			document.Find("span.byline__authors").Each(func(i int, s *goquery.Selection) {
				author = strings.TrimSpace(s.Text())

			})
	}
	*/

	return author, published_at
}

func (t *Template) CbsnewsWorldGetPublishedAtTimestampSingleJson(document *goquery.Document) int64 {

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
			log.Printf("currentPublishedAtStr %s", currentPublishedAtStr[0:len(currentPublishedAtStr)-5]+"Z")
			publishedAtTimestamp = ConvertStringTimeToTimestamp(currentPublishedAtStr[0:len(currentPublishedAtStr)-5] + "Z")
			if publishedAtTimestamp != 0 {
				publishedAtTimestamp = publishedAtTimestamp + 5*3600
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
