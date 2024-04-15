package templates

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// site url  https://www.euronews.com/news/international
// rss  https://www.euronews.com/rss?format=mrss&level=theme&name=news

/*
{
   "@context":"https://schema.org/",
   "@graph":[
      {
         "@type":"NewsArticle",
         "mainEntityOfPage":{
            "@type":"Webpage",
            "url":"https://www.euronews.com/2024/01/30/spanish-lawmakers-vote-on-a-highly-divisive-law-to-grant-amnesty-to-catalan-separatists"
         },
         "headline":"Spanish lawmakers vote against highly divisive law to grant amnesty to Catalan separatists",
         "description":"The law could have paved the way for the return of fugitive ex-Catalan president Carles Puigdemont.",

         "dateCreated":"2024-01-30 16:55:12",
         "dateModified":"2024-01-30 22:18:09",
         "datePublished":"2024-01-30 17:23:33",
         "image":{
            "@type":"ImageObject",
            "url":"https://static.euronews.com/articles/stories/08/20/64/42/1440x810_cmsv2_75c24ebd-18df-596d-9ddb-9be2c6106621-8206442.jpg",
            "width":"1440px",
            "height":"810px",
            "caption":"Demonstrators waves Spanish flags at Plaza del castillo square during a protest called by Spain's Conservative Popular Party.",
            "thumbnail":"https://static.euronews.com/articles/stories/08/20/64/42/385x202_cmsv2_75c24ebd-18df-596d-9ddb-9be2c6106621-8206442.jpg",
            "publisher":{
               "@type":"Organization",
               "name":"euronews",
               "url":"https://static.euronews.com/website/images/euronews-logo-main-blue-403x60.png"
            }
         },
         "author":{
            "@type":"Person",
            "name":"Euronews",
            "url":"https://www.euronews.com/",
            "sameAs":"https://twitter.com/euronews"
         },
         "publisher":{
            "@type":"Organization",
            "name":"Euronews",
            "legalName":"Euronews",
            "url":"https://www.euronews.com/",
            "logo":{
               "@type":"ImageObject",
               "url":"https://static.euronews.com/website/images/euronews-logo-main-blue-403x60.png",
               "width":"403px",
               "height":"60px"
            },
            "sameAs":[
               "https://www.facebook.com/euronews",
               "https://twitter.com/euronews",
               "https://flipboard.com/@euronews",
               "https://www.instagram.com/euronews.tv/",
               "https://www.linkedin.com/company/euronews"
            ]
         },
         "video":{
            "@type":"VideoObject",
            "contentUrl":"https://video.euronews.com/mp4/med/EN/NW/SU/24/01/30/en/240130_NWSU_54671157_54671183_5280_221219_en.mp4",
            "description":"The law could have paved the way for the return of fugitive ex-Catalan president Carles Puigdemont.",
            "duration":"PT5S",
            "embedUrl":"https://www.euronews.com/embed/2467318",
            "height":"202px",
            "name":"Spain's congress votes amnesty law for Catalan separatists",
            "thumbnailUrl":"https://static.euronews.com/articles/stories/08/20/64/42/385x202_cmsv2_75c24ebd-18df-596d-9ddb-9be2c6106621-8206442.jpg",
            "uploadDate":"2024-01-30 17:23:33",
            "videoQuality":"hd",
            "width":"385px",
            "inLanguage":{
               "name":"en-GB",
               "alternateName":"en",
               "description":"https://www.euronews.com",
               "identifier":"en",
               "url":"https://www.euronews.com",
               "inLanguage":"en-GB"
            },
            "publisher":{
               "@type":"Organization",
               "name":"Euronews",
               "legalName":"Euronews",
               "url":"https://www.euronews.com/",
               "logo":{
                  "@type":"ImageObject",
                  "url":"https://static.euronews.com/website/images/euronews-logo-main-blue-403x60.png",
                  "width":"403px",
                  "height":"60px"
               },
               "sameAs":[
                  "https://www.facebook.com/euronews",
                  "https://twitter.com/euronews",
                  "https://flipboard.com/@euronews",
                  "https://www.instagram.com/euronews.tv/",
                  "https://www.linkedin.com/company/euronews"
               ]
            }
         },
         "speakable":{
            "@type":"SpeakableSpecification",
            "xPath":[
               "/html/head/title",
               "/html/head/meta[@name='description']/@content"
            ],
            "url":"https://www.euronews.com/2024/01/30/spanish-lawmakers-vote-on-a-highly-divisive-law-to-grant-amnesty-to-catalan-separatists"
         }
      },
      {
         "@type":"WebSite",
         "name":"Euronews.com",
         "url":"https://www.euronews.com/",
         "potentialAction":{
            "@type":"SearchAction",
            "target":"https://www.euronews.com/search?query={search_term_string}",
            "query-input":"required name=search_term_string"
         },
         "sameAs":[
            "https://www.facebook.com/euronews",
            "https://twitter.com/euronews",
            "https://flipboard.com/@euronews",
            "https://www.instagram.com/euronews.tv/",
            "https://www.linkedin.com/company/euronews"
         ]
      }
   ]
}
*/

