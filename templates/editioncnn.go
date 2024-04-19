package templates

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func extractEditionCnnValues(text string) (names []string, dates []string) {
	nameRegex := `"name"\s*:\s*"([^"]+)"`
	dateRegex := `"dateCreated"\s*:\s*"([^"]+)"`

	namePattern := regexp.MustCompile(nameRegex)
	datePattern := regexp.MustCompile(dateRegex)

	nameMatches := namePattern.FindAllStringSubmatch(text, -1)
	dateMatches := datePattern.FindAllStringSubmatch(text, -1)

	for _, match := range nameMatches {
		names = append(names, match[1])
	}
	for _, match := range dateMatches {
		dates = append(dates, match[1])
	}

	return names, dates
}

func (t *Template) EditionCnnScrapMetaData(document *goquery.Document) (string, string) {

	author := ""
	published_at := ""
	scriptSelectorFirst := "head > script[type=\"application/ld+json\"]"
	scriptSelectorSecond := "body > script[type=\"application/ld+json\"]"
	scriptSelectorThird := "script[type=\"application/ld+json\"]"

	scriptSelectorList := make([]string, 100)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorFirst)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorSecond)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorThird)

	for _, scriptSelector := range scriptSelectorList {

		document.Find(scriptSelector).Each(func(i int, s *goquery.Selection) {
			if author != "" {
				return
			}
			scriptContent := strings.TrimSpace(s.Text())
			nameList, _ := extractEditionCnnValues(scriptContent)
			for _, currentName := range nameList {
				if len(author) != 0 {
					author = author + " & " + currentName
				} else {
					author = currentName
				}
			}

		})
		if author != "" {
			break
		}
	}
	log.Printf("author last: %s", author)
	return author, published_at
}

func (t *Template) EditionCnnPublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

	var publishedAt int64 = 0

	scriptSelectorFirst := "head > script[type=\"application/ld+json\"]"
	scriptSelectorSecond := "body > script[type=\"application/ld+json\"]"
	scriptSelectorThird := "script[type=\"application/ld+json\"]"

	scriptSelectorList := make([]string, 100)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorFirst)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorSecond)
	scriptSelectorList = append(scriptSelectorList, scriptSelectorThird)

	for _, scriptSelector := range scriptSelectorList {

		document.Find(scriptSelector).Each(func(i int, s *goquery.Selection) {
			if publishedAt != 0 {
				return
			}
			scriptContent := strings.TrimSpace(s.Text())
			_, timeList := extractEditionCnnValues(scriptContent)
			fmt.Println(timeList)
			//publishedAt = firstTypeMetaData[0].DatePublished.Unix()
		})

	}
	return publishedAt
}
