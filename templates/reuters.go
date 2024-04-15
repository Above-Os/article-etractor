package templates

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (t *Template) ReutersScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("p").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.HasPrefix(text, "Read More:") {
			RemoveNodes(s)
		}

	})

	document.Find("article").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
