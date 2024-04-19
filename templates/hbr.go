package templates

import (
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func (t *Template) HBRScrapContent(document *goquery.Document) string {

	contents := ""

	document.Find("div.article-body").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}

func (t *Template) HBRScrapMetaData(document *goquery.Document) (string, string) {

	author := ""
	published_at := ""
	document.Find("content[js-target='article-content']").Each(func(i int, s *goquery.Selection) {
		pText := s.Find("p").First().Text()

		if strings.HasPrefix(pText, "By ") {
			author = strings.TrimPrefix(pText, "By ")
		}
	})
	return author, published_at
}

func (t *Template) HBRPublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

	var publishedAt int64 = 0

	content, exists := document.Find("meta[property='article:published_time']").Attr("content")
	if !exists {
		log.Printf("Specified tag not found")
	} else {

		publishedTime, err := time.Parse(time.RFC3339, content)
		if err != nil {
			log.Printf("parse hbr time err")
		} else {

			publishedAt = publishedTime.Unix()
		}
	}
	return publishedAt
}
