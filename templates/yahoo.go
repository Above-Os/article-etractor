package templates

import (
	"encoding/json"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (t *Template) YahooNewsScrapMetaData(document *goquery.Document) (string, string) {

	author := ""
	published_at := ""
	author, published_at = t.AuthorExtractFromScriptMetadata(document)
	if author == "" {
		document.Find("span.caas-author-byline-collapse").Each(func(i int, s *goquery.Selection) {
			reg := regexp.MustCompile(`(?:\n\s+)`)
			author = reg.ReplaceAllString(s.Text(), "")

		})
		document.Find("div.caas-attr-time-style>time").Each(func(i int, s *goquery.Selection) {
			published_at, _ = s.Attr("datetime")
		})
	}

	return author, published_at
}

func (t *Template) YahoocrapContent(document *goquery.Document) string {

	contents := ""
	//aside.caas-aside-section https://nypost.com/2024/04/08/business/tsmc-to-boost-computer-chip-production-in-arizona-with-11-6-billion-in-federal-grants-loans/
	document.Find("header.caas-header,div.caas-content-byline-wrapper,button,div.xray-error-wrapper,div.caas-xray-pills-container,aside.caas-aside-section").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	/*currentCoverImageUrl := t.YahooConverImageUrlExtractFromScriptMetadata(document)
	// log.Printf("current yahoo covert image url %s",currentCoverImageUrl)
	if currentCoverImageUrl != "" {
		currentImageTag := fmt.Sprintf("<figure><img src=\"%s\"/></figure>",currentCoverImageUrl)
		contents = contents + currentImageTag
	}*/
	document.Find("div#module-article").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}

type YahooCoverImage struct {
	ThumbnailCovertImageUrl string `json:"thumbnailUrl"`
}

func (t *Template) YahooConverImageUrlExtractFromScriptMetadata(document *goquery.Document) string {
	url := ""
	// scriptSelectorFirst := "head > script[type=\"application/ld+json\"]"
	// scriptSelectorSecond := "body > script[type=\"application/ld+json\"]"
	scriptSelectorThird := "article > script[type=\"application/ld+json\"]"

	// #caas-art-51cf82b6-f0e5-3999-a910-ce4fb658efb4 > article > script:nth-child(1)
	scriptSelectorList := make([]string, 100)
	// scriptSelectorList = append(scriptSelectorList, scriptSelectorFirst)
	// scriptSelectorList = append(scriptSelectorList, scriptSelectorSecond)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorThird)

	for _, scriptSelector := range scriptSelectorList {

		document.Find(scriptSelector).Each(func(i int, s *goquery.Selection) {
			if url != "" {
				return
			}
			scriptContent := strings.TrimSpace(s.Text())

			var currentYahooImage YahooCoverImage
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &currentYahooImage)
			if unmarshalErr != nil {
				log.Printf("unmarshal error")
				return
			}

			urlList := strings.Split(currentYahooImage.ThumbnailCovertImageUrl, "-/")
			if len(urlList) == 3 {
				url = urlList[2]
			}

			// logger.Info("script content %s  author length %d",scriptContent, len(currentDWMetadata.Author))
		})
		if url != "" {
			break
		}
	}

	return url
}
