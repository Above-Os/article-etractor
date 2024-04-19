package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) CNNScrapContent(document *goquery.Document) string {

	contents := ""

	document.Find("main.article__main").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
