package templates

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (t *Template) BBCScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("header,div[data-component=topic-list],div[data-component=links-block],span.visually-hidden,h1#main-heading").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("div.description").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents = content
	})
	if contents == "" {
		document.Find("article,#main-content").Each(func(i int, s *goquery.Selection) {
			var content string
			content, _ = goquery.OuterHtml(s)
			contents = content
		})
	}
	return contents
}

func (t *Template) BBCScrapMetaData(document *goquery.Document) (string, string) {

	author := ""
	published_at := ""
	author, published_at = t.AuthorExtractFromScriptMetadata(document)
	if author != "" {
		byPrefix := "By "
		exist := strings.HasPrefix(author, byPrefix)
		if exist {
			author = author[len(byPrefix)-1:]
		}
	}

	return author, published_at
}
