package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) CFainstituteScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("div.date,h1.post-title,div.cfa_meta,figure.wp-block-image,div.author-details-list,div#comments,div.cfai_author_header,div.cfai_social_share,ul.share-links").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("div.post-item").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
