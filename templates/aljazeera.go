package templates

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type AljazeeraImageList struct {
	ImageList []Image `json:"image"`
}

func (t *Template) AljazeeraScrapMetaData(document *goquery.Document) (string, string) {
	author := ""
	published_at := ""
	author, published_at = t.AuthorExtractFromScriptMetadata(document)

	return author, published_at
}

func (t *Template) AljazeeraScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("div.more-on").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	/*currentCoverImageUrl := t.AljazeeraImageExtractFromScriptMetadata(document)
	if currentCoverImageUrl != "" {
		currentImageTag := fmt.Sprintf("<figure><img src=\"%s\"/></figure>",currentCoverImageUrl)
		contents = contents + currentImageTag
	}*/

	document.Find("figure.article-featured-image,figure.gallery-featured-image,div.wysiwyg--all-content,figure.gallery-image").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}

func (t *Template) AljazeeraImageExtractFromScriptMetadata(document *goquery.Document) string {

	currentImage := ""

	scriptSelectorFirst := "head > script[type=\"application/ld+json\"]"
	scriptSelectorSecond := "body > script[type=\"application/ld+json\"]"
	scriptSelectorList := make([]string, 100)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorFirst)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorSecond)

	for _, scriptSelector := range scriptSelectorList {

		document.Find(scriptSelector).Each(func(i int, s *goquery.Selection) {
			if currentImage != "" {
				return
			}
			scriptContent := strings.TrimSpace(s.Text())

			var currentImageList AljazeeraImageList
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &currentImageList)
			if unmarshalErr != nil {
				log.Printf("unmarshal AljazeeraImageList error")
				return
			}
			if len(currentImageList.ImageList) > 0 {
				currentImage = currentImageList.ImageList[0].Url
			}

			// var jsonMap map[string]interface{}

			// logger.Info("script content %s  author length %d",scriptContent, len(currentDWMetadata.Author))
		})
		if currentImage != "" {
			break
		}
	}

	return currentImage
}
