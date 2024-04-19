package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) VarietyScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("aside,div.article-tags,div.o-comments-link").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("img").Each(func(i int, img *goquery.Selection) {
		src, exists := img.Attr("data-lazy-src")
		if exists {
			img.SetAttr("src", src)
		}
	})

	document.Find("div.article-header__feature,div.vy-cx-page-content").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
