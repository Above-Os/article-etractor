package templates

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type A16ZMetaDataWithAuthorName struct {
	Context          string `json:"@context"`
	Type             string `json:"@type"`
	Headline         string `json:"headline"`
	URL              string `json:"url"`
	MainEntityOfPage struct {
		Type string `json:"@type"`
		ID   string `json:"@id"`
	} `json:"mainEntityOfPage"`
	ThumbnailURL string `json:"thumbnailUrl"`
	Image        struct {
		Type string `json:"@type"`
		URL  string `json:"url"`
	} `json:"image"`
	ArticleSection string `json:"articleSection"`
	Author         []struct {
		Type string `json:"@type"`
		Name string `json:"name"`
	} `json:"author"`
	Creator   []string `json:"creator"`
	Publisher struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		Logo string `json:"logo"`
	} `json:"publisher"`
	Keywords      []any  `json:"keywords"`
	DateCreated   string `json:"dateCreated"`
	DatePublished string `json:"datePublished"`
	DateModified  string `json:"dateModified"`
}

type A16ZMetadataWithAuthorID struct {
	Context string `json:"@context"`
	Graph   []struct {
		Type     string `json:"@type"`
		ID       string `json:"@id"`
		IsPartOf struct {
			ID string `json:"@id"`
		} `json:"isPartOf,omitempty"`
		Author []struct {
			ID string `json:"@id"`
		} `json:"author,omitempty"`
		Headline         string `json:"headline,omitempty"`
		DatePublished    string `json:"datePublished,omitempty"`
		DateModified     string `json:"dateModified,omitempty"`
		MainEntityOfPage struct {
			ID string `json:"@id"`
		} `json:"mainEntityOfPage,omitempty"`
		WordCount    int `json:"wordCount,omitempty"`
		CommentCount int `json:"commentCount,omitempty"`
		Publisher    struct {
			ID string `json:"@id"`
		} `json:"publisher,omitempty"`
		ArticleSection  []string `json:"articleSection,omitempty"`
		InLanguage      string   `json:"inLanguage,omitempty"`
		PotentialAction []struct {
			Type   string   `json:"@type"`
			Name   string   `json:"name"`
			Target []string `json:"target"`
		} `json:"potentialAction,omitempty"`
		URL         string `json:"url,omitempty"`
		Name        string `json:"name,omitempty"`
		Description string `json:"description,omitempty"`
		Logo        struct {
			Type       string `json:"@type"`
			InLanguage string `json:"inLanguage"`
			ID         string `json:"@id"`
			URL        string `json:"url"`
			ContentURL string `json:"contentUrl"`
			Width      int    `json:"width"`
			Height     int    `json:"height"`
			Caption    string `json:"caption"`
		} `json:"logo,omitempty"`
		Image struct {
			ID string `json:"@id"`
		} `json:"image,omitempty"`
		SameAs []string `json:"sameAs,omitempty"`
		Image0 struct {
			Type       string `json:"@type"`
			InLanguage string `json:"inLanguage"`
			ID         string `json:"@id"`
			URL        string `json:"url"`
			ContentURL string `json:"contentUrl"`
			Width      int    `json:"width"`
			Height     int    `json:"height"`
			Caption    string `json:"caption"`
		} `json:"image,omitempty"`
	} `json:"@graph"`
}

func (t *Template) A16ZScrapMetaData(document *goquery.Document) (string, string) {

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
			var firstTypeMetaData A16ZMetaDataWithAuthorName
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert a16metadata  unmarshalError %v", unmarshalErr)
			} else {
				for _, currentAuthro := range firstTypeMetaData.Author {
					if len(currentAuthro.Name) != 0 {
						if len(author) != 0 {
							author = author + " & " + currentAuthro.Name
						} else {
							author = currentAuthro.Name
						}
					}
				}

			}

		})
		if author != "" {

			break
		}
	}
	if len(author) == 0 {
		author = "a16z editorial"
	}
	log.Printf("author last: %s", author)
	return author, published_at
}

func (t *Template) A16ZPublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

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
			var firstTypeMetaData A16ZMetaDataWithAuthorName
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert a16zmetadata unmarshalError %v", unmarshalErr)
				return

			} else {
				currentParsedPublishedAt, parsePublishedErr := ConvertToTimestampA16Z(firstTypeMetaData.DatePublished)
				if parsePublishedErr != nil {
					log.Printf("convert time fail")
					return
				}
				publishedAt = currentParsedPublishedAt
				// publishedAt = firstTypeMetaData.DateCreated.Unix()
			}

			if publishedAt != 0 {
				return
			}

			var secondTypeMetaData A16ZMetadataWithAuthorID
			unmarshalErr = json.Unmarshal([]byte(scriptContent), &secondTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert a16zmetadata unmarshalError %v", unmarshalErr)
				return

			} else {
				for _, currentGraph := range secondTypeMetaData.Graph {
					if publishedAt != 0 {
						break
					}
					currentParsedPublishedAt, parsePublishedErr := ConvertToTimestampA16Z(currentGraph.DatePublished)
					if parsePublishedErr != nil {
						log.Printf("convert time fail")
						return
					}
					publishedAt = currentParsedPublishedAt
				}
				//publishedAt = firstTypeMetaData.DateCreated.Unix()
			}
			//publishedAt = firstTypeMetaData[0].DatePublished.Unix()
		})

	}
	return publishedAt
}

func ConvertToTimestampA16Z(timeStr string) (int64, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}
