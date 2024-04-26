package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) PolygonScrapContent(document *goquery.Document) string {
	contents := ""

	document.Find("aside,div.loopnav-a").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("figure.e-image--hero,div.c-entry-content").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
