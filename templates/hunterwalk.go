package templates

import "github.com/PuerkitoBio/goquery"

func (t *Template) HunterWalkScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("div.sharedaddy").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	document.Find("div.entry-content").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
