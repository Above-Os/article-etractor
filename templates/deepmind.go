package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) DeepMindScrapContent(document *goquery.Document) string {

	contents := ""
	document.Find("div.article-cover__header,aside.related-posts").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	document.Find("article").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
