package templates

import "github.com/PuerkitoBio/goquery"

func (t *Template) MediumScrapContent(document *goquery.Document) string {

	document.Find("h1.pw-post-title,h2.pw-subtitle-paragraph,div.speechify-ignore").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	contents := ""

	document.Find("section").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
