package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) DazeddigitalScrapContent(document *goquery.Document) string {

	contents := ""
	document.Find("div.show-when-aside-content-hides,div.gallery-label,div.read-more").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("figure.main-img,div.article-body-container").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
