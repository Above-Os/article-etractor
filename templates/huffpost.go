package templates

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// site url  https://www.huffpost.com/news/world-news
// rss  https://chaski.huffpost.com/us/auto/vertical/world-news
func (t *Template) HuffPostScrapMetaData(document *goquery.Document) (string, string) {

	author := ""
	published_at := ""
	cssSelectorFirst := "#entry-header > header > div.bottom-header.js-cet-subunit > div.bottom-header__left > div.entry__wirepartner.entry-wirepartner > span"
	cssSelectorSecond := "#entry-footer > div.entry__author-cards > div > div > div > h2 > a > span"

	cssSelectorList := make([]string, 100)
	cssSelectorList = append(cssSelectorList, cssSelectorFirst)
	cssSelectorList = append(cssSelectorList, cssSelectorSecond)

	for _, cssSelector := range cssSelectorList {
		document.Find(cssSelector).Each(func(i int, s *goquery.Selection) {
			author = strings.TrimSpace(s.Text())
		})
		if author != "" {
			break
		}
	}

	return author, published_at
}

func (t *Template) HuffPostScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("div.cli-advertisement,aside,div.loading-message,div#support-huffpost-entry").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	document.Find("section#entry-body").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
