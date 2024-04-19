package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) DeadlineScrapContent(document *goquery.Document) string {

	contents := ""
	document.Find("div.injected-related-story,section.toaster,div.article-tags,div#comments-loading,h2#comments-title,p.subscribe-to").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("article").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
