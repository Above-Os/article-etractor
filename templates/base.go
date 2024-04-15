package templates

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Template struct{}

func ScrapContentUseRules(document *goquery.Document, rules string) (string, error) {
	contents := ""
	document.Find(rules).Each(func(i int, s *goquery.Selection) {
		var content string

		content, _ = goquery.OuterHtml(s)
		contents += content
	})

	return contents, nil
}

func GetArticleByDivClass(document *goquery.Document) string {
	content := ""
	/*cnn  article__content
	  fortune.com articleContent
	*/
	document.Find("div.entry-content,div.content-entry,div.entry-body,div.article-detail,div.entry,div.entry__content,div.article__content,div.articleContent").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		len := usefulContentLen(text)
		if len > 300 {
			s.Children().Each(func(i int, childsection *goquery.Selection) {
				checkUnusedlDiv(childsection)
			})
			content, _ = goquery.OuterHtml(s)
		}
	})

	return content
}

func checkUnusedlDiv(s *goquery.Selection) bool {
	is_remove := false
	node := s.Get(0)
	d1 := node.Data
	content := s.Text()
	if strings.ToLower(strings.TrimSpace(content)) == "advertisement" {
		RemoveNodes(s)
		is_remove = true
	}
	if d1 == "ul" && len(content) == 0 {
		RemoveNodes(s)
		is_remove = true
	}
	return is_remove
}

func usefulContentLen(text string) int {
	content := strings.Replace(text, " ", "", -1)
	content = strings.Replace(content, "\n", "", -1)
	content = strings.Replace(content, "\t", "", -1)
	return len(content)
}

func RemoveNodes(s *goquery.Selection) {
	s.Each(func(i int, s *goquery.Selection) {
		parent := s.Parent()
		if parent.Length() > 0 {
			parent.Get(0).RemoveChild(s.Get(0))
		}
	})
}
