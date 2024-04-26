package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) TechCrunchScrapContent(document *goquery.Document) string {

	contents := ""

	document.Find("h1.article__title,div.article__byline").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("article.article-container").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
