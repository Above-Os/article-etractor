package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) ScreenrantScrapContent(document *goquery.Document) string {
	contents := ""

	document.Find("div.adsninja-ad-zone,div.active-content").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("div.heading_image,section.article-body").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
