package templates

import "github.com/PuerkitoBio/goquery"

func (t *Template) TheStarScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("div.hidden-print").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("div.asset-photo,div#article-body").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
