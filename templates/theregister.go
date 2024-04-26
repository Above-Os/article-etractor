package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) TheRegisterScrapContent(document *goquery.Document) string {

	contents := ""
	document.Find("ul.listinks,div[aria-hidden=true]").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("div#body").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
