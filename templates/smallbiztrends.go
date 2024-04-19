package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) SmallBizTrendsScrapContent(document *goquery.Document) string {

	contents := ""
	document.Find("h1.entry-title,span.byline").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	document.Find("div.post-inner").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents

}
