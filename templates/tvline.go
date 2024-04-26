package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) TVLineScrapContent(document *goquery.Document) string {
	contents := ""

	document.Find("div[data-component=social-media],div[data-component=cards-related-content]").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)

	})

	document.Find("div[data-component=featured-media],div[data-component=gutenberg-content]").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
