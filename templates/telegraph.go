package templates

import "github.com/PuerkitoBio/goquery"

func (t *Template) TelegraphScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("header,aside").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("article.grid").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
