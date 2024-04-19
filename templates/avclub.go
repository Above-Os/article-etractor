package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) AVClubScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("div.instream-native-video,div.advertisement,div.js_ad-mobile-dynamic,div.js_related-stories-inset,div.js_related-stories-inset-mobile,div.js_postbottom-waypoint-hook,div.js_comments-iframe").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	/*document.Find("img").Each(func(i int, img *goquery.Selection) {
		src, exists := img.Attr("data-lazy-src")
		if exists {
			img.SetAttr("src", src)
		}
	})*/

	document.Find("div.article-header__feature,div.js_post-content").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
