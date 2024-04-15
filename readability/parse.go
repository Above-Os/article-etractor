package readability

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

func FromReader(input io.Reader, pageURL string) (Article, error) {
	parser := NewReadability()
	return parser.ParseDocument(input, pageURL)
}

// extract published date from url if it's in the format of yyyy/mm/dd or yyyy-mm-dd
func extractPublishedDateFromUrl(url string) *time.Time {
	if url == "" {
		return nil
	}
	regex := regexp.MustCompile(`(\d{4})(\/|-)(\d{2})(\/|-)(\d{2})`)
	match := regex.FindStringSubmatch(url)
	if match != nil {
		year, _ := strconv.Atoi(match[1])
		month, _ := strconv.Atoi(match[3])
		day, _ := strconv.Atoi(match[5])
		date := time.Date(year, time.Month(month-1), day, 0, 0, 0, 0, time.UTC)
		return &date
	}
	return nil
}

func extractPublishedDateFromAuthor(author string) (string, *time.Time) {
	if author == "" {
		return "", nil
	}
	authorName := regexp.MustCompile(`^by\s+`).ReplaceAllString(author, "")
	regex := regexp.MustCompile(`(January|February|March|April|May|June|July|August|September|October|November|December)\s\d{1,2},\s\d{2,4}`)
	chineseDateRegex := regexp.MustCompile(`(\d{2,4})年(\d{1,2})月(\d{1,2})日`)
	// English date
	if regex.MatchString(author) {
		match := regex.FindStringSubmatch(author)
		t, _ := ParseTime(match[0])
		return regex.ReplaceAllString(authorName, ""), &t
	}
	// Chinese date
	if chineseDateRegex.MatchString(author) {
		match := chineseDateRegex.FindStringSubmatch(author)
		if match != nil {
			year, _ := strconv.Atoi(match[1])
			month, _ := strconv.Atoi(match[2])
			day, _ := strconv.Atoi(match[3])
			publishedAt := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
			return chineseDateRegex.ReplaceAllString(authorName, ""), &publishedAt
		}
	}
	return authorName, nil
}

// ParseDocument parses the specified document and find the main readable content.
func (r *Readability) ParseDocument(input io.Reader, pageURL string) (Article, error) {
	var err error

	// Reset parser data
	r.articleTitle = ""
	r.articleByline = ""
	r.articleSiteName = ""

	if r.documentURI, err = url.ParseRequestURI(pageURL); err != nil {
		return Article{}, fmt.Errorf("failed to parse URL: %v", err)
	}

	r.attempts = []parseAttempt{}
	r.flags.stripUnlikelys = true
	r.flags.useWeightClasses = true
	r.flags.cleanConditionally = true

	// Avoid parsing too large documents, as per configuration option.
	if r.MaxElemsToParse > 0 {
		numTags := len(getElementsByTagName(r.doc, "*"))

		if numTags > r.MaxElemsToParse {
			return Article{}, fmt.Errorf("too many elements: %d", numTags)
		}
	}
	if r.doc, err = html.Parse(input); err != nil {
		return Article{}, fmt.Errorf("failed to parse input: %v", err)
	}

	if r.documentURI.Host == "mp.weixin.qq.com" {
		r.wechatHandle(r.doc)
	}

	//html1 := innerHTML(r.doc)
	//InsertToFile("before Prepare1.html", html1)

	// Unwrap image from noscript
	r.unwrapNoscriptImages(r.doc)

	// Extract JSON-LD metadata before removing scripts
	var jsonLd map[string]string
	if !r.DisableJSONLD {
		jsonLd, _ = r.getJSONLD()
	}

	// Remove script tags from the document.
	r.removeScripts(r.doc)

	// Prepares the HTML document.
	r.prepDocument()

	// Fetch metadata.
	metadata := r.getArticleMetadata(jsonLd)
	r.articleTitle = metadata["title"]

	// Try to grab article content.
	finalHTMLContent := ""
	finalTextContent := ""
	readableNode := &html.Node{}
	articleContent := r.grabArticle()

	if articleContent != nil {
		r.postProcessContent(articleContent)

		// If we have not found an excerpt in the article's metadata, use the
		// article's first paragraph as the excerpt. This is used for displaying
		// a preview of the article's content.
		if metadata["excerpt"] == "" {
			paragraphs := getElementsByTagName(articleContent, "p")
			if len(paragraphs) > 0 {
				metadata["excerpt"] = strings.TrimSpace(textContent(paragraphs[0]))
			}
		}

		readableNode = firstElementChild(articleContent)

		finalHTMLContent = innerHTML(articleContent)
		finalTextContent = textContent(articleContent)
		finalTextContent = strings.TrimSpace(finalTextContent)
	}

	finalByline := metadata["byline"]
	if finalByline == "" {
		finalByline = r.articleByline
	}
	finalByline = strings.ToValidUTF8(finalByline, "")

	publishedData := extractPublishedDateFromUrl(pageURL)
	finalAuthor, publishedDateFromAuthor := extractPublishedDateFromAuthor(finalByline)

	metaPublishedData, dateParseErr := ParseTime(metadata["publishedTime"])

	if dateParseErr == nil {
		publishedData = &metaPublishedData
	}
	if publishedData == nil {
		publishedData = publishedDateFromAuthor
	}
	if publishedData == nil {
		publishedData = r.publishedDate
	}
	var finalPublishedData *time.Time
	if publishedData != nil {
		finalPublishedData = publishedData
	}

	// Excerpt is an supposed to be short and concise,
	// so it shouldn't have any new line
	excerpt := strings.TrimSpace(metadata["excerpt"])
	excerpt = strings.Join(strings.Fields(excerpt), " ")
	validTitle := strings.ToValidUTF8(r.articleTitle, "")

	validExcerpt := strings.ToValidUTF8(excerpt, "")
	return Article{
		Title:         validTitle,
		Byline:        finalAuthor,
		Node:          readableNode,
		Content:       finalHTMLContent,
		PublishedDate: finalPublishedData,
		TextContent:   finalTextContent,
		Length:        charCount(finalTextContent),
		Excerpt:       validExcerpt,
		SiteName:      metadata["siteName"],
		Image:         metadata["image"],
		Favicon:       metadata["favicon"],
		Language:      r.articleLang,
	}, nil
}
