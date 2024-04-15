package templates

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)
func (t *Template) ThemoscowtimesScrapMetaDataOpinionPart(document *goquery.Document) (string, string) {
	author := ""
	published_at := ""
    // optionion 部分
	     childernFormat := "div.container.article-container > div > div.col > article > div.article__byline.byline.byline--opinion > div > div.col > div > div > div > a:nth-child(%d)"
	singleChildSelector := "div.container.article-container > div > div.col > article > div.article__byline.byline.byline--opinion > div > div.col > div > div > div > a"

    index := 1
	for {
		currentAuthor := ""
		currentChildern := fmt.Sprintf(childernFormat,index)
		document.Find(currentChildern).Each(func(i int, s *goquery.Selection) {
			currentAuthor = strings.TrimSpace(s.Text())
			// s.Find("b")
		})
		if  currentAuthor != "" {
			if index != 1 {
				author = author + " & "
			}
			author = author + currentAuthor
			index = index + 1
			continue;
		}
		if currentAuthor == "" {
			break;
		}
		
	}

	if author == "" {
		document.Find(singleChildSelector).Each(func(i int, s *goquery.Selection) {
			author = strings.TrimSpace(s.Text())
		})
	}

    // news part

	return author, published_at
}


func (t *Template) ThemoscowtimesScrapMetaDataNewsPart(document *goquery.Document) (string, string) {
	author := ""
	published_at := ""
    // optionion 部分
		childernFormat := "div.container.article-container > div > div.col > article > div.article__byline.byline > div > div.col > div > div > div > a:nth-child(%d)"
	singleChildSelector := "div.container.article-container > div > div.col > article > div.article__byline.byline > div > div.col > div > div > div > a"

    index := 1
	for {
		currentAuthor := ""
		currentChildern := fmt.Sprintf(childernFormat,index)
		document.Find(currentChildern).Each(func(i int, s *goquery.Selection) {
			currentAuthor = strings.TrimSpace(s.Text())
			// s.Find("b")
		})
		if  currentAuthor != "" {
			if index != 1 {
				author = author + " & "
			}
			author = author + currentAuthor
			index = index + 1
			continue;
		}
		if currentAuthor == "" {
			break;
		}
		
	}

	if author == "" {
		document.Find(singleChildSelector).Each(func(i int, s *goquery.Selection) {
			author = strings.TrimSpace(s.Text())
		})
	}

    // news part

	return author, published_at
}

// site url https://www.nbcnews.com/world
// rss  http://feeds.nbcnews.com/feeds/worldnews
func (t *Template) ThemoscowtimesScrapMetaData(document *goquery.Document) (string, string) {

	author := ""
	published_at := ""
    // optionion part
	author,published_at = t.ThemoscowtimesScrapMetaDataOpinionPart(document);
	if author == "" {
		author,published_at = t.ThemoscowtimesScrapMetaDataNewsPart(document);
	}


    // news part

	return author, published_at
}