package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) MessengerScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("div.border-grey-3,div.custom-twitter-embed").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	articleSectionBody := document.Find("article.mx-auto").First()
	articleBody := articleSectionBody.Find("article").First()
	articleDiv := articleBody.Find("div").First()
	childrenElements := articleDiv.Children()
	headerNode := childrenElements.Eq(1)
	headerNode.Find("h2.text-navy,figure.text-grey-3").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	document.Find("article.prose").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
