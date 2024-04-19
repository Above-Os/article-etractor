package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) MattturckScrapContent(document *goquery.Document) string {

	contents := ""

	document.Find("aside").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("div.entry-content").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
