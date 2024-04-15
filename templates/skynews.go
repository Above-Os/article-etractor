package templates

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (t *Template) SkyNewsScrapMetaData(document *goquery.Document) (string, string) {
	author := ""
	published_at := ""
	document.Find("span.sdc-article-author__name>a,div.sdc-article-author>p").Each(func(i int, s *goquery.Selection) {
		author = strings.TrimSpace(s.Text())

	})

	return author, published_at
}

func (t *Template) SkyNewsScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("div.sdc-article-related-stories,div.sdc-site-video,a,span[data-label-text=Advertisement]").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("p").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "Read more:") {
			RemoveNodes(s)
		}
	})
	document.Find("figure.sdc-article-image__figure,div.sdc-article-body").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
