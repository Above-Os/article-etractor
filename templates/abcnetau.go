package templates

import "github.com/PuerkitoBio/goquery"

func (t *Template) AbcNetAUScrapMetaData(document *goquery.Document) (string, string) {
	author := "abc.net.au"
	published_at := ""

	return author, published_at
}

func (t *Template) AbcNetAUScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("aside,div.Headline_meta__ZgyGe,div[data-component=RelatedTopics],div[data-component=ShareUtility],div[data-component=Dateline]").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	document.Find("div.Article_layoutMain__eBEMA").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