func ConvertInterfaceToArrayOfMaps(value interface{}) ([]map[string]interface{}, error) {
	if valueSlice, ok := value.([]interface{}); ok {
		currentMapSlice := []map[string]interface{}{}
		for _, currentMapInterface := range valueSlice {
			currentMap, currentMapErr := ConvertInterfaceToMap(currentMapInterface)
			if currentMapErr == nil {
				currentMapSlice = append(currentMapSlice, currentMap)
			} else {

			}
		}
		return currentMapSlice, nil
	} else if valueMap, ok := value.(map[string]interface{}); ok {
		valueSlice := make([]map[string]interface{}, 1)
		valueSlice[0] = valueMap
		return valueSlice, nil
	} else {
		return nil, fmt.Errorf("value is not a []map[string]interface{} or map[string]interface{}")
	}
}

func (t *Template) EuroNewsGetPublishedAtTimeStamp(document *goquery.Document) int64 {
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
				log.Printf("EuroNewsGetAuthor unmarshal string to json fail")
				return
			}
			graphJson, ok := jsonMap["@graph"]
			if !ok {
				log.Printf("EuroNewsGetAuthor @graph not exist")
				return
			}
			// var jsonMap map[string]interface{}
			graphJsonListMap, convertErr := ConvertInterfaceToArrayOfMaps(graphJson)

			if convertErr != nil {
				log.Printf("EuroNewsGetAuthor convert @graph to map fail %v", convertErr)
				return
			}
			for _, currentJsonMap := range graphJsonListMap {
				if publishedAtTimestamp != 0 {
					break
				}
				datePublishedInterface, datePublishedInterfaceOk := currentJsonMap["datePublished"]
				if datePublishedInterfaceOk {
					currentPublishedAtStr := datePublishedInterface.(string)
					publishedAtTimestamp = ConvertStringTimeToTimestampForEuroNews(currentPublishedAtStr)
				}

			}

			// logger.Info("script content %s  author length %d",scriptContent, len(currentDWMetadata.Author))
		})
		if publishedAtTimestamp != 0 {
			break
		}
	}

	return publishedAtTimestamp

}

func (t *Template) EuroNewsGetAuthor(document *goquery.Document) string {
	author := ""
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

			var jsonMap map[string]interface{}
			unmarshalErr := json.Unmarshal([]byte(scriptContent), &jsonMap)
			if unmarshalErr != nil {
				log.Printf("EuroNewsGetAuthor unmarshal string to json fail")
				return
			}
			graphJson, ok := jsonMap["@graph"]
			if !ok {
				log.Printf("EuroNewsGetAuthor @graph not exist")
				return
			}
			// var jsonMap map[string]interface{}
			graphJsonListMap, convertErr := ConvertInterfaceToArrayOfMaps(graphJson)

			if convertErr != nil {
				log.Printf("EuroNewsGetAuthor convert @graph to map fail %v", convertErr)
				return
			}
			for _, currentJsonMap := range graphJsonListMap {
				if author != "" {
					break
				}
				authorInterface, authorInterfaceOk := currentJsonMap["author"]
				if authorInterfaceOk {
					authorInterfaceMap, authorInterfaceMapErr := ConvertInterfaceToMap(authorInterface)
					if authorInterfaceMapErr != nil {
						continue
					}
					authorValue, authorValueExist := authorInterfaceMap["name"]
					if authorValueExist {
						authorValueName, authorValueConvertStrOk := authorValue.(string)
						if authorValueConvertStrOk {
							author = authorValueName
						}
					}
				}

			}

			// logger.Info("script content %s  author length %d",scriptContent, len(currentDWMetadata.Author))
		})
		if author != "" {
			break
		}
	}

	return author

}

