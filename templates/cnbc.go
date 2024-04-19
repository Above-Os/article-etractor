package templates

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type CnbcMetaData struct {
	Context     string `json:"@context"`
	Type        string `json:"@type"`
	Description string `json:"description"`
	Speakable   struct {
		Type        string   `json:"@type"`
		Xpath       []string `json:"xpath"`
		CSSSelector []string `json:"cssSelector"`
	} `json:"speakable"`
	MainEntityOfPage string `json:"mainEntityOfPage"`
	URL              string `json:"url"`
	Headline         string `json:"headline"`
	DateCreated      string `json:"dateCreated"`
	DatePublished    string `json:"datePublished"`
	DateModified     string `json:"dateModified"`
	Author           []struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"author"`
	Publisher struct {
		Type         string `json:"@type"`
		Name         string `json:"name"`
		URL          string `json:"url"`
		FoundingDate string `json:"foundingDate"`
		Logo         struct {
			Type   string `json:"@type"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"logo"`
		SameAs []string `json:"sameAs"`
	} `json:"publisher"`
	ArticleSection string   `json:"articleSection"`
	Keywords       []string `json:"keywords"`
	Image          struct {
		Type   string `json:"@type"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"image"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

type CnbcMetaDataSecond struct {
	Context     string `json:"@context"`
	Type        string `json:"@type"`
	Description string `json:"description"`
	Speakable   struct {
		Type        string   `json:"@type"`
		Xpath       []string `json:"xpath"`
		CSSSelector []string `json:"cssSelector"`
	} `json:"speakable"`
	MainEntityOfPage string   `json:"mainEntityOfPage"`
	URL              string   `json:"url"`
	Headline         string   `json:"headline"`
	DateCreated      string   `json:"dateCreated"`
	DatePublished    string   `json:"datePublished"`
	DateModified     string   `json:"dateModified"`
	Author           []string `json:"author"`
	Publisher        struct {
		Type         string `json:"@type"`
		Name         string `json:"name"`
		URL          string `json:"url"`
		FoundingDate string `json:"foundingDate"`
		Logo         struct {
			Type   string `json:"@type"`
			URL    string `json:"url"`
			Width  int    `json:"width"`
			Height int    `json:"height"`
		} `json:"logo"`
		SameAs []string `json:"sameAs"`
	} `json:"publisher"`
	ArticleSection string   `json:"articleSection"`
	Keywords       []string `json:"keywords"`
	Image          struct {
		Type   string `json:"@type"`
		URL    string `json:"url"`
		Width  int    `json:"width"`
		Height int    `json:"height"`
	} `json:"image"`
	ThumbnailURL string `json:"thumbnailUrl"`
}

func (t *Template) CNBCScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("div#RegularArticle-RelatedQuotes,div[data-test=PlayButton],div.InlineVideo-videoFooter,div.InlineImage-imageEmbedCaption,div.InlineImage-imageEmbedCredit,div.RelatedContent-relatedContent").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})

	document.Find("div.InlineImage-imageContainer,div.RenderKeyPoints-list,div.ArticleBody-articleBody").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})

	return contents
}

func (t *Template) CnbcScrapMetaData(document *goquery.Document) (string, string) {

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
			var firstTypeMetaData CnbcMetaData
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)

			} else {
				for _, currentAuthor := range firstTypeMetaData.Author {
					if len(currentAuthor.Name) != 0 {
						if len(author) != 0 {
							author = author + " & " + currentAuthor.Name
						} else {
							author = currentAuthor.Name
						}
					}
				}
			}

			if len(author) != 0 {
				return
			}

			var secondTypeMetaData CnbcMetaDataSecond
			unmarshalErr = json.Unmarshal([]byte(scriptContent), &secondTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)

			} else {
				for _, currentAuthor := range secondTypeMetaData.Author {
					if len(currentAuthor) != 0 {
						if len(author) != 0 {
							author = author + " & " + currentAuthor
						} else {
							author = currentAuthor
						}
					}
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

func (t *Template) CnbcPublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

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
			var firstTypeMetaData CnbcMetaData
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert CnbcPublishedAtTimeFromScriptMetadata unmarshalError %v", unmarshalErr)
				return

			}
			//publishedAt = firstTypeMetaData[0].DatePublished.Unix()
			fmt.Println(firstTypeMetaData.DatePublished)
			convertPublishedAt, parsePublishedAtErr := parseCnbcTimestamp(firstTypeMetaData.DatePublished)
			if parsePublishedAtErr != nil {
				log.Printf("convert CnbcTimestamp str format to timestamp %v", unmarshalErr)

			} else {
				publishedAt = convertPublishedAt
			}

			if publishedAt != 0 {
				return
			}

			var secondTypeMetaData CnbcMetaDataSecond
			unmarshalErr = json.Unmarshal([]byte(scriptContent), &secondTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert SkyNewsScrap unmarshalError %v", unmarshalErr)

			} else {
				convertPublishedAt, parsePublishedAtErr := parseCnbcTimestamp(secondTypeMetaData.DatePublished)
				if parsePublishedAtErr != nil {
					log.Printf("convert CnbcTimestamp str format to timestamp %v", unmarshalErr)

				} else {
					publishedAt = convertPublishedAt
				}
				if publishedAt != 0 {
					return
				}
				convertPublishedAt, parsePublishedAtErr = parseCnbcTimestampSecond(secondTypeMetaData.DatePublished)
				if parsePublishedAtErr != nil {
					log.Printf("convert CnbcTimestamp str format to timestamp %v", unmarshalErr)

				} else {
					publishedAt = convertPublishedAt
				}

			}
		})

	}
	return publishedAt
}

func parseCnbcTimestamp(timeStr string) (int64, error) {
	// Define the custom layout matching the input string format.
	// Note that "2006-01-02T15:04:05-0700" is the reference time format needed by Go.
	const layout = "2006-01-02T15:04:05-0700"

	// Parse the time string using the custom layout
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return 0, err
	}

	// Return the Unix timestamp (seconds since January 1, 1970 UTC)
	return t.Unix(), nil
}

func parseCnbcTimestampSecond(timeStr string) (int64, error) {
	// Adjusted layout for parsing the time string without a colon in the timezone offset
	const layout = "2006-01-02T15:04:05Z0700"

	// Parse the time string using the adjusted layout
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return 0, err
	}

	// Return the Unix timestamp (seconds since January 1, 1970 UTC)
	return t.Unix(), nil
}
