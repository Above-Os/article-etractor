package templates

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (t *Template) CnbcScrapMetaData(document *goquery.Document) (string, string) {
	author := ""
	published_at := ""
	author, published_at = t.AuthorExtractFromScriptMetadata(document)

	if author == "" {
		document.Find("a.Author-authorName").Each(func(i int, s *goquery.Selection) {
			author = strings.TrimSpace(s.Text())

		})
	}

	return author, published_at
}

func (t *Template) CNBCScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("div#RegularArticle-RelatedQuotes,div[data-test=PlayButton],div.InlineVideo-videoFooter,div.InlineImage-imageEmbedCaption,div.InlineImage-imageEmbedCredit").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	/*currentCoverImageUrl := t.CnbcConverImageUrlExtractFromScriptMetadata(document)
	if currentCoverImageUrl != "" {
		currentImageTag := fmt.Sprintf("<figure><img src=\"%s\"/></figure>",currentCoverImageUrl)
		contents = contents + currentImageTag
	}*/
	// log.Printf("current image url %s",currentCoverImageUrl)
	document.Find("div.RenderKeyPoints-list,div.ArticleBody-articleBody").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}

type CNBCCoverImage struct {
	ThumbnailCovertImageUrl string `json:"thumbnailUrl"`
}

func (t *Template) CnbcConverImageUrlExtractFromScriptMetadata(document *goquery.Document) string {
	url := ""
	scriptSelectorFirst := "head > script[type=\"application/ld+json\"]"
	scriptSelectorSecond := "body > script[type=\"application/ld+json\"]"
	// #caas-art-51cf82b6-f0e5-3999-a910-ce4fb658efb4 > article > script:nth-child(1)
	scriptSelectorList := make([]string, 100)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorFirst)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorSecond)

	for _, scriptSelector := range scriptSelectorList {

		document.Find(scriptSelector).Each(func(i int, s *goquery.Selection) {
			if url != "" {
				return
			}
			scriptContent := strings.TrimSpace(s.Text())

			var currentCNBCCovertImage CNBCCoverImage
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &currentCNBCCovertImage)
			if unmarshalErr != nil {
				log.Println("unmarshal error")
				return
			}

			url = currentCNBCCovertImage.ThumbnailCovertImageUrl
			// logger.Info("script content %s  author length %d",scriptContent, len(currentDWMetadata.Author))
		})
		if url != "" {
			break
		}
	}

	return url
}

func (t *Template) CnbcWorldGetPublishedAtTimestampSingleJson(document *goquery.Document) int64 {

	var publishedAtTimestamp int64 = 0
	scriptSelectorFirst := "head > script[type=\"application/ld+json\"]"
	scriptSelectorSecond := "body > script[type=\"application/ld+json\"]"
	scriptSelectorThird := "script[type=\"application/ld+json\"]"
	scriptSelectorList := make([]string, 100)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorFirst)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorSecond)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorThird)

	for _, scriptSelector := range scriptSelectorList {

		document.Find(scriptSelector).Each(func(i int, s *goquery.Selection) {
			if publishedAtTimestamp != 0 {
				return
			}
			scriptContent := strings.TrimSpace(s.Text())

			var jsonMap map[string]interface{}
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &jsonMap)
			if unmarshalErr != nil {
				log.Printf("unmarshal error")
				return
			}
			currentPublishedAt, ok := jsonMap["datePublished"]
			if !ok {

				return
			}
			currentPublishedAtStr := currentPublishedAt.(string)
			log.Printf("currentPublishedAtStr %s", currentPublishedAtStr[0:len(currentPublishedAtStr)-5]+"Z")
			publishedAtTimestamp = ConvertStringTimeToTimestamp(currentPublishedAtStr[0:len(currentPublishedAtStr)-5] + "Z")

			// var jsonMap map[string]interface{}

			// logger.Info("script content %s  author length %d",scriptContent, len(currentDWMetadata.Author))
		})
		if publishedAtTimestamp != 0 {
			break
		}
	}

	return publishedAtTimestamp

}
