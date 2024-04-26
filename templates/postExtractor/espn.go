package postExtractor

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (t PostExtractorTemplate) EspnPostExtractor(content, feedUrl string) string {
	templateData := strings.NewReader(content)

	doc, _ := goquery.NewDocumentFromReader(templateData)
	doc.Find("header,figcaption").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	newContent, _ := doc.Html()
	return newContent
}
