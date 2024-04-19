package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) DeadspinScrapContent(document *goquery.Document) string {

	contents := ""

	document.Find("div.js_post-content").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
