package templates

import (
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

func (t *Template) BleadherReportScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("div.twitterEmbed").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	var rxPicture = regexp.MustCompile(`resize=(\d+):`)
	document.Find("picture").Each(func(i int, pic *goquery.Selection) {
		firstChild := pic.Children().First()
		childNode := firstChild.Get(0)
		if childNode.Data == "source" {
			src, exists := firstChild.Attr("srcset")
			if exists {
				replaceVal := rxPicture.ReplaceAllString(src, "resize=650:")
				firstChild.SetAttr("srcset", replaceVal)
			}

		}
	})

	document.Find("div.contentStream").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
