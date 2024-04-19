package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) ThisisGoingtobeBIGScrapContent(document *goquery.Document) string {
	contents := ""

	document.Find("h2").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "Where to Find Me" {
			RemoveNodes(s)
		}

	})

	document.Find("div.sqs-html-content").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
