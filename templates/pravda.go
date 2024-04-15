package templates

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// site url  https://www.pravda.com.ua/eng/news/
// rss  https://www.pravda.com.ua/eng/rss/view_news/
func (t *Template) PravdaScrapMetaData(document *goquery.Document) (string, string) {
	author := ""
	published_at := ""
	cssSelectorFirst := "body > div.main_content > div > div.container_sub_post_news > article > header > div > span > a"
	// cssSelectorSecond := "#root > div.sc-dkjKgF.Oldwm > div.sc-dsLnd.eEVWYF > div.sc-Fiojb.ofioZ.sc-eVAmPc.iJgGA-d.container-fluid > section:nth-child(2) > div > article > div > header > div.sc-fmqvdQ.bGNhbS.sc-fkekHa.kfthnc.author-details > span > a > span"

	cssSelectorList := make([]string, 100)
	cssSelectorList = append(cssSelectorList, cssSelectorFirst)
	// cssSelectorList = append(cssSelectorList, cssSelectorSecond)

	for _, cssSelector := range cssSelectorList {
		document.Find(cssSelector).Each(func(i int, s *goquery.Selection) {
			author = strings.TrimSpace(s.Text())
		})
		if author != "" {
			break
		}
	}

	return author, published_at
}

func (t *Template) PravdaScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("div.advtext_mob").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	document.Find("div.post_photo_news,div.post_text").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
