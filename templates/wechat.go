package templates

import (
	"github.com/PuerkitoBio/goquery"
)

func (t *Template) WechatScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("div#meta_content,div#js_pc_qr_code,div#wx_stream_article_slide_tip,div#wx_expand_article").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("div.rich_media_content").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
