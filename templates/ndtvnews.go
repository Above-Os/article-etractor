package templates

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// site url  https://www.ndtv.com/world-news/
// rss  https://feeds.feedburner.com/ndtvnews-world-news
func (t *Template) NdtvNewsScrapMetaData(document *goquery.Document) (string, string) {

	author := ""
	published_at := ""
	cssSelectorFirst := "body > div.content > div > div > section > div > div.sp-hd > nav > div > span:nth-child(2) > a > span > span"
	cssSelectorSecond := "body > div.content > div > div > section > div > div.sp-hd > nav > div > span:nth-child(3) > a > span > span"

	cssSelectorThird := "body > div.content > div > div > section > div > div.sp-hd > nav > div > span:nth-child(4) > a > span > span"
	cssSelectorList := make([]string, 100)
	cssSelectorList = append(cssSelectorList, cssSelectorFirst)
	cssSelectorList = append(cssSelectorList, cssSelectorSecond)
	cssSelectorList = append(cssSelectorList, cssSelectorThird)

	for _, cssSelector := range cssSelectorList {
		document.Find(cssSelector).Each(func(i int, s *goquery.Selection) {
			itempropValue, itempropExist := s.Attr("itemprop")
			if itempropExist {
				if itempropValue == "name" {
					author = strings.TrimSpace(s.Text())
				}
			}

		})
		if author != "" {
			break
		}
	}

	return author, published_at
}

func (t *Template) NdtvGetPublishedAtTimestamp(document *goquery.Document) int64 {
	var  publishedAtTimeStamp  int64 = 0
	cssSelector := "meta[itemprop=\"datePublished\"]"
	document.Find(cssSelector).Each(func(i int, s *goquery.Selection) {
		itempropValue, itempropExist := s.Attr("content")
		if itempropExist {
			// logger.Info("itempropValue %s",itempropValue)
			publishedAtTimeStamp = ConvertStringTimeToTimestampRFC33399(itempropValue)
		}

	})
	return publishedAtTimeStamp
}