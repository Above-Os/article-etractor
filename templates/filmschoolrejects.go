package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) FilmSchoolRejectsScrapContent(document *goquery.Document) string {

	contents := ""
	document.Find("section.recommended").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("article.article").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
