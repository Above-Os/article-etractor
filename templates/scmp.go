package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) SCMPScrapContent(document *goquery.Document) string {

	contents := ""
	/*document.Find("div.sdc-article-related-stories,div.sdc-article-factbox,div.sdc-article-strapline,div.site-share-wrapper,div[data-format=floated-mpu]").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})*/

	document.Find("div[data-qa=GenericArticle-Leading],div[data-qa=GenericArticle-Content]").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents

}
