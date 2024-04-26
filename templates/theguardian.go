package templates

import (
	"encoding/json"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type TheGuardianCoverImage struct {
	TheImageUrlList []string `json:"image"`
}

func (t *Template) TheguardianScrapContent(document *goquery.Document) string {
	contents := ""

	document.Find("gu-island,aside,a.dcr-porppu,p.dcr-porppu,p#EmailSignup-skip-link-8").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	document.Find("a").Each(func(i int, s *goquery.Selection) {
		if s.Text() == "skip past newsletter promotion" {
			RemoveNodes(s)
		}

	})
	/*currentCoverImageUrl := t.TheGuardianImageExtractFromListScriptMetadata(document)
	// log.Printf("current the guardian image %s",currentCoverImageUrl)
	if currentCoverImageUrl != "" {
		currentImageTag := fmt.Sprintf("<figure><img src=\"%s\"/></figure>",currentCoverImageUrl)
		contents = contents + currentImageTag
	}*/

	document.Find("div#img-1,div#maincontent").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}

func (t *Template) TheguardianScrapMetaData(document *goquery.Document) (string, string) {

	author := ""
	published_at := ""

	author, published_at = t.AuthorExtractFromListScriptMetadata(document)

	return author, published_at
}

func (t *Template) TheGuardianImageExtractFromListScriptMetadata(document *goquery.Document) string {

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
			// log.Printf("script content %s ",scriptContent)
			jsonMap := []map[string]interface{}{}

			unmarshalErr := json.Unmarshal([]byte(scriptContent), &jsonMap)
			if unmarshalErr != nil {
				log.Printf("unmarshal error  %v", unmarshalErr)
				return
			}

			for _, currentJsonMap := range jsonMap {

				currentImageInterface, ok := currentJsonMap["image"]
				if ok {
					currentImageList := make([]string, 0)
					for _, val := range currentImageInterface.([]interface{}) {
						currentImageList = append(currentImageList, val.(string))
					}
					if len(currentImageList) > 0 {
						currentImage = currentImageList[0]
						break
					}

				}
			}

		})
		if currentImage != "" {
			break
		}
	}

	return currentImage
}