func ConvertInterfaceToMap(value interface{}) (map[string]interface{}, error) {
	if valueMap, ok := value.(map[string]interface{}); ok {
		return valueMap, nil
	} else {
		return nil, fmt.Errorf("value is not a map[string]interface{}")
	}
}
func (t *Template) EuroNewsScrapMetaData(document *goquery.Document) (string, string) {

	author := ""
	published_at := ""
	author, published_at = t.AuthorExtractFromScriptMetadata(document)
	if author == "" {
		author = t.EuroNewsGetAuthor(document)
		log.Printf("author********************** [%s]", author)
	}
	if author == "" {
		// cssSelectorList := make([]string,100);
		cssSelectorFirst := "article > div > div.o-article-newsy__main__body.u-article-content.u-article-grid > div.c-article-contributors > b:nth-child(1)"
		//cssSelectorList = append(cssSelectorList, cssSelectorFirst)
		// cssSelectorList = append(cssSelectorList, cssSelectorSecond)

		document.Find(cssSelectorFirst).Each(func(i int, s *goquery.Selection) {
			author = strings.TrimSpace(s.Text())
			// s.Find("b")
		})
		prev_selector_template := "#articlesSwiper > div.swiper-wrapper.swiper-wrapper--article > div.jsArticleFirst.u-transform-none.swiper-slide.swiper-slide-prev > article > div > div.o-article-newsy__main__body.u-article-content.u-article-grid > div.c-article-contributors > b:nth-child(%d)"
		active_selector_template := "#articlesSwiper > div.swiper-wrapper.swiper-wrapper--article > div.jsArticleFirst.u-transform-none.swiper-slide.swiper-slide-active > article > div > div.o-article-newsy__main__body.u-article-content.u-article-grid > div.c-article-contributors > b:nth-child(%d)"
		index := 1
		for {
			current_author := ""
			prev_template := fmt.Sprintf(prev_selector_template, index)
			activate_template := fmt.Sprintf(active_selector_template, index)
			document.Find(prev_template).Each(func(i int, s *goquery.Selection) {
				current_author = strings.TrimSpace(s.Text())
				// s.Find("b")
			})
			if current_author != "" {
				if index != 1 {
					author = author + " & "
				}
				author = author + current_author
				index = index + 1
				continue
			}
			document.Find(activate_template).Each(func(i int, s *goquery.Selection) {
				current_author = strings.TrimSpace(s.Text())
				// s.Find("b")
			})
			if current_author != "" {
				if index != 1 {
					author = author + " & "
				}
				author = author + current_author
				index = index + 1
				continue
			}
			if current_author == "" {
				break
			}

		}
	}

	link_template := "#articlesSwiper > div.swiper-wrapper.swiper-wrapper--article > div.jsArticleFirst.u-transform-none.swiper-slide.swiper-slide-active > article > div > div.o-article-newsy__main__body.u-article-content.u-article-grid > div.c-article-contributors > a"
	if author == "" {
		document.Find(link_template).Each(func(i int, s *goquery.Selection) {
			author = strings.TrimSpace(s.Text())
			// s.Find("b")
		})
	}

	return author, published_at
}

func (t *Template) EuroNewsScrapContent(document *goquery.Document) string {
	contents := ""
	document.Find("nav,h1.c-article-redesign-title,div.c-article-you-might-also-like,div.c-article-contributors,time.c-article-publication-date,div.c-ad__placeholder,a.c-article-partage-commentaire__links,div.c-article-caption,div.c-article-partage-commentaire-popup-overlay").Each(func(i int, s *goquery.Selection) {
		RemoveNodes(s)
	})
	document.Find("div.o-article-newsy__main").Each(func(i int, s *goquery.Selection) {
		var content string
		content, _ = goquery.OuterHtml(s)
		contents += content
	})
	return contents
}
