package templates

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (t *Template) A16ZCrptoScrapMetaData(document *goquery.Document) (string, string) {

	author := ""
	published_at := ""
	var authors []string

	// 使用 goquery 选择器找到所有具有 name="parsely-author" 的 meta 标签
	document.Find("meta[name='parsely-author']").Each(func(i int, s *goquery.Selection) {
		// 对于每个找到的标签，获取其 content 属性
		if content, exists := s.Attr("content"); exists {
			authors = append(authors, content)
		}
	})
	if len(authors) > 0 {
		author = strings.Join(authors," & ")
	}

	return author, published_at
}