package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) IndependentUKScrapContent(document *goquery.Document) string {

	contents := ""
	document.Find("div.newsletter-component,nav.topics,div.show-comments,div.lightbox").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	document.Find("header#articleHeader>div:nth-child(1)>div:nth-child(3),div#main").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
