package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) DigitalTrendsScrapContent(document *goquery.Document) string {

	contents := ""
	document.Find("div.b-related-links,ul.h-editors-recs,h4.h-editors-recs-title,div#dt-toc").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	document.Find("article").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
