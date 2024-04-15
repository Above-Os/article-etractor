package templates

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// site url https://www.nbcnews.com/world
// rss  http://feeds.nbcnews.com/feeds/worldnews
func (t *Template) NbcNewsScrapMetaData(document *goquery.Document) (string, string) {

	author := ""
	published_at := ""
	author, published_at = t.AuthorExtractFromScriptMetadata(document)
	if author == "" {
		childernFormat := "#content > div:nth-child(7) > div > article > div.article-body.breaking > div > div.article-body__section.layout-grid-container.article-body__first-section > div.article-body.layout-grid-item.layout-grid-item--with-gutter-s-only.grid-col-10-m.grid-col-push-1-m.grid-col-6-xl.grid-col-push-2-xl.article-body--custom-column > section > div.article-inline-byline > span:nth-child(%s)"
		childernFormatLink := "#content > div:nth-child(7) > div > article > div.article-body.breaking > div > div.article-body__section.layout-grid-container.article-body__first-section > div.article-body.layout-grid-item.layout-grid-item--with-gutter-s-only.grid-col-10-m.grid-col-push-1-m.grid-col-6-xl.grid-col-push-2-xl.article-body--custom-column > section > div.article-inline-byline > span:nth-child(%s) > a"

		index := 1
		for {
			currentAuthor := ""
			currentChildern := fmt.Sprintf(childernFormat, index)
			currentChildernLink := fmt.Sprintf(childernFormatLink, index)
			document.Find(currentChildern).Each(func(i int, s *goquery.Selection) {
				currentAuthor = strings.TrimSpace(s.Text())
				// s.Find("b")
			})
			if currentAuthor != "" {
				if index != 1 {
					author = author + " & "
				}
				author = author + currentAuthor
				index = index + 1
				continue
			}
			document.Find(currentChildernLink).Each(func(i int, s *goquery.Selection) {
				currentAuthor = strings.TrimSpace(s.Text())
			})
			if currentAuthor != "" {
				if index != 1 {
					author = author + " & "
				}
				author = author + currentAuthor
				index = index + 1
				continue
			}
			if currentAuthor == "" {
				break
			}

		}
	}

	if author == "" {
		directAuthorSelector := "#content > div:nth-child(7) > div > article > div.article-body > div > div.article-body__section.layout-grid-container.article-body__last-section.article-body__first-section > div.article-body.layout-grid-item.layout-grid-item--with-gutter-s-only.grid-col-10-m.grid-col-push-1-m.grid-col-6-xl.grid-col-push-2-xl.article-body--custom-column > section > div.article-inline-byline > span > a"
		document.Find(directAuthorSelector).Each(func(i int, s *goquery.Selection) {
			author = strings.TrimSpace(s.Text())
		})
	}

	return author, published_at
}

func (t *Template) NbcNewsScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("h1.article-hero-headline__htag,aside").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("section.article-hero__container").Each(func(i int, s *goquery.Selection) {
		content, _ := goquery.OuterHtml(s)
		contents += content
	})
	articleSectionBody := document.Find("div.article-body__first-section").First()
	articleBody := articleSectionBody.Find("div").First()
	firstArticle := articleBody.Find("div.article-body__content").First()

	articleContent, _ := goquery.OuterHtml(firstArticle)
	contents += articleContent

	return contents
}
