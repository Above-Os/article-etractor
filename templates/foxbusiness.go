package templates

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type FoxbusinessMetadata struct {
	Context string `json:"@context"`
	Graph   []struct {
		Context string `json:"@context"`
		Type    string `json:"@type"`
		ID      string `json:"@id"`
		Name    string `json:"name"`
		URL     string `json:"url"`
	} `json:"@graph,omitempty"`
	Type            string `json:"@type,omitempty"`
	ID              string `json:"@id,omitempty"`
	ItemListElement []struct {
		Type     string `json:"@type"`
		Position int    `json:"position"`
		Item     struct {
			ID   string `json:"@id"`
			Name string `json:"name"`
		} `json:"item"`
	} `json:"itemListElement,omitempty"`
	URL              string    `json:"url,omitempty"`
	Headline         string    `json:"headline,omitempty"`
	MainEntityOfPage string    `json:"mainEntityOfPage,omitempty"`
	DatePublished    time.Time `json:"datePublished,omitempty"`
	DateModified     time.Time `json:"dateModified,omitempty"`
	Description      string    `json:"description,omitempty"`
	ArticleSection   string    `json:"articleSection,omitempty"`
	ArticleBody      string    `json:"articleBody,omitempty"`
	Keywords         string    `json:"keywords,omitempty"`
	Name             string    `json:"name,omitempty"`
	ThumbnailURL     string    `json:"thumbnailUrl,omitempty"`
	WordCount        string    `json:"wordCount,omitempty"`
	TimeRequired     string    `json:"timeRequired,omitempty"`
	MainEntity       struct {
		Type string `json:"@type"`
		ID   string `json:"@id"`
	} `json:"mainEntity,omitempty"`
	Author struct {
		Type        string   `json:"@type"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		URL         string   `json:"url"`
		SameAs      []string `json:"sameAs"`
		Image       struct {
			Type   string `json:"@type"`
			URL    string `json:"url"`
			Height int    `json:"height"`
			Width  int    `json:"width"`
		} `json:"image"`
	} `json:"author,omitempty"`
	Editor struct {
		Type        string   `json:"@type"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		URL         string   `json:"url"`
		SameAs      []string `json:"sameAs"`
		Image       struct {
			Type   string `json:"@type"`
			URL    string `json:"url"`
			Height int    `json:"height"`
			Width  int    `json:"width"`
		} `json:"image"`
	} `json:"editor,omitempty"`
	Publisher struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		URL  string `json:"url"`
		Logo struct {
			Type   string `json:"@type"`
			URL    string `json:"url"`
			Width  string `json:"width"`
			Height string `json:"height"`
		} `json:"logo"`
	} `json:"publisher,omitempty"`
}

type FoxbusinessMetadataSecond struct {
	Context          string `json:"@context"`
	Type             string `json:"@type"`
	MainEntityOfPage string `json:"mainEntityOfPage"`
	Headline         string `json:"headline"`
	ArticleBody      string `json:"articleBody"`
	DatePublished    string `json:"datePublished"`
	DateModified     string `json:"dateModified"`
	Description      string `json:"description"`
	Author           []struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"author"`
	Publisher struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		Logo struct {
			Type   string `json:"@type"`
			URL    string `json:"url"`
			Width  string `json:"width"`
			Height string `json:"height"`
		} `json:"logo"`
	} `json:"publisher"`
	Image struct {
		Type string `json:"@type"`
		URL  string `json:"url"`
	} `json:"image"`
}

func (t *Template) FoxbusinessScrapMetaData(document *goquery.Document) (string, string) {

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
			var firstTypeMetaData FoxbusinessMetadata
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert foxbusiness unmarshalError %v", unmarshalErr)
			} else {
				author = firstTypeMetaData.Author.Name
			}
			if len(author) != 0 {
				return
			}
			var secondTypeMetaData FoxbusinessMetadataSecond
			unmarshalErr = json.Unmarshal([]byte(scriptContent), &secondTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert foxbusiness unmarshalError %v", unmarshalErr)
				return
			}
			for _, currentAuthor := range secondTypeMetaData.Author {
				if len(currentAuthor.Name) != 0 {
					if len(author) != 0 {
						author = author + " & " + currentAuthor.Name
					} else {
						author = currentAuthor.Name
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

func (t *Template) FoxbusinessPublishedAtTimeFromScriptMetadata(document *goquery.Document) int64 {

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
			var firstTypeMetaData FoxbusinessMetadata
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &firstTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert foxbusiness unmarshalError %v", unmarshalErr)

			} else {
				publishedAt = firstTypeMetaData.DatePublished.Unix()
			}
			//publishedAt = firstTypeMetaData[0].DatePublished.Unix()
			if publishedAt != 0 {
				return
			}
			var secondTypeMetaData FoxbusinessMetadataSecond
			unmarshalErr = json.Unmarshal([]byte(scriptContent), &secondTypeMetaData)
			if unmarshalErr != nil {
				log.Printf("convert foxbusiness unmarshalError %v", unmarshalErr)
				return
			}
			parsedPublishedAt, parseErr := ConvertToTimestampFoxbusiness(secondTypeMetaData.DatePublished)
			if parseErr != nil {
				log.Printf("convert foxbusiness error %v", parseErr)
				return
			}
			publishedAt = parsedPublishedAt
		})

	}
	return publishedAt
}

func ConvertToTimestampFoxbusiness(timeStr string) (int64, error) {
	t, err := time.Parse(time.RFC3339, timeStr)
	if err != nil {
		return 0, err
	}

	return t.Unix(), nil
}
