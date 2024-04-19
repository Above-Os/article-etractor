package templates

import (
	"github.com/PuerkitoBio/goquery"
	//"recommend.common/logger"
)

func (t *Template) ThevergeScrapContent(document *goquery.Document) string {

	contents := ""

	document.Find("button").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("div.duet--article--article-body-component-container").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
