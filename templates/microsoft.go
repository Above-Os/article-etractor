package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) MicrosoftScrapContent(document *goquery.Document) string {
	contents := ""

	document.Find("header").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)

	})
	document.Find("article.m-blog-post").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
