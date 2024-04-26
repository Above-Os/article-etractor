package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) NYpostScrapContent(document *goquery.Document) string {

	contents := ""
	document.Find("div.single__header,aside,div[data-component=floatingShare],div.single__footer").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	document.Find("article").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
