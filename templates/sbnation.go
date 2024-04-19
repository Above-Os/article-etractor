package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) SbnationScrapContent(document *goquery.Document) string {

	contents := ""

	document.Find("div#comments,section.c-nextclick").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("div.l-col__main").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
