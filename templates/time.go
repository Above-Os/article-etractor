package templates

import (

	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (t *Template) TimeScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("p").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if strings.HasPrefix(text, "Read More:") {
			RemoveNodes(s)
		}

	})

	document.Find("div.featured-media,div#article-body-main").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}


func (t *Template) TimesScrapMetaData(document *goquery.Document) (string, string) {
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
	// scriptSelector := "body > script:nth-child(2)"
	author, published_at=t.AuthorExtractFromListScriptMetadata(document)

	


	return author, published_at
}