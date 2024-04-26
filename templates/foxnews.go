package templates

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// site url  https://www.foxnews.com/world
// rss  https://moxie.foxnews.com/google-publisher/world.xml
func (t *Template) FoxNewsScrapMetaData(document *goquery.Document) (string, string) {
	author := ""
	published_at := ""

	authorSelectorFirst := "#wrapper > div.page-content > div.row.full > main > article > header > div.author-byline > span:nth-child(1) > span > a"
	authorSelectorSecond := "#wrapper > div.page-content > div.row.full > main > article > header > div.author-byline > span > a"

	cssSelectorList := make([]string, 100)
	cssSelectorList = append(cssSelectorList, authorSelectorFirst)
	cssSelectorList = append(cssSelectorList, authorSelectorSecond)
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

func (t *Template) FoxNewsScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("div.article-gating-wrapper,div.ad-container").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("a").Each(func(i int, s *goquery.Selection) {
		text := s.Text()

		if strings.HasPrefix(text, "MORE:") || strings.HasPrefix(text, "close") {
			RemoveNodes(s)
		}

	})

	document.Find("div.paywall,div.article-body").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
