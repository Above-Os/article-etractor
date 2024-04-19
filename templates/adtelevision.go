package templates

import "github.com/PuerkitoBio/goquery"

func (t *Template) AdTelevisionScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("h3.h1,div#post-meta-tags,div.col-md,div.addtoany_share_save_container,div.yarpp-related").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	document.Find("div.py-2").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
