package templates

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (t *Template) DWScrapMetaDataSingleChild(document *goquery.Document) (string, string) {
	author := ""
	published_at := ""
	cssSelectorFirst := "#root > div.sc-dkjKgF.Oldwm > div.sc-dsLnd.eEVWYF > div.sc-Fiojb.ofioZ.sc-eVAmPc.iJgGA-d.container-fluid > section:nth-child(2) > div > article > div > header > div.sc-fmqvdQ.bGNhbS.sc-fkekHa.kfthnc.author-details > span > span"
	cssSelectorSecond := "#root > div.sc-dkjKgF.Oldwm > div.sc-dsLnd.eEVWYF > div.sc-Fiojb.ofioZ.sc-eVAmPc.iJgGA-d.container-fluid > section:nth-child(2) > div > article > div > header > div.sc-fmqvdQ.bGNhbS.sc-fkekHa.kfthnc.author-details > span > a > span"

	cssSelectorList := make([]string, 100)
	cssSelectorList = append(cssSelectorList, cssSelectorFirst)
	cssSelectorList = append(cssSelectorList, cssSelectorSecond)

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

// site url  http://www.dw.com/english/?maca=en-rss-en-world-4025-rdf
// rss  https://rss.dw.com/rdf/rss-en-world
func (t *Template) DWScrapMetaData(document *goquery.Document) (string, string) {
	//xpath /html/body/div[2]/div[1]/div[2]/div[2]/section[1]/div/article/div/header/div[2]/span
	//selector #root > div.sc-dkjKgF.Oldwm > div.sc-dsLnd.eEVWYF > div.sc-Fiojb.ofioZ.sc-eVAmPc.iJgGA-d.container-fluid > section:nth-child(2) > div > article > div > header > div.sc-fmqvdQ.bGNhbS.sc-fkekHa.kfthnc.author-details > span > span
	// #root > div.sc-dkjKgF.Oldwm > div.sc-dsLnd.eEVWYF > div.sc-Fiojb.ofioZ.sc-eVAmPc.iJgGA-d.container-fluid > section:nth-child(2) > div > article > div > header > div.sc-fmqvdQ.bGNhbS.sc-fkekHa.kfthnc.author-details > span > span
	//         #root > div.sc-dkjKgF.Oldwm > div.sc-dsLnd.eEVWYF > div.sc-Fiojb.ofioZ.sc-eVAmPc.iJgGA-d.container-fluid > section:nth-child(2) > div > article > div > header > div.sc-fmqvdQ.bGNhbS.sc-fkekHa.kfthnc.author-details > span > a > span
	// #root > div.sc-dkjKgF.Oldwm > div.sc-dsLnd.eEVWYF > div.sc-Fiojb.ofioZ.sc-eVAmPc.iJgGA-d.container-fluid > section:nth-child(2) > div > article > div > header > div.sc-fmqvdQ.bGNhbS.sc-fkekHa.kfthnc.author-details > span:nth-child(1) > span
	// #root > div.sc-dkjKgF.Oldwm > div > div.sc-fHIGvW.jTgZVM.sc-dkCBEl.ftkcGG.container-fluid > section:nth-child(2) > div > article > div > header > div.sc-ihxOxR.exeCrb.sc-eVRJAz.gAHmDO.author-details > span > span

	// head > script:nth-child(127)
	// head > script:nth-child(126)
	author := ""
	published_at := ""
	author, published_at = t.AuthorExtractFromScriptMetadata(document)

	return author, published_at
}

func (t *Template) DWScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("header>div,header>h1,header>span,section[data-tracking-name=sharing-icons-inline],div.advertisement,footer").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	document.Find("article").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	/*document.Find("header>p,div.rich-text").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})*/
	return contents
}
