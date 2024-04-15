package templates

import (
	"github.com/PuerkitoBio/goquery"
)

/*
<div class="QHbl nkdH hTos whbO " data-testid="prism-byline">

	<div class="VZTD mLAS OcxM oJce ">
	    <div class="kKfX VZTD rEPu ">
	        <div class="TQPv HUca ucZk WxHI HhZO yaUf VOJB XSba Umfi ukdD ">
	            <span class="tChG zbFa ">By</span>
	            <span>GERALD IMRAY Associated Press</span>
	        </div>
	        <div class="VZTD mLAS ">
	            <div class="xAPp Zdbe  jTKb pCRh ">January 14, 2024, 7:12 AM</div>
	        </div>
	    </div>
	</div>

</div>
*/
func (t *Template) AbcNewsScrapMetaData(document *goquery.Document) (string, string) {
	author := ""
	published_at := ""
	document.Find("[data-testid='prism-byline']").Each(func(i int, s *goquery.Selection) {
		s.Find("div:nth-child(1)>div").Each(func(i int, firstDiv *goquery.Selection) {
			firstDiv.Find("div").Each(func(subindex int, subDiv *goquery.Selection) {
				if subindex == 0 {
					subDiv.Find("span:nth-child(2)").Each(func(i int, authSpan *goquery.Selection) {
						author = authSpan.Text()
					})
				}
				/*if subindex == 1 {
					subDiv.Find("div").Each(func(i int, timeDiv *goquery.Selection) {
						published_at = timeDiv.Text()
					})
				}*/
			})
		})

	})

	return author, published_at
}

func (t *Template) AbcNewsScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("div[data-testid=prism-byline],div[data-testid=prism-headline]>h1,div[data-testid=prism-tags]").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	document.Find("p").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "___" {
			RemoveNodes(s)
		}

	})
	document.Find("div.FITT_Article_main__body").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
