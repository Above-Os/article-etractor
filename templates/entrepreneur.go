package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) EntrepreneurScrapContent(document *goquery.Document) string {

	contents := ""
	document.Find("article>figure,div.prose").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
