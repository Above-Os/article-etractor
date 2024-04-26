package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) ThevergeScrapContent(document *goquery.Document) string {

	contents := ""

	document.Find("button,div.duet--recirculation--related-list").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("figure.duet--article--lede-image,div.duet--article--article-body-component-container").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
