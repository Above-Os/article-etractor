package readability

import (
	"encoding/json"
	"fmt"
	shtml "html"
	"math"
	nurl "net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var (
	rxPositive           = regexp.MustCompile(`(?i)article|body|content|entry|hentry|h-entry|main|page|pagination|post|text|blog|story|tweet(-\w+)?|instagram|image|container-banners|player|commentOnSelection`)
	rxNegative           = regexp.MustCompile(`(?i)-ad-|hidden|^hid$| hid$| hid |^hid |banner|combx|comment|com-|contact|foot|footer|footnote|gdpr|masthead|media|meta|outbrain|promo|related|scroll|share|shoutbox|sidebar|skyscraper|sponsor|shopping|tags|tool|widget|controls|video-controls`)
	rxUnlikelyCandidates = regexp.MustCompile(`(?i)-ad-|ai2html|banner|breadcrumbs|combx|comment|community|cover-wrap|disqus|extra|footer|gdpr|header|legends|menu|related|remark|replies|rss|shoutbox|sidebar|skyscraper|social|sponsor|supplemental|ad-break|agegate|pagination|pager|popup|yom-remote|copyright|keywords|outline|infinite-list|beta|recirculation|site-index|hide-for-print|post-end-share-cta|post-end-cta-full|post-footer|post-head|post-tag|li-date|main-navigation|programtic-ads|outstream_article|hfeed|comment-holder|back-to-top|show-up-next|onward-journey|topic-tracker|list-nav|block-ad-entity|adSpecs|gift-article-button|modal-title|in-story-masthead|share-tools|standard-dock|expanded-dock|margins-h|subscribe-dialog|icon|bumped|dvz-social-media-buttons|post-toc|mobile-menu|mobile-navbar|tl_article_header|mvp(-post)*-(add-story|soc(-mob)*-wrap)|w-condition-invisible|rich-text-block main w-richtext|rich-text-block_ataglance at-a-glance test w-richtext|PostsPage-commentsSection|hide-text`)
	/*
			articleNegativeLookBehindCandidates: /breadcrumbs|breadcrumb|utils|trilist|_header/i,
		    articleNegativeLookAheadCandidates: /outstream(.?)_|sub(.?)_|m_|omeda-promo-|in-article-advert|block-ad-.*|tl_/i,
			get okMaybeItsACandidate() {
		      return new RegExp(`and|(?<!${this.articleNegativeLookAheadCandidates.source})article(?!-(${this.articleNegativeLookBehindCandidates.source}))|body|column|content|^(?!main-navigation|main-header)main|shadow|post-header|hfeed site|blog-posts hfeed|container-banners|menu-opacity|header-with-anchor-widget|commentOnSelection`, 'i')
		    }*/
	rxOkMaybeItsACandidate = regexp.MustCompile(`(?i)and|article|body|column|content|main|shadow`)
	rxByline               = regexp.MustCompile(`(?i)byline|author|dateline|writtenby|p-author`)
	rxVideos               = regexp.MustCompile(`(?i)//(www\.)?((dailymotion|youtube|youtube-nocookie|player\.vimeo|v\.qq)\.com|(archive|upload\.wikimedia)\.org|player\.twitch\.tv|piped\.mha\.fi)`)

	rxHashURL            = regexp.MustCompile(`(?i)^#.+`)
	rxTokenize           = regexp.MustCompile(`(?i)\W+`)
	rxWhitespace         = regexp.MustCompile(`(?i)^\s*$`)
	rxHasContent         = regexp.MustCompile(`(?i)\S$`)
	rxNormalize          = regexp.MustCompile(`(?i)\s{2,}`)
	rxShareElements      = regexp.MustCompile(`(?i)(\b|_)(share|sharedaddy|post-tags)(\b|_)`)
	rxJsonLdArticleTypes = regexp.MustCompile(`(?i)^Article|AdvertiserContentArticle|NewsArticle|AnalysisNewsArticle|AskPublicNewsArticle|BackgroundNewsArticle|OpinionNewsArticle|ReportageNewsArticle|ReviewNewsArticle|Report|SatiricalArticle|ScholarlyArticle|MedicalScholarlyArticle|SocialMediaPosting|BlogPosting|LiveBlogPosting|DiscussionForumPosting|TechArticle|APIReference$`)
	rxSrcsetURL          = regexp.MustCompile(`(?i)(\S+)(\s+[\d.]+[xw])?(\s*(?:,|$))`)
	rxB64DataURL         = regexp.MustCompile(`(?i)^data:\s*([^\s;,]+)\s*;\s*base64\s*,`)

	rxLazyLoadingElements = regexp.MustCompile(`(?i)\S*loading\S*`)
	rxPublishedDate       = regexp.MustCompile(`(?i)published|modified|created|updated`)
	rxDates               = []*regexp.Regexp{
		regexp.MustCompile(`([0-9]{4}[-\/]?((0[13-9]|1[012])[-\/]?(0[1-9]|[12][0-9]|30)|(0[13578]|1[02])[-\/]?31|02[-\/]?(0[1-9]|1[0-9]|2[0-8]))|([0-9]{2}(([2468][048]|[02468][48])|[13579][26])|([13579][26]|[02468][048]|0[0-9]|1[0-6])00)[-\/]?02[-\/]?29)`),
		regexp.MustCompile(`(((0[13-9]|1[012])[-/]?(0[1-9]|[12][0-9]|30)|(0[13578]|1[02])[-/]?31|02[-/]?(0[1-9]|1[0-9]|2[0-8]))[-/]?[0-9]{4}|02[-/]?29[-/]?([0-9]{2}(([2468][048]|[02468][48])|[13579][26])|([13579][26]|[02468][048]|0[0-9]|1[0-6])00))`),
		regexp.MustCompile(`(((0[1-9]|[12][0-9]|30)[-/]?(0[13-9]|1[012])|31[-/]?(0[13578]|1[02])|(0[1-9]|1[0-9]|2[0-8])[-/]?02)[-/]?[0-9]{4}|29[-/]?02[-/]?([0-9]{2}(([2468][048]|[02468][48])|[13579][26])|([13579][26]|[02468][048]|0[0-9]|1[0-6])00))`),
	}
	rxLongDate    = regexp.MustCompile(`(?i)^(Jan(uary)?|Feb(ruary)?|Mar(ch)?|Apr(il)?|May|Jun(e)?|Jul(y)?|Aug(ust)?|Sep(tember)?|Oct(ober)?|Nov(ember)?|Dec(ember)?)\s\d{1,2}(?:st|nd|rd|th)?(,)?\s\d{2,4}$`)
	rxChineseDate = regexp.MustCompile(`(?i)^\d{2,4}年\d{1,2}月\d{1,2}日$`)

	rxPropertyPattern      = regexp.MustCompile(`(?i)\s*(dc|dcterm|og|twitter)\s*:\s*(author|creator|description|title|site_name|image\S*)\s*`)
	rxNamePattern          = regexp.MustCompile(`(?i)^\s*(?:(dc|dcterm|og|twitter|weibo:(article|webpage))\s*[\.:]\s*)?(author|creator|description|title|site_name|image)\s*$`)
	rxTitleSeparator       = regexp.MustCompile(`(?i) [\|\-\\/>»] `)
	rxTitleHierarchySep    = regexp.MustCompile(`(?i) [\\/>»] `)
	rxTitleRemoveFinalPart = regexp.MustCompile(`(?i)(.*)[\|\-\\/>»] .*`)
	rxTitleRemove1stPart   = regexp.MustCompile(`(?i)[^\|\-\\/>»]*[\|\-\\/>»](.*)`)
	rxTitleAnySeparator    = regexp.MustCompile(`(?i)[\|\-\\/>»]+`)
	rxDisplayNone          = regexp.MustCompile(`(?i)display\s*:\s*none`)
	rxSentencePeriod       = regexp.MustCompile(`(?i)\.( |$)`)
	rxFaviconSize          = regexp.MustCompile(`(?i)(\d+)x(\d+)`)
	rxLazyImageSrcset      = regexp.MustCompile(`(?i)\.(jpg|jpeg|png|webp)\s+\d`)
	rxLazyImageSrc         = regexp.MustCompile(`(?i)^\s*\S+\.(jpg|jpeg|png|webp)\S*\s*$`)
	rxImgExtensions        = regexp.MustCompile(`(?i)\.(jpg|jpeg|png|webp)`)
	rxCDATA                = regexp.MustCompile(`^\s*<!\[CDATA\[|\]\]>\s*$`)
	rxSchemaOrg            = regexp.MustCompile(`(?i)^https?\:\/\/schema\.org$`)
)

// Constants that used by readability.
var (
	unlikelyRoles                = sliceToMap("menu", "menubar", "complementary", "navigation", "alert", "alertdialog", "dialog")
	divToPElems                  = sliceToMap("blockquote", "dl", "div", "img", "ol", "p", "pre", "table", "ul", "select")
	alterToDivExceptions         = []string{"div", "article", "section", "p"}
	presentationalAttributes     = []string{"align", "background", "bgcolor", "border", "cellpadding", "cellspacing", "frame", "hspace", "rules", "style", "valign", "vspace"}
	deprecatedSizeAttributeElems = []string{"table", "th", "td", "hr", "pre"}
	phrasingElems                = []string{
		"abbr", "audio", "b", "bdo", "br", "button", "cite", "code", "data",
		"datalist", "dfn", "em", "embed", "i", "img", "input", "kbd", "label",
		"mark", "math", "meter", "noscript", "object", "output", "progress", "q",
		"ruby", "samp", "script", "select", "small", "span", "strong", "sub",
		"sup", "textarea", "time", "var", "wbr"}
	placeholerClasses = []string{"tweet-placeholder", "instagram-placeholder"}
)

// flags is flags that used by parser.
type flags struct {
	stripUnlikelys     bool
	useWeightClasses   bool
	cleanConditionally bool
}

// parseAttempt is container for the result of previous parse attempts.
type parseAttempt struct {
	articleContent *html.Node
	textLength     int
}

// Article is the final readable content.
type Article struct {
	Title         string
	Byline        string
	Node          *html.Node
	Content       string
	TextContent   string
	Length        int
	PublishedDate *time.Time
	Excerpt       string
	SiteName      string
	Image         string
	Favicon       string
	Language      string
}

// Parser is the parser that parses the page to get the readable content.
type Readability struct {
	// MaxElemsToParse is the max number of nodes supported by this
	// parser. Default: 0 (no limit)
	MaxElemsToParse int
	// NTopCandidates is the number of top candidates to consider when
	// analysing how tight the competition is among candidates.
	NTopCandidates int
	// CharThresholds is the default number of chars an article must
	// have in order to return a result
	CharThresholds int
	// ClassesToPreserve are the classes that readability sets itself.
	ClassesToPreserve []string
	// KeepClasses specify whether the classes should be stripped or not.
	KeepClasses bool
	// TagsToScore is element tags to score by default.
	TagsToScore []string
	// Debug determines if the log should be printed or not. Default: false.
	Debug bool
	// DisableJSONLD determines if metadata in JSON+LD will be extracted
	// or not. Default: false.
	DisableJSONLD bool
	// AllowedVideoRegex is a regular expression that matches video URLs that should be
	// allowed to be included in the article content. If undefined, it will use default filter.
	AllowedVideoRegex *regexp.Regexp

	doc           *html.Node
	documentURI   *nurl.URL
	articleTitle  string
	articleByline string
	publishedDate *time.Time
	//articleDir      string
	articleSiteName string
	articleLang     string
	attempts        []parseAttempt
	flags           flags
}

// NewParser returns new Parser which set up with default value.
func NewReadability() Readability {
	return Readability{
		MaxElemsToParse:   0,
		NTopCandidates:    5,
		CharThresholds:    500,
		ClassesToPreserve: []string{"page"},
		KeepClasses:       false,
		TagsToScore:       []string{"section", "h2", "h3", "h4", "h5", "h6", "p", "td", "pre"},
		Debug:             false,
	}
}
func (r *Readability) prepArticle(articleContent *html.Node) {
	r.cleanStyles(articleContent)

	// Check for data tables before we continue, to avoid removing
	// items in those tables, which will often be isolated even
	// though they're visually linked to other content-ful elements
	// (text, images, etc.).
	r.markDataTables(articleContent)

	r.fixLazyImages(articleContent)

	// Clean out junk from the article content
	r.cleanConditionally(articleContent, "form")
	r.cleanConditionally(articleContent, "fieldset")
	r.clean(articleContent, "object")
	r.clean(articleContent, "embed")
	r.clean(articleContent, "footer")
	r.clean(articleContent, "link")
	r.clean(articleContent, "aside")

	// Clean out elements have "share" in their id/class combinations
	// from final top candidates, which means we don't remove the top
	// candidates even they have "share".

	r.forEachNode(children(articleContent), func(topCandidate *html.Node, _ int) {
		r.cleanMatchedNodes(topCandidate, func(node *html.Node, nodeClassID string) bool {
			return rxShareElements.MatchString(nodeClassID) && len(textContent(node)) < r.CharThresholds
		})
	})

	r.clean(articleContent, "iframe")
	r.clean(articleContent, "input")
	r.clean(articleContent, "textarea")
	r.clean(articleContent, "select")
	r.clean(articleContent, "button")
	r.cleanHeaders(articleContent)
	r.cleanSmallImg(articleContent)

	// Do these last as the previous stuff may have removed junk
	// that will affect these
	r.cleanConditionally(articleContent, "table")
	r.cleanConditionally(articleContent, "ul")
	r.cleanConditionally(articleContent, "div")

	// Replace H1 with H2 as H1 should be only title that is displayed separately
	r.replaceNodeTags(r.getAllNodesWithTag(articleContent, "h1"), "h2")

	// Remove extra paragraphs
	r.removeNodes(getElementsByTagName(articleContent, "p"), func(p *html.Node) bool {
		imgCount := len(getElementsByTagName(p, "img")) + len(getElementsByTagName(p, "picture"))
		embedCount := len(getElementsByTagName(p, "embed"))
		objectCount := len(getElementsByTagName(p, "object"))

		// Nasty iframes have been removed, only remain embedded videos.
		iframeCount := len(getElementsByTagName(p, "iframe"))
		totalCount := imgCount + embedCount + objectCount + iframeCount
		return totalCount == 0 && r.getInnerText(p, false) == ""
	})

	r.forEachNode(getElementsByTagName(articleContent, "br"), func(br *html.Node, _ int) {
		next := r.nextElement(br.NextSibling)

		if next != nil && tagName(next) == "p" {
			br.Parent.RemoveChild(br)
		}
	})

	// Remove single-cell tables
	r.forEachNode(getElementsByTagName(articleContent, "table"), func(table *html.Node, _ int) {
		tbody := table
		if r.hasSingleTagInsideElement(table, "tbody") {
			tbody = firstElementChild(table)
		}

		if r.hasSingleTagInsideElement(tbody, "tr") {
			row := firstElementChild(tbody)
			if r.hasSingleTagInsideElement(row, "td") {
				cell := firstElementChild(row)

				newTag := "div"
				if r.everyNode(childNodes(cell), r.isPhrasingContent) {
					newTag = "p"
				}

				r.setNodeTag(cell, newTag)
				replaceChild(table.Parent, cell, table)
			}
		}
	})
}

// nextElement finds the next element, starting from the given node, and
// ignoring whitespace in between. If the given node is an element, the same
// node is returned.
func (r *Readability) nextElement(node *html.Node) *html.Node {
	next := node

	for next != nil &&
		next.Type != html.ElementNode &&
		rxWhitespace.MatchString(textContent(next)) {
		next = next.NextSibling
	}

	return next
}

// prepDocument prepares the HTML document for readability to scrape it.
// This includes things like stripping javascript, CSS, and handling
// terrible markup.
func (r *Readability) prepDocument() {
	doc := r.doc

	// ADDITIONAL, not exist in readability.js:
	// Remove all comments,
	r.removeComments(doc)

	// Remove all style tags in head
	r.removeNodes(getElementsByTagName(doc, "style"), nil)

	if nodes := getElementsByTagName(doc, "body"); len(nodes) > 0 && nodes[0] != nil {
		r.replaceBrs(nodes[0])
	}

	r.replaceNodeTags(getElementsByTagName(doc, "font"), "span")

	r.transformTagName(doc)

}

// postProcessContent runs any post-process modifications to article
// content as necessary.
func (r *Readability) postProcessContent(articleContent *html.Node) {
	// Readability cannot open relative uris so we convert them to absolute uris.
	r.fixRelativeURIs(articleContent)

	r.simplifyNestedElements(articleContent)

	// Remove classes.
	if !r.KeepClasses {
		r.cleanClasses(articleContent)
	}

	// Remove readability attributes.
	r.clearReadabilityAttr(articleContent)
}

func (r *Readability) grabArticle() *html.Node {
	for {
		doc := cloneNode(r.doc)

		var page *html.Node
		if nodes := getElementsByTagName(doc, "body"); len(nodes) > 0 {
			page = nodes[0]
		}

		// We can not grab an article if we do not have a page.
		if page == nil {
			return nil
		}

		// First, node prepping. Trash nodes that look cruddy (like ones with
		// the class name "comment", etc), and turn divs into P tags where they
		// have been used inappropriately (as in, where they contain no other
		// block level elements).
		var elementsToScore []*html.Node
		var node = documentElement(doc)
		shouldRemoveTitleHeader := true
		for node != nil {
			matchString := className(node) + "\x20" + id(node)

			if tagName(node) == "html" {
				r.articleLang = getAttribute(node, "lang")
			}
			if !r.isProbablyVisible(node) {
				node = r.removeAndGetNext(node)
				continue
			}
			// User is not able to see elements applied with both "aria-modal = true" and "role = dialog"
			if getAttribute(node, "aria-modal") == "true" && getAttribute(node, "role") == "dialog" {
				node = r.removeAndGetNext(node)
				continue
			}

			// Remove Node if it is a Byline.
			if r.checkByline(node, matchString) || r.checkPublishedDate(node, matchString) {
				node = r.removeAndGetNext(node)
				continue
			}

			if shouldRemoveTitleHeader && r.headerDuplicatesTitle(node) {
				shouldRemoveTitleHeader = false
				// Replacing title with the heading if the title includes heading but heading is smaller
				// Example article: http://jsomers.net/i-should-have-loved-biology
				// Or if there is the specific attribute that we can lean on.
				// For example "headline" in this article - https://nymag.com/intelligencer/2020/12/four-seasons-total-landscaping-the-full-est-possible-story.html
				headingText := trim(textContent(node))
				titleText := trim(r.articleTitle)
				if titleText != headingText && strings.Contains(titleText, headingText) || r.someNodeAttribute(node, func(name, value string) bool { return value == "headline" }) {
					r.articleTitle = headingText
				}
				node = r.removeAndGetNext(node)
				continue
			}

			// Remove unlikely candidates.
			nodeTagName := tagName(node)
			if r.flags.stripUnlikelys {
				if rxUnlikelyCandidates.MatchString(matchString) &&
					!rxOkMaybeItsACandidate.MatchString(matchString) &&
					!r.hasAncestorTag(node, "table", 3, nil) &&
					!r.hasAncestorTag(node, "code", 3, nil) &&
					// Example article - https://www.bbc.com/future/article/20170817-nasas-ambitious-plan-to-save-earth-from-a-supervolcano
					// Blockquote removed text example: "the conclusion that the supervolcano threat"
					!r.hasAncestorTag(node, "blockquote", 1, nil) &&
					nodeTagName != "body" &&
					nodeTagName != "a" {
					node = r.removeAndGetNext(node)
					continue
				}

				role := getAttribute(node, "role")
				if _, include := unlikelyRoles[role]; include {
					node = r.removeAndGetNext(node)
					continue
				}
			}

			// todo skip content modifications with embeds

			// Remove DIV, SECTION and HEADER nodes without any content.
			switch nodeTagName {
			case "div",
				"section",
				"header",
				"h1",
				"h2",
				"h3",
				"h4",
				"h5",
				"h6":
				if r.isElementWithoutContent(node) {
					node = r.removeAndGetNext(node)
					continue
				}
			}

			if indexOf(r.TagsToScore, nodeTagName) != -1 {
				elementsToScore = append(elementsToScore, node)
			}

			// Convert <div> without children block level elements into <p>.
			if nodeTagName == "div" {
				// Put phrasing content into paragraphs.
				var p *html.Node
				childNode := node.FirstChild

				for childNode != nil {
					nextSibling := childNode.NextSibling

					if r.isPhrasingContent(childNode) {
						if p != nil {
							appendChild(p, childNode)
						} else if !r.isWhitespace(childNode) {
							p = createElement("p")
							appendChild(p, cloneNode(childNode))
							replaceChild(node, p, childNode)
						}
					} else if p != nil {
						for p.LastChild != nil && r.isWhitespace(p.LastChild) {
							p.RemoveChild(p.LastChild)
						}
						p = nil
					}

					childNode = nextSibling
				}

				// Sites like http://mobile.slate.com encloses each paragraph
				// with a DIV element. DIVs with only a P element inside and no
				// text content can be safely converted into plain P elements to
				// avoid confusing the scoring algorithm with DIVs with are, in
				// practice, paragraphs.
				if r.hasSingleTagInsideElement(node, "p") && r.getLinkDensity(node) < 0.25 {
					newNode := children(node)[0]
					node, _ = replaceChild(node.Parent, newNode, node)
					elementsToScore = append(elementsToScore, node)
				} else if !r.hasChildBlockElement(node) {
					r.setNodeTag(node, "p")
					elementsToScore = append(elementsToScore, node)
				}
			}

			node = r.getNextNode(node, false)
		}

		//html1 := innerHTML(r.doc)
		//InsertToFile("before Prepare1.html", html1)

		// Loop through all paragraphs and assign a score to them based on how
		// much relevant content they have. Then add their score to their parent
		// node. A score is determined by things like number of commas, class
		// names, etc. Maybe eventually link density.
		var candidates []*html.Node
		r.forEachNode(elementsToScore, func(elementToScore *html.Node, _ int) {
			if elementToScore.Parent == nil || tagName(elementToScore.Parent) == "" {
				return
			}

			// If this paragraph is less than 25 characters, don't even count it.
			innerText := r.getInnerText(elementToScore, true)
			if charCount(innerText) < 25 {
				return
			}

			// Exclude nodes with no ancestor.
			ancestors := r.getNodeAncestors(elementToScore, 5)
			if len(ancestors) == 0 {
				return
			}

			// Add a point for the paragraph itself as a base.
			contentScore := 1

			// Add points for any commas within this paragraph.
			contentScore += strings.Count(innerText, ",")

			// For every 100 characters in this paragraph, add another point. Up to 3 points.
			contentScore += int(math.Min(math.Floor(float64(charCount(innerText))/100.0), 3.0))

			// Initialize and score ancestors.
			r.forEachNode(ancestors, func(ancestor *html.Node, level int) {
				if tagName(ancestor) == "" || ancestor.Parent == nil || ancestor.Parent.Type != html.ElementNode {
					return
				}

				if !r.hasContentScore(ancestor) {
					r.initializeNode(ancestor)
					candidates = append(candidates, ancestor)
				}

				// Node score divider:
				// - parent:             1 (no division)
				// - grandparent:        2
				// - great grandparent+: ancestor level * 3
				scoreDivider := 1
				switch level {
				case 0:
					scoreDivider = 1
				case 1:
					scoreDivider = 2
				default:
					scoreDivider = level * 3
				}

				ancestorScore := r.getContentScore(ancestor)
				ancestorScore += float64(contentScore) / float64(scoreDivider)

				r.setContentScore(ancestor, ancestorScore)
			})
		})

		// These lines are a bit different compared to Readability.js.
		//
		// In Readability.js, they fetch NTopCandidates utilising array method
		// like `splice` and `pop`. In Go, array method like that is not as
		// simple, especially since we are working with pointer. So, here we
		// simply sort top candidates, and limit it to max NTopCandidates.

		// Scale the final candidates score based on link density. Good
		// content should have a relatively small link density (5% or
		// less) and be mostly unaffected by this operation.
		for i := 0; i < len(candidates); i++ {
			candidate := candidates[i]
			candidateScore := r.getContentScore(candidate) * (1 - r.getLinkDensity(candidate))
			r.setContentScore(candidate, candidateScore)
		}

		// After we have calculated scores, sort through all of the possible
		// candidate nodes we found and find the one with the highest score.
		sort.Slice(candidates, func(i int, j int) bool {
			return r.getContentScore(candidates[i]) > r.getContentScore(candidates[j])
		})

		var topCandidates []*html.Node

		if len(candidates) > r.NTopCandidates {
			topCandidates = candidates[:r.NTopCandidates]
		} else {
			topCandidates = candidates
		}

		var topCandidate, parentOfTopCandidate *html.Node
		neededToCreateTopCandidate := false
		if len(topCandidates) > 0 {
			topCandidate = topCandidates[0]
		}

		// If we still have no top candidate, just use the body as a last
		// resort. We also have to copy the body node so it is something
		// we can modify.
		if topCandidate == nil || tagName(topCandidate) == "body" {
			// Move all of the page's children into topCandidate
			topCandidate = createElement("div")
			neededToCreateTopCandidate = true
			// Move everything (not just elements, also text nodes etc.)
			// into the container so we even include text directly in the body:
			kids := childNodes(page)
			for i := 0; i < len(kids); i++ {
				appendChild(topCandidate, kids[i])
			}

			appendChild(page, topCandidate)
			r.initializeNode(topCandidate)
		} else if topCandidate != nil {
			// Find a better top candidate node if it contains (at least three)
			// nodes which belong to `topCandidates` array and whose scores are
			// quite closed with current `topCandidate` node.
			topCandidateScore := r.getContentScore(topCandidate)
			var alternativeCandidateAncestors [][]*html.Node
			for i := 1; i < len(topCandidates); i++ {
				if r.getContentScore(topCandidates[i])/topCandidateScore >= 0.75 {
					topCandidateAncestors := r.getNodeAncestors(topCandidates[i], 0)
					alternativeCandidateAncestors = append(alternativeCandidateAncestors, topCandidateAncestors)
				}
			}

			minimumTopCandidates := 3
			if len(alternativeCandidateAncestors) >= minimumTopCandidates {
				parentOfTopCandidate = topCandidate.Parent
				for parentOfTopCandidate != nil && tagName(parentOfTopCandidate) != "body" {
					listContainingThisAncestor := 0
					for ancestorIndex := 0; ancestorIndex < len(alternativeCandidateAncestors) && listContainingThisAncestor < minimumTopCandidates; ancestorIndex++ {
						if includeNode(alternativeCandidateAncestors[ancestorIndex], parentOfTopCandidate) {
							listContainingThisAncestor++
						}
					}

					if listContainingThisAncestor >= minimumTopCandidates {
						topCandidate = parentOfTopCandidate
						break
					}

					parentOfTopCandidate = parentOfTopCandidate.Parent
				}
			}

			if !r.hasContentScore(topCandidate) {
				r.initializeNode(topCandidate)
			}

			// Because of our bonus system, parents of candidates might
			// have scores themselves. They get half of the node. There
			// won't be nodes with higher scores than our topCandidate,
			// but if we see the score going *up* in the first few steps *
			// up the tree, that's a decent sign that there might be more
			// content lurking in other places that we want to unify in.
			// The sibling stuff below does some of that - but only if
			// we've looked high enough up the DOM tree.
			parentOfTopCandidate = topCandidate.Parent
			lastScore := r.getContentScore(topCandidate)
			// The scores shouldn't get too lor.
			scoreThreshold := lastScore / 3.0
			for parentOfTopCandidate != nil && tagName(parentOfTopCandidate) != "body" {
				if !r.hasContentScore(parentOfTopCandidate) {
					parentOfTopCandidate = parentOfTopCandidate.Parent
					continue
				}

				parentScore := r.getContentScore(parentOfTopCandidate)
				if parentScore < scoreThreshold {
					break
				}

				if parentScore > lastScore {
					// Alright! We found a better parent to use.
					topCandidate = parentOfTopCandidate
					break
				}

				lastScore = parentScore
				parentOfTopCandidate = parentOfTopCandidate.Parent
			}

			// If the top candidate is the only child, use parent
			// instead. This will help sibling joining logic when
			// adjacent content is actually located in parent's
			// sibling node.
			parentOfTopCandidate = topCandidate.Parent
			for parentOfTopCandidate != nil && tagName(parentOfTopCandidate) != "body" && len(children(parentOfTopCandidate)) == 1 {
				topCandidate = parentOfTopCandidate
				parentOfTopCandidate = topCandidate.Parent
			}

			if !r.hasContentScore(topCandidate) {
				r.initializeNode(topCandidate)
			}
		}

		// Now that we have the top candidate, look through its siblings
		// for content that might also be related. Things like preambles,
		// content split by ads that we removed, etc.
		articleContent := createElement("div")
		siblingScoreThreshold := math.Max(10, r.getContentScore(topCandidate)*0.2)

		// Keep potential top candidate's parent node to try to get text direction of it later.
		topCandidateScore := r.getContentScore(topCandidate)
		topCandidateClassName := className(topCandidate)

		parentOfTopCandidate = topCandidate.Parent
		siblings := children(parentOfTopCandidate)
		for s := 0; s < len(siblings); s++ {
			sibling := siblings[s]
			appendNode := false

			if sibling == topCandidate {
				appendNode = true
			} else {
				contentBonus := float64(0)

				// Give a bonus if sibling nodes and top candidates have the example same classname
				if className(sibling) == topCandidateClassName && topCandidateClassName != "" {
					contentBonus += topCandidateScore * 0.2
				}

				if r.hasContentScore(sibling) && r.getContentScore(sibling)+contentBonus >= siblingScoreThreshold {
					appendNode = true
				} else if tagName(sibling) == "p" {
					linkDensity := r.getLinkDensity(sibling)
					nodeContent := r.getInnerText(sibling, true)
					nodeLength := len(nodeContent)

					if nodeLength > 80 && linkDensity < 0.25 {
						appendNode = true
					} else if nodeLength < 80 && nodeLength > 0 && linkDensity == 0 &&
						rxSentencePeriod.MatchString(nodeContent) {
						appendNode = true
					}
				}
			}

			if appendNode {
				// We have a node that is not a common block level element,
				// like a FORM or TD tag. Turn it into a DIV so it does not get
				// filtered out later by accident.
				if indexOf(alterToDivExceptions, tagName(sibling)) == -1 {
					r.setNodeTag(sibling, "div")
				}

				appendChild(articleContent, sibling)
			}
		}

		// So we have all of the content that we need. Now we clean
		// it up for presentation.
		r.prepArticle(articleContent)

		if neededToCreateTopCandidate {
			// We already created a fake DIV thing, and there would not have
			// been any siblings left for the previous loop, so there is no
			// point trying to create a new DIV and then move all the children
			// over. Just assign IDs and CSS class names here. No need to append
			// because that already happened anyway.
			//
			// By the way, this line is different with Readability.js.
			//
			// In Readability.js, when using `appendChild`, the node is still
			// referenced. Meanwhile here, our `appendChild` will clone the
			// node, put it in the new place, then delete the original.
			firstChild := firstElementChild(articleContent)
			if firstChild != nil && tagName(firstChild) == "div" {
				setAttribute(firstChild, "id", "readability-page-1")
				setAttribute(firstChild, "class", "page")
			}
		} else {
			div := createElement("div")
			setAttribute(div, "id", "readability-page-1")
			setAttribute(div, "class", "page")
			childs := childNodes(articleContent)

			for i := 0; i < len(childs); i++ {
				appendChild(div, childs[i])
			}
			appendChild(articleContent, div)
		}

		parseSuccessful := true

		// Now that we've gone through the full algorithm, check to see if we
		// got any meaningful content. If we did not, we may need to re-run
		// grabArticle with different flags set. This gives us a higher
		// likelihood of finding the content, and the sieve approach gives us a
		// higher likelihood of finding the -right- content.
		textLength := charCount(r.getInnerText(articleContent, true))
		if textLength < r.CharThresholds {
			parseSuccessful = false

			if r.flags.stripUnlikelys {
				r.flags.stripUnlikelys = false
				r.attempts = append(r.attempts, parseAttempt{
					articleContent: articleContent,
					textLength:     textLength,
				})
			} else if r.flags.useWeightClasses {
				r.flags.useWeightClasses = false
				r.attempts = append(r.attempts, parseAttempt{
					articleContent: articleContent,
					textLength:     textLength,
				})
			} else if r.flags.cleanConditionally {
				r.flags.cleanConditionally = false
				r.attempts = append(r.attempts, parseAttempt{
					articleContent: articleContent,
					textLength:     textLength,
				})
			} else {
				r.attempts = append(r.attempts, parseAttempt{
					articleContent: articleContent,
					textLength:     textLength,
				})

				// No luck after removing flags, just return the
				// longest text we found during the different loops *
				sort.Slice(r.attempts, func(i, j int) bool {
					return r.attempts[i].textLength > r.attempts[j].textLength
				})

				// But first check if we actually have something
				if r.attempts[0].textLength == 0 {
					return nil
				}

				articleContent = r.attempts[0].articleContent
				parseSuccessful = true
			}
		}

		if parseSuccessful {
			return articleContent
		}
	}
}

// getJSONLD try to extract metadata from JSON-LD object.
// For now, only Schema.org objects of type Article or its subtypes are supported.
func (r *Readability) getJSONLD() (map[string]string, error) {
	var metadata map[string]string

	scripts := querySelectorAll(r.doc, `script[type="application/ld+json"]`)
	r.forEachNode(scripts, func(jsonLdElement *html.Node, _ int) {
		if metadata != nil {
			return
		}

		// Strip CDATA markers if present
		content := rxCDATA.ReplaceAllString(textContent(jsonLdElement), "")

		// Decode JSON
		var parsed map[string]interface{}
		err := json.Unmarshal([]byte(content), &parsed)
		if err != nil {
			return
		}

		// Check context
		strContext, isString := parsed["@context"].(string)
		if !isString || !rxSchemaOrg.MatchString(strContext) {
			return
		}

		// If parsed doesn't have any @type, find it in its graph list
		if _, typeExist := parsed["@type"]; !typeExist {
			graphList, isArray := parsed["@graph"].([]interface{})
			if !isArray {
				return
			}

			for _, graph := range graphList {
				objGraph, isObj := graph.(map[string]interface{})
				if !isObj {
					continue
				}

				strType, isString := objGraph["@type"].(string)
				if isString && rxJsonLdArticleTypes.MatchString(strType) {
					parsed = objGraph
					break
				}
			}
		}

		// Once again, make sure parsed has valid @type
		strType, isString := parsed["@type"].(string)
		if !isString || !rxJsonLdArticleTypes.MatchString(strType) {
			return
		}

		// Initiate metadata
		metadata = make(map[string]string)

		// Title
		name, nameIsString := parsed["name"].(string)
		headline, headlineIsString := parsed["headline"].(string)

		if nameIsString && headlineIsString && name != headline {
			// We have both name and headline element in the JSON-LD. They should both be the same
			// but some websites like aktualne.cz put their own name into "name" and the article
			// title to "headline" which confuses Readability. So we try to check if either "name"
			// or "headline" closely matches the html title, and if so, use that one. If not, then
			// we use "name" by default.
			title := r.getArticleTitle()
			nameMatches := r.textSimilarity(name, title) > 0.75
			headlineMatches := r.textSimilarity(headline, title) > 0.75

			if headlineMatches && !nameMatches {
				metadata["title"] = headline
			} else {
				metadata["title"] = name
			}
		} else if name, isString := parsed["name"].(string); isString {
			metadata["title"] = strings.TrimSpace(name)
		} else if headline, isString := parsed["headline"].(string); isString {
			metadata["title"] = strings.TrimSpace(headline)
		}

		// Author
		switch val := parsed["author"].(type) {
		case map[string]interface{}:
			if name, isString := val["name"].(string); isString {
				metadata["byline"] = strings.TrimSpace(name)
			}

		case []interface{}:
			var authors []string
			for _, author := range val {
				objAuthor, isObj := author.(map[string]interface{})
				if !isObj {
					continue
				}

				if name, isString := objAuthor["name"].(string); isString {
					authors = append(authors, strings.TrimSpace(name))
				}
			}
			metadata["byline"] = strings.Join(authors, ", ")
		}

		// Description
		if description, isString := parsed["description"].(string); isString {
			metadata["excerpt"] = strings.TrimSpace(description)
		}

		// Publisher
		if objPublisher, isObj := parsed["publisher"].(map[string]interface{}); isObj {
			if name, isString := objPublisher["name"].(string); isString {
				metadata["siteName"] = strings.TrimSpace(name)
			}
		}
	})

	return metadata, nil
}

func (r *Readability) simplifyNestedElements(articleContent *html.Node) {
	node := articleContent

	for node != nil {
		nodeID := id(node)
		nodeTagName := tagName(node)

		if node.Parent != nil && (nodeTagName == "div" || nodeTagName == "section") &&
			!strings.HasPrefix(nodeID, "readability") {
			if r.isElementWithoutContent(node) {
				node = r.removeAndGetNext(node)
				continue
			}

			if r.hasSingleTagInsideElement(node, "div") || r.hasSingleTagInsideElement(node, "section") {
				child := children(node)[0]
				for _, attr := range node.Attr {
					setAttribute(child, attr.Key, attr.Val)
				}

				replaceChild(node.Parent, child, node)
				node = child
				continue
			}
		}

		node = r.getNextNode(node, false)
	}
}

// getArticleTitle attempts to get the article title.
func (r *Readability) getArticleTitle() string {
	doc := r.doc
	curTitle := ""
	origTitle := ""
	titleHadHierarchicalSeparators := false

	// If they had an element with tag "title" in their HTML
	if nodes := getElementsByTagName(doc, "title"); len(nodes) > 0 {
		origTitle = r.getInnerText(nodes[0], true)
		curTitle = origTitle
	}

	// If there's a separator in the title, first remove the final part
	if rxTitleSeparator.MatchString(curTitle) {
		titleHadHierarchicalSeparators = rxTitleHierarchySep.MatchString(curTitle)
		curTitle = rxTitleRemoveFinalPart.ReplaceAllString(origTitle, "$1")

		// If the resulting title is too short (3 words or fewer), remove
		// the first part instead:
		if wordCount(curTitle) < 3 {
			curTitle = rxTitleRemove1stPart.ReplaceAllString(origTitle, "$1")
		}
	} else if strings.Contains(curTitle, ": ") {
		// Check if we have an heading containing this exact string, so
		// we could assume it's the full title.
		headings := r.concatNodeLists(
			getElementsByTagName(doc, "h1"),
			getElementsByTagName(doc, "h2"),
		)

		trimmedTitle := strings.TrimSpace(curTitle)
		match := r.someNode(headings, func(heading *html.Node) bool {
			return strings.TrimSpace(textContent(heading)) == trimmedTitle
		})

		// If we don't, let's extract the title out of the original
		// title string.
		if !match {
			curTitle = origTitle[strings.LastIndex(origTitle, ":")+1:]

			// If the title is now too short, try the first colon instead:
			if wordCount(curTitle) < 3 {
				curTitle = origTitle[strings.Index(origTitle, ":")+1:]
				// But if we have too many words before the colon there's
				// something weird with the titles and the H tags so let's
				// just use the original title instead
			} else if wordCount(origTitle[:strings.Index(origTitle, ":")]) > 5 {
				curTitle = origTitle
			}
		}
	} else if charCount(curTitle) > 150 || charCount(curTitle) < 15 {
		if hOnes := getElementsByTagName(doc, "h1"); len(hOnes) == 1 {
			curTitle = r.getInnerText(hOnes[0], true)
		}
	}

	curTitle = strings.TrimSpace(curTitle)
	curTitle = rxNormalize.ReplaceAllString(curTitle, " ")
	// If we now have 4 words or fewer as our title, and either no
	// 'hierarchical' separators (\, /, > or ») were found in the original
	// title or we decreased the number of words by more than 1 word, use
	// the original title.
	curTitleWordCount := wordCount(curTitle)
	tmpOrigTitle := rxTitleAnySeparator.ReplaceAllString(origTitle, "")

	if curTitleWordCount <= 4 &&
		(!titleHadHierarchicalSeparators ||
			curTitleWordCount != wordCount(tmpOrigTitle)-1) {
		curTitle = origTitle
	}

	return curTitle
}

// replaceBrs replaces two or more successive <br> elements with a single <p>.
// Whitespace between <br> elements are ignored. For example:
//
//	<div>foo<br>bar<br> <br><br>abc</div>
//
// will become:
//
//	<div>foo<br>bar<p>abc</p></div>
func (r *Readability) replaceBrs(elem *html.Node) {
	r.forEachNode(r.getAllNodesWithTag(elem, "br"), func(br *html.Node, _ int) {
		next := br.NextSibling

		// Whether two or more <br> elements have been found and replaced with
		// a <p> block.
		replaced := false

		// If we find a <br> chain, remove the <br> nodes until we hit another
		// element or non-whitespace. This leaves behind the first <br> in the
		// chain (which will be replaced with a <p> later).
		for {
			next = r.nextElement(next)

			if next == nil || tagName(next) != "br" {
				break
			}

			replaced = true
			brSibling := next.NextSibling
			next.Parent.RemoveChild(next)
			next = brSibling
		}

		// If we removed a <br> chain, replace the remaining <br> with a <p>.
		// Add all sibling nodes as children of the <p> until we hit another
		// <br> chain.
		if replaced {
			p := createElement("p")
			replaceChild(br.Parent, p, br)

			next = p.NextSibling
			for next != nil {
				// If we have hit another <br><br>, we are done adding children
				// to this <p>.
				if tagName(next) == "br" {
					nextElem := r.nextElement(next.NextSibling)
					if nextElem != nil && tagName(nextElem) == "br" {
						break
					}
				}

				if !r.isPhrasingContent(next) {
					break
				}

				// Otherwise, make this node a child of the new <p>.
				sibling := next.NextSibling
				appendChild(p, next)
				next = sibling
			}

			for p.LastChild != nil && r.isWhitespace(p.LastChild) {
				p.RemoveChild(p.LastChild)
			}

			if tagName(p.Parent) == "p" {
				r.setNodeTag(p.Parent, "div")
			}
		}
	})
}

func (r *Readability) setNodeTag(node *html.Node, newTagName string) {
	if node.Type == html.ElementNode {
		node.Data = newTagName
	}

	// NOTES(cixtor): the original function in Readability.js is a bit longer
	// because it contains a fallback mechanism to set the node tag name just
	// in case JSDOMParser is not available, there is no need to implement this
	// here.
}

// compares second text to first one
// 1 = same text, 0 = completely different text
// works the way that it splits both texts into words and then finds words that are unique in second text
// the result is given by the lower length of unique parts
func (r *Readability) textSimilarity(textA, textB string) float64 {
	tokensA := rxTokenize.Split(strings.ToLower(textA), -1)
	tokensA = strFilter(tokensA, func(s string) bool { return s != "" })
	mapTokensA := sliceToMap(tokensA...)

	tokensB := rxTokenize.Split(strings.ToLower(textB), -1)
	tokensB = strFilter(tokensB, func(s string) bool { return s != "" })
	uniqueTokensB := strFilter(tokensB, func(s string) bool {
		_, existInA := mapTokensA[s]
		return !existInA
	})

	mergedB := strings.Join(tokensB, " ")
	mergedUniqueB := strings.Join(uniqueTokensB, " ")
	distanceB := float64(charCount(mergedUniqueB)) / float64(charCount(mergedB))

	return 1 - distanceB
}

/**
* Check if this node is an H1 or H2 element whose content is mostly
* the same as the article title.
*
* @param Element  the node to check.
* @return boolean indicating whether this is a title-like header.
 */
func (r *Readability) headerDuplicatesTitle(node *html.Node) bool {
	tag := strings.ToLower(tagName(node))
	if tag != "h1" && tag != "h2" {
		return false
	}
	heading := r.getInnerText(node, false)
	return r.textSimilarity(r.articleTitle, heading) > 0.75
}

func (r *Readability) checkPublishedDate(node *html.Node, matchString string) bool {
	// Skipping meta tags,and don't want to check for dates in the URL's
	tag := strings.ToLower(tagName(node))
	if tag == "meta" || tag == "a" {
		return false
	}

	// get the datetime from time element
	if tag == "time" {
		datetime := getAttribute(node, "datetime")
		t, err := ParseTime(datetime)
		if err == nil {
			r.publishedDate = &t
			return true
		}
	}

	content := strings.TrimSpace(textContent(node))
	if strings.Contains(content, "2023-12-21") {
		print('a')
	}
	var dateFound string
	var dateRegExpFound *regexp.Regexp
	for _, regexp := range rxDates {
		if regexp.MatchString(content) {
			dateRegExpFound = regexp
			break
		}
	}
	if dateRegExpFound != nil {
		dateFound = dateRegExpFound.FindString(content)
	} else if rxLongDate.MatchString(content) {
		dateFound = rxLongDate.FindString(content)
		dateFound = strings.ReplaceAll(dateFound, "st", "")
		dateFound = strings.ReplaceAll(dateFound, "nd", "")
		dateFound = strings.ReplaceAll(dateFound, "rd", "")
		dateFound = strings.ReplaceAll(dateFound, "th", "")
	} else if rxChineseDate.MatchString(content) {
		dateFound = rxChineseDate.FindString(content)
		dateFound = strings.ReplaceAll(dateFound, "年", "-")
		dateFound = strings.ReplaceAll(dateFound, "月", "-")
		dateFound = strings.ReplaceAll(dateFound, "日", "")
	}
	publishedDateParsed, dateContentParseErr := ParseTime(content)

	if (r.someNodeAttribute(node, func(name, value string) bool {
		if regexp.MustCompile(`href|uri|url`).MatchString(name) {
			return false
		}
		return rxPublishedDate.MatchString(value)
	}) || dateFound != "" || (strings.Contains(matchString, "date") && dateContentParseErr == nil)) && r.isValidPublishedDate(content) {
		var dateFoundParseErr error
		if dateContentParseErr != nil {
			publishedDateParsed, dateFoundParseErr = ParseTime(dateFound)
		}
		if dateFoundParseErr == nil && r.publishedDate == nil {
			r.publishedDate = &publishedDateParsed
		}
		return true
	}
	return false
}

// checkByline determines if a node is used as byline.
func (r *Readability) checkByline(node *html.Node, matchString string) bool {
	if r.articleByline != "" {
		return false
	}

	rel := getAttribute(node, "rel")
	itemprop := getAttribute(node, "itemprop")
	nodeText := textContent(node)
	if (rel == "author" || strings.Contains(itemprop, "author") || rxByline.MatchString(matchString)) && r.isValidByline(nodeText) {
		nodeText = strings.TrimSpace(nodeText)
		nodeText = strings.Join(strings.Fields(nodeText), "\x20")
		r.articleByline = nodeText
		return true
	}

	return false

}

func (r *Readability) isValidPublishedDate(publishedDate string) bool {
	publishedDate = strings.TrimSpace(publishedDate)
	return len(publishedDate) > 0 && len(publishedDate) < 50
}

// isValidByline checks whether the input string could be a byline.
func (r *Readability) isValidByline(byline string) bool {
	byline = strings.TrimSpace(byline)
	return len(byline) > 0 && len(byline) < 100
}

// cleanConditionally cleans an element of all tags of type "tag" if they look
// fishy. "Fishy" is an algorithm based on content length, classnames, link
// density, number of images & embeds, etc.
func (r *Readability) cleanConditionally(element *html.Node, tag string) {
	if !r.flags.cleanConditionally {
		return
	}

	// Prepare regex video filter
	rxVideoVilter := r.AllowedVideoRegex
	if rxVideoVilter == nil {
		rxVideoVilter = rxVideos
	}

	// Gather counts for other typical elements embedded within.
	// Traverse backwards so we can remove nodes at the same time
	// without effecting the traversal.
	// TODO: Consider taking into account original contentScore here.
	r.removeNodes(getElementsByTagName(element, tag), func(node *html.Node) bool {
		// First check if this node IS data table, in which case don't remove it.
		if tag == "table" && r.isReadabilityDataTable(node) {
			return false
		}

		//espn remove author https://www.espn.com/mlb/story/_/id/40016415/giants-blake-snell-placed-15-day-il-adductor-strain
		className := className(node)
		if className == "article-meta" {
			return true
		}
		// Do not clean placeholders
		if indexOf(placeholerClasses, className) != -1 {
			return false
		}

		isList := tag == "ul" || tag == "ol"

		if isList && isProbablyNavigation(node) {
			return true
		}

		if !isList {
			var listLength int
			listNodes := r.getAllNodesWithTag(node, "ul", "ol")
			r.forEachNode(listNodes, func(list *html.Node, _ int) {
				listLength += charCount(r.getInnerText(list, true))
			})

			nodeLength := charCount(r.getInnerText(node, true))
			isList = float64(listLength)/float64(nodeLength) > 0.9
		}

		// Next check if we're inside a data table, in which case don't remove it as well.
		if r.hasAncestorTag(node, "table", -1, r.isReadabilityDataTable) {
			return false
		}

		if r.hasAncestorTag(node, "code", 3, nil) {
			return false
		}

		// Avoiding lazyloaded images container removing
		// TODO: Rework this logic to work in a more robust and flexible way, this solution is fragile
		// Article example: https://nymag.com/intelligencer/2020/12/four-seasons-total-landscaping-the-full-est-possible-story.html
		childs := children(node)
		if len(childs) == 1 && strings.ToLower(tagName(childs[0])) == "picture" {
			return false
		}

		var contentScore int
		weight := r.getClassWeight(node)
		if weight+contentScore < 0 {
			return true
		}

		if r.getCharCount(node, ",") < 10 {
			// If there are not very many commas, and the number of
			// non-paragraph elements is more than paragraphs or other
			// ominous signs, remove the element.
			p := float64(len(getElementsByTagName(node, "p")))
			img := float64(len(getElementsByTagName(node, "img")))
			li := float64(len(getElementsByTagName(node, "li")) - 100)
			input := float64(len(getElementsByTagName(node, "input")))
			headingDensity := r.getTextDensity(node, "h1", "h2", "h3", "h4", "h5", "h6")

			embedCount := 0
			embeds := r.getAllNodesWithTag(node, "object", "embed", "iframe")

			for _, embed := range embeds {
				// If this embed has attribute that matches video regex,
				// don't delete it.
				for _, attr := range embed.Attr {
					if rxVideoVilter.MatchString(attr.Val) {
						return false
					}
				}

				// For embed with <object> tag, check inner HTML as well.
				if tagName(embed) == "object" && rxVideoVilter.MatchString(innerHTML(embed)) {
					return false
				}

				embedCount++
			}

			linkDensity := r.getLinkDensity(node)
			contentLength := charCount(r.getInnerText(node, true))
			haveToRemove := (img > 1 && p/img < 0.5 && !r.hasAncestorTag(node, "figure", 3, nil)) ||
				(!isList && li > p) ||
				(input > math.Floor(p/3)) ||
				(!isList && headingDensity < 0.9 && contentLength < 25 && (img == 0 || img > 2) && !r.hasAncestorTag(node, "figure", 3, nil)) ||
				(!isList && weight < 25 && linkDensity > 0.2) ||
				(weight >= 25 && linkDensity > 0.5) ||
				((embedCount == 1 && contentLength < 75) || embedCount > 1)

			// Allow simple lists of images to remain in pages
			if isList && haveToRemove {
				for _, child := range children(node) {
					// Don't filter in lists with li's that contain more than one child
					if len(children(child)) > 1 {
						return haveToRemove
					}
				}

				// Only allow the list to remain if every li contains an image
				liCount := len(getElementsByTagName(node, "li"))
				if int(img) == liCount {
					return false
				}
			}

			return haveToRemove
		}

		return false
	})
}

// cleanMatchedNodes cleans out elements whose ID and CSS class combinations
// match specific string.
func (r *Readability) cleanMatchedNodes(e *html.Node, filter func(*html.Node, string) bool) {
	endOfSearchMarkerNode := r.getNextNode(e, true)
	next := r.getNextNode(e, false)

	for next != nil && next != endOfSearchMarkerNode {
		if filter != nil && filter(next, className(next)+"\x20"+id(next)) {
			next = r.removeAndGetNext(next)
		} else {
			next = r.getNextNode(next, false)
		}
	}
}

// cleanHeaders cleans out spurious headers from an Element. Checks things like
// classnames and link density.
func (r *Readability) cleanHeaders(e *html.Node) {
	for headerIndex := 1; headerIndex < 3; headerIndex++ {
		headerTag := fmt.Sprintf("h%d", headerIndex)

		r.removeNodes(getElementsByTagName(e, headerTag), func(header *html.Node) bool {
			return r.getClassWeight(header) < 0
		})
	}
}

// fixRelativeURIs converts each <a> and <img> uri in the given element to an
// absolute URI, ignoring #ref URIs.
func (r *Readability) fixRelativeURIs(articleContent *html.Node) {
	links := r.getAllNodesWithTag(articleContent, "a")
	r.forEachNode(links, func(link *html.Node, _ int) {
		href := getAttribute(link, "href")
		if href == "" {
			return
		}

		// Remove links with javascript: URIs, since they won't
		// work after scripts have been removed from the page.
		if strings.HasPrefix(href, "javascript:") {
			linkChilds := childNodes(link)

			if len(linkChilds) == 1 && linkChilds[0].Type == html.TextNode {
				// If the link only contains simple text content,
				// it can be converted to a text node
				text := createTextNode(textContent(link))
				replaceChild(link.Parent, text, link)
			} else {
				// If the link has multiple children, they should
				// all be preserved
				container := createElement("span")
				for link.FirstChild != nil {
					appendChild(container, link.FirstChild)
				}
				replaceChild(link.Parent, container, link)
			}
		} else {
			newHref := toAbsoluteURI(href, r.documentURI)
			if newHref == "" {
				removeAttribute(link, "href")
			} else {
				setAttribute(link, "href", newHref)
			}
		}
	})

	medias := r.getAllNodesWithTag(articleContent, "img", "picture", "figure", "video", "audio", "source")
	r.forEachNode(medias, func(media *html.Node, _ int) {
		src := getAttribute(media, "src")
		poster := getAttribute(media, "poster")
		srcset := getAttribute(media, "srcset")

		if src != "" {
			newSrc := toAbsoluteURI(src, r.documentURI)
			setAttribute(media, "src", newSrc)
		}

		if poster != "" {
			newPoster := toAbsoluteURI(poster, r.documentURI)
			setAttribute(media, "poster", newPoster)
		}

		if srcset != "" {
			newSrcset := rxSrcsetURL.ReplaceAllStringFunc(srcset, func(s string) string {
				p := rxSrcsetURL.FindStringSubmatch(s)
				return toAbsoluteURI(p[1], r.documentURI) + p[2] + p[3]
			})

			setAttribute(media, "srcset", newSrcset)
		}
	})

}

// cleanClasses removes the class="" attribute from every element in the given
// subtree, except those that match CLASSES_TO_PRESERVE and classesToPreserve
// array from the options object.
func (r *Readability) cleanClasses(node *html.Node) {
	nodeClassName := className(node)
	preservedClassName := []string{}

	for _, class := range strings.Fields(nodeClassName) {
		if indexOf(r.ClassesToPreserve, class) != -1 {
			preservedClassName = append(preservedClassName, class)
		}
	}

	if len(preservedClassName) > 0 {
		setAttribute(node, "class", strings.Join(preservedClassName, "\x20"))
	} else {
		removeAttribute(node, "class")
	}

	for child := firstElementChild(node); child != nil; child = nextElementSibling(child) {
		r.cleanClasses(child)
	}
}

// clearReadabilityAttr removes Readability attribute created by the parser.
func (r *Readability) clearReadabilityAttr(node *html.Node) {
	removeAttribute(node, "data-readability-score")
	removeAttribute(node, "data-readability-table")

	for child := firstElementChild(node); child != nil; child = nextElementSibling(child) {
		r.clearReadabilityAttr(child)
	}
}

// isProbablyVisible determines if a node is visible.
func (r *Readability) isProbablyVisible(node *html.Node) bool {
	nodeStyle := getAttribute(node, "style")
	nodeAriaHidden := getAttribute(node, "aria-hidden")
	className := getAttribute(node, "class")

	return (nodeStyle == "" || !rxDisplayNone.MatchString(nodeStyle)) &&
		!hasAttribute(node, "hidden") &&
		(nodeAriaHidden == "" ||
			nodeAriaHidden != "true" ||
			strings.Contains(className, "fallback-image"))
}

// initializeNode initializes a node with the readability score. Also checks
// the className/id for special names to add to its score.
func (r *Readability) initializeNode(node *html.Node) {
	contentScore := float64(r.getClassWeight(node))

	switch tagName(node) {
	case "div":
		contentScore += 5
	case "pre", "td", "blockquote":
		contentScore += 3
	case "address", "ol", "ul", "dl", "dd", "dt", "li", "form":
		contentScore -= 3
	case "h1", "h2", "h3", "h4", "h5", "h6", "th":
		contentScore -= 5
	}

	r.setContentScore(node, contentScore)
}

// removeAndGetNext remove node and returns its next node.
func (r *Readability) removeAndGetNext(node *html.Node) *html.Node {
	nextNode := r.getNextNode(node, true)

	if node.Parent != nil {
		node.Parent.RemoveChild(node)
	}

	return nextNode
}

// getNextNode traverses the DOM from node to node, starting at the node passed
// in. Pass true for the second parameter to indicate this node itself (and its
// kids) are going away, and we want the next node over. Calling this in a loop
// will traverse the DOM depth-first.
//
// In Readability.js, ignoreSelfAndKids default to false.
func (r *Readability) getNextNode(node *html.Node, ignoreSelfAndKids bool) *html.Node {
	// First check for kids if those are not being ignored
	if firstChild := firstElementChild(node); !ignoreSelfAndKids && firstChild != nil {
		return firstChild
	}

	// Then for siblings...
	if sibling := nextElementSibling(node); sibling != nil {
		return sibling
	}

	// And finally, move up the parent chain *and* find a sibling
	// (because this is depth-first traversal, we will have already
	// seen the parent nodes themselves).
	for {
		node = node.Parent
		if node == nil || nextElementSibling(node) != nil {
			break
		}
	}

	if node != nil {
		return nextElementSibling(node)
	}

	return nil
}

// removeNodes iterates over a collection of HTML elements, calls the optional
// filter function on each node, and removes the node if function returns True.
// If function is not passed, removes all the nodes in the list.
func (r *Readability) removeNodes(list []*html.Node, filter func(*html.Node) bool) {
	var node *html.Node
	var parentNode *html.Node

	for i := len(list) - 1; i >= 0; i-- {
		node = list[i]
		parentNode = node.Parent

		if parentNode != nil && (filter == nil || filter(node)) {
			parentNode.RemoveChild(node)
		}
	}
}

// replaceNodeTags iterates over a list, and calls setNodeTag for each node.
func (r *Readability) replaceNodeTags(list []*html.Node, newTagName string) {
	for i := len(list) - 1; i >= 0; i-- {
		r.setNodeTag(list[i], newTagName)
	}
}

// forEachNode iterates over a list of HTML nodes, which doesn’t natively fully
// implement the Array interface. For convenience, the current object context
// is applied to the provided iterate function.
func (r *Readability) forEachNode(list []*html.Node, fn func(*html.Node, int)) {
	for idx, node := range list {
		fn(node, idx)
	}
}

// someNode iterates over a NodeList, return true if any of the
// provided iterate function calls returns true, false otherwise.
func (r *Readability) someNode(nodeList []*html.Node, fn func(*html.Node) bool) bool {
	for i := 0; i < len(nodeList); i++ {
		if fn(nodeList[i]) {
			return true
		}
	}

	return false
}

/**
 * Iterate over the attributes of the Element, return true if any of the provided iterate
 * function calls returns true, false otherwise.
 * @param {Element} node - Node to check for attributes
 * @param {function({name: string, value: string})} fn - The iterate function. Accepts object with name and value of the attribute
 */
func (r *Readability) someNodeAttribute(node *html.Node, fn func(name, value string) bool) bool {

	for i := 0; i < len(node.Attr); i++ {
		if fn(node.Attr[i].Key, node.Attr[i].Val) {
			return true
		}
	}

	return false
}

// concatNodeLists concats all nodelists passed as arguments.
func (r *Readability) concatNodeLists(nodeLists ...[]*html.Node) []*html.Node {
	var result []*html.Node

	for i := 0; i < len(nodeLists); i++ {
		result = append(result, nodeLists[i]...)
	}

	return result
}

func (r *Readability) getAllNodesWithTag(node *html.Node, tagNames ...string) []*html.Node {
	var list []*html.Node

	for _, tag := range tagNames {
		list = append(list, getElementsByTagName(node, tag)...)
	}

	return list
}

// everyNode iterates over a collection of nodes, returns true if all of the
// provided iterator function calls return true, otherwise returns false. For
// convenience, the current object context is applied to the provided iterator
// function.
func (r *Readability) everyNode(list []*html.Node, fn func(*html.Node) bool) bool {
	for _, node := range list {
		if !fn(node) {
			return false
		}
	}

	return true
}

// hasAncestorTag checks if a given node has one of its ancestor tag name
// matching the provided one.
//
// In Readability.js, default value for maxDepth is 3.
func (r *Readability) hasAncestorTag(node *html.Node, tag string, maxDepth int, filterFn func(*html.Node) bool) bool {
	depth := 0

	for node.Parent != nil {
		if maxDepth > 0 && depth > maxDepth {
			return false
		}

		if tagName(node.Parent) == tag && (filterFn == nil || filterFn(node.Parent)) {
			return true
		}

		node = node.Parent

		depth++
	}

	return false
}

// hasSingleTagInsideElement check if the node has only whitespace and a single
// element with given tag. Returns false if the DIV Node contains non-empty text
// nodes or if it contains no element with given tag or more than 1 element.
func (r *Readability) hasSingleTagInsideElement(element *html.Node, tag string) bool {
	// There should be exactly 1 element child with given tag
	if childs := children(element); len(childs) != 1 || tagName(childs[0]) != tag {
		return false
	}

	// And there should be no text nodes with real content
	return !r.someNode(childNodes(element), func(node *html.Node) bool {
		return node.Type == html.TextNode && rxHasContent.MatchString(textContent(node))
	})
}

func (r *Readability) getTextDensity(node *html.Node, tags ...string) float64 {
	textLength := charCount(r.getInnerText(node, true))
	if textLength == 0 {
		return 0
	}

	var childrenLength int
	children := r.getAllNodesWithTag(node, tags...)
	r.forEachNode(children, func(child *html.Node, _ int) {
		childrenLength += charCount(r.getInnerText(child, true))
	})

	return float64(childrenLength) / float64(textLength)
}

// getNodeAncestors gets the node's direct parent and grandparents.
//
// In Readability.js, maxDepth default to 0.
func (r *Readability) getNodeAncestors(node *html.Node, maxDepth int) []*html.Node {
	level := 0
	ancestors := []*html.Node{}

	for node.Parent != nil {
		level++
		ancestors = append(ancestors, node.Parent)

		if maxDepth > 0 && level == maxDepth {
			break
		}

		node = node.Parent
	}

	return ancestors
}

// isElementWithoutContent determines if node is empty. A node is considered
// empty is there is nothing inside or if the only things inside are HTML break
// tags <br> and HTML horizontal rule tags <hr>.
func (r *Readability) isElementWithoutContent(node *html.Node) bool {
	brs := getElementsByTagName(node, "br")
	hrs := getElementsByTagName(node, "hr")
	childs := children(node)

	return node.Type == html.ElementNode &&
		strings.TrimSpace(textContent(node)) == "" &&
		(len(childs) == 0 || len(childs) == len(brs)+len(hrs))
}

// hasChildBlockElement determines whether element has any children block level
// elements.

// setContentScore sets the readability score for a node.
func (r *Readability) setContentScore(node *html.Node, score float64) {
	setAttribute(node, "data-readability-score", fmt.Sprintf("%.4f", score))
}

// hasContentScore checks if node has readability score.
func (r *Readability) hasContentScore(node *html.Node) bool {
	return hasAttribute(node, "data-readability-score")
}

// getContentScore gets the readability score of a node.
func (r *Readability) getContentScore(node *html.Node) float64 {
	strScore := getAttribute(node, "data-readability-score")
	strScore = strings.TrimSpace(strScore)

	if strScore == "" {
		return 0
	}

	score, err := strconv.ParseFloat(strScore, 64)

	if err != nil {
		return 0
	}

	return score
}

func (r *Readability) removeComments(doc *html.Node) {
	var comments []*html.Node
	var finder func(*html.Node)

	finder = func(node *html.Node) {
		if node.Type == html.CommentNode {
			comments = append(comments, node)
		}

		for child := node.FirstChild; child != nil; child = child.NextSibling {
			finder(child)
		}
	}

	for child := doc.FirstChild; child != nil; child = child.NextSibling {
		finder(child)
	}

	r.removeNodes(comments, nil)
}

// isSingleImage checks if node is image, or if node contains exactly
// only one image whether as a direct child or as its descendants.
func (r *Readability) isSingleImage(node *html.Node) bool {
	if tagName(node) == "img" {
		return true
	}

	children := children(node)
	textContent := textContent(node)
	if len(children) != 1 || strings.TrimSpace(textContent) != "" {
		return false
	}

	return r.isSingleImage(children[0])
}

// getArticleFavicon attempts to get high quality favicon
// that used in article. It will only pick favicon in PNG
// format, so small favicon that uses ico file won't be picked.
// Using algorithm by philippe_b.
func (r *Readability) getArticleFavicon() string {
	favicon := ""
	faviconSize := -1
	linkElements := getElementsByTagName(r.doc, "link")

	r.forEachNode(linkElements, func(link *html.Node, _ int) {
		linkRel := strings.TrimSpace(getAttribute(link, "rel"))
		linkType := strings.TrimSpace(getAttribute(link, "type"))
		linkHref := strings.TrimSpace(getAttribute(link, "href"))
		linkSizes := strings.TrimSpace(getAttribute(link, "sizes"))

		if linkHref == "" || !strings.Contains(linkRel, "icon") {
			return
		}

		if linkType != "image/png" && !strings.Contains(linkHref, ".png") {
			return
		}

		size := 0
		for _, sizesLocation := range []string{linkSizes, linkHref} {
			sizeParts := rxFaviconSize.FindStringSubmatch(sizesLocation)
			if len(sizeParts) != 3 || sizeParts[1] != sizeParts[2] {
				continue
			}

			size, _ = strconv.Atoi(sizeParts[1])
			break
		}

		if size > faviconSize {
			faviconSize = size
			favicon = linkHref
		}
	})

	return toAbsoluteURI(favicon, r.documentURI)
}

// getArticleMetadata attempts to get excerpt and byline
// metadata for the article.
func (r *Readability) getArticleMetadata(jsonLd map[string]string) map[string]string {
	values := make(map[string]string)
	metaElements := getElementsByTagName(r.doc, "meta")

	// Find description tags.
	r.forEachNode(metaElements, func(element *html.Node, _ int) {
		elementName := getAttribute(element, "name")
		elementProperty := getAttribute(element, "property")
		content := getAttribute(element, "content")
		if content == "" {
			return
		}
		matches := []string{}
		name := ""

		if elementProperty != "" {
			matches = rxPropertyPattern.FindAllString(elementProperty, -1)
			for i := len(matches) - 1; i >= 0; i-- {
				// Convert to lowercase, and remove any whitespace
				// so we can match belops.
				name = strings.ToLower(matches[i])
				name = strings.Join(strings.Fields(name), "")
				// multiple authors
				values[name] = strings.TrimSpace(content)
			}
		}

		if len(matches) == 0 && elementName != "" && rxNamePattern.MatchString(elementName) {
			// Convert to lowercase, remove any whitespace, and convert
			// dots to colons so we can match belops.
			name = strings.ToLower(elementName)
			name = strings.Join(strings.Fields(name), "")
			name = strings.Replace(name, ".", ":", -1)
			values[name] = strings.TrimSpace(content)
		}
	})

	// get title
	metadataTitle := strOr(
		jsonLd["title"],
		values["dc:title"],
		values["dcterm:title"],
		values["og:title"],
		values["weibo:article:title"],
		values["weibo:webpage:title"],
		values["title"],
		values["twitter:title"])

	if metadataTitle == "" {
		metadataTitle = r.getArticleTitle()
	}

	// get author
	metadataByline := strOr(
		jsonLd["byline"],
		values["dc:creator"],
		values["dcterm:creator"],
		values["author"])

	// get description
	metadataExcerpt := strOr(
		jsonLd["excerpt"],
		values["dc:description"],
		values["dcterm:description"],
		values["og:description"],
		values["weibo:article:description"],
		values["weibo:webpage:description"],
		values["description"],
		values["twitter:description"])

	// get site name
	metadataSiteName := strOr(jsonLd["siteName"], values["og:site_name"])

	// get image thumbnail
	metadataImage := strOr(
		values["og:image"],
		values["image"],
		values["twitter:image"])

	// get favicon
	metadataFavicon := r.getArticleFavicon()

	// in many sites the meta value is escaped with HTML entities,
	// so here we need to unescape it
	metadataTitle = shtml.UnescapeString(metadataTitle)
	metadataByline = shtml.UnescapeString(metadataByline)
	metadataExcerpt = shtml.UnescapeString(metadataExcerpt)
	metadataSiteName = shtml.UnescapeString(metadataSiteName)

	return map[string]string{
		"title":    metadataTitle,
		"byline":   metadataByline,
		"excerpt":  metadataExcerpt,
		"siteName": metadataSiteName,
		"image":    metadataImage,
		"favicon":  metadataFavicon,
	}
}

// unwrapNoscriptImages finds all <noscript> that are located after <img> nodes,
// and which contain only one <img> element. Replace the first image with
// the image from inside the <noscript> tag, and remove the <noscript> tag.
// This improves the quality of the images we use on some sites (e.g. Medium).
func (r *Readability) unwrapNoscriptImages(doc *html.Node) {
	// Find img without source or attributes that might contains image, and
	// remove it. This is done to prevent a placeholder img is replaced by
	// img from noscript in next step.
	imgs := getElementsByTagName(doc, "img")
	r.forEachNode(imgs, func(img *html.Node, _ int) {
		for _, attr := range img.Attr {
			switch attr.Key {
			case "src", "data-src", "srcset", "data-srcset":
				return
			}

			if rxImgExtensions.MatchString(attr.Val) {
				return
			}
		}
		if tagName(img.Parent) != "picture" {
			img.Parent.RemoveChild(img)
		}
	})

	// Next find noscript and try to extract its image
	noscripts := getElementsByTagName(doc, "noscript")
	r.forEachNode(noscripts, func(noscript *html.Node, _ int) {
		// Parse content of noscript and make sure it only contains image
		noscriptContent := textContent(noscript)
		tmpDoc, err := html.Parse(strings.NewReader(noscriptContent))
		if err != nil {
			return
		}

		tmpBodyElems := getElementsByTagName(tmpDoc, "body")
		if len(tmpBodyElems) == 0 {
			return
		}

		tmpBody := tmpBodyElems[0]
		if !r.isSingleImage(tmpBodyElems[0]) {
			return
		}

		// If noscript has previous sibling and it only contains image,
		// replace it with noscript content. However we also keep old
		// attributes that might contains image.
		prevElement := previousElementSibling(noscript)
		if prevElement != nil && r.isSingleImage(prevElement) {
			prevImg := prevElement
			if tagName(prevImg) != "img" {
				prevImg = getElementsByTagName(prevElement, "img")[0]
			}

			newImg := getElementsByTagName(tmpBody, "img")[0]
			for _, attr := range prevImg.Attr {
				if attr.Val == "" {
					continue
				}

				if attr.Key == "src" || attr.Key == "srcset" || rxImgExtensions.MatchString(attr.Val) {
					if getAttribute(newImg, attr.Key) == attr.Val {
						continue
					}

					attrName := attr.Key
					if hasAttribute(newImg, attrName) {
						attrName = "data-old-" + attrName
					}

					setAttribute(newImg, attrName, attr.Val)
				}
			}

			replaceChild(noscript.Parent, firstElementChild(tmpBody), prevElement)
		}
	})
}

// removeScripts removes script tags from the document.
func (r *Readability) removeScripts(doc *html.Node) {
	r.removeNodes(getElementsByTagName(doc, "script"), nil)
	r.removeNodes(getElementsByTagName(doc, "noscript"), nil)
}

func (r *Readability) transformTagName(doc *html.Node) {
	r.forEachNode(getElementsByTagName(doc, "progressive-image"), func(imageNode *html.Node, _ int) {
		srcAttr := getAttribute(imageNode, "src")
		altAttr := getAttribute(imageNode, "alt")
		replaceWithHtml(imageNode, `<img src="`+srcAttr+`" alt="`+altAttr+`"/>`)

	})
}

// hasChildBlockElement determines whether element has any children block level
// elements.
func (r *Readability) hasChildBlockElement(element *html.Node) bool {
	return r.someNode(childNodes(element), func(node *html.Node) bool {
		_, exist := divToPElems[tagName(node)]
		return exist || r.hasChildBlockElement(node)
	})

}

// isPhrasingContent determines if a node qualifies as phrasing content.
func (r *Readability) isPhrasingContent(node *html.Node) bool {
	if node.Type == html.TextNode {
		return true
	}

	tag := tagName(node)
	if indexOf(phrasingElems, tag) != -1 {
		return true
	}

	return ((tag == "a" || tag == "del" || tag == "ins") &&
		r.everyNode(childNodes(node), r.isPhrasingContent))
}

// isWhitespace determines if a node only used as whitespace.
func (r *Readability) isWhitespace(node *html.Node) bool {
	return (node.Type == html.TextNode && strings.TrimSpace(textContent(node)) == "") ||
		(node.Type == html.ElementNode && tagName(node) == "br")
}

// getInnerText gets the inner text of a node.
// This also strips * out any excess whitespace to be found.
// In Readability.js, normalizeSpaces default to true.
func (r *Readability) getInnerText(node *html.Node, normalizeSpaces bool) string {
	textContent := strings.TrimSpace(textContent(node))
	if normalizeSpaces {
		textContent = rxNormalize.ReplaceAllString(textContent, " ")
	}
	return textContent
}

// getCharCount returns the number of times a string s
// appears in the node.
func (r *Readability) getCharCount(node *html.Node, s string) int {
	innerText := r.getInnerText(node, true)
	return strings.Count(innerText, s)
}

// cleanStyles removes the style attribute on every node and under.
func (r *Readability) cleanStyles(node *html.Node) {
	nodeTagName := tagName(node)
	if node == nil || nodeTagName == "svg" {
		return
	}

	if nodeTagName == "img" {
		//nodeStyle := getAttribute(node, "style")

	}

	// Remove `style` and deprecated presentational attributes
	for i := 0; i < len(presentationalAttributes); i++ {
		removeAttribute(node, presentationalAttributes[i])
	}

	if indexOf(deprecatedSizeAttributeElems, nodeTagName) != -1 {
		removeAttribute(node, "width")
		removeAttribute(node, "height")
	}

	for child := firstElementChild(node); child != nil; child = nextElementSibling(child) {
		r.cleanStyles(child)
	}
}

// getLinkDensity gets the density of links as a percentage of the
// content. This is the amount of text that is inside a link divided
// by the total text in the node.
func (r *Readability) getLinkDensity(element *html.Node) float64 {
	textLength := charCount(r.getInnerText(element, true))
	if textLength == 0 {
		return 0
	}

	var linkLength float64
	r.forEachNode(getElementsByTagName(element, "a"), func(linkNode *html.Node, _ int) {
		href := getAttribute(linkNode, "href")
		href = strings.TrimSpace(href)

		coefficient := 1.0
		if href != "" && rxHashURL.MatchString(href) {
			coefficient = 0.3
		}

		nodeLength := charCount(r.getInnerText(linkNode, true))
		linkLength += float64(nodeLength) * coefficient
	})

	return linkLength / float64(textLength)
}

// getClassWeight gets an elements class/id weight. Uses regular
// expressions to tell if this element looks good or bad.
func (r *Readability) getClassWeight(node *html.Node) int {
	if !r.flags.useWeightClasses {
		return 0
	}

	weight := 0

	// Look for a special classname
	if nodeClassName := className(node); nodeClassName != "" {
		if rxNegative.MatchString(nodeClassName) {
			weight -= 25
		}

		if rxPositive.MatchString(nodeClassName) {
			weight += 25
		}
	}

	// Look for a special ID
	if nodeID := id(node); nodeID != "" {
		if rxNegative.MatchString(nodeID) {
			weight -= 25
		}

		if rxPositive.MatchString(nodeID) {
			weight += 25
		}
	}

	return weight
}

// clean cleans a node of all elements of type "tag".
// (Unless it's a youtube/vimeo video. People love movies.)
func (r *Readability) clean(node *html.Node, tag string) {
	isEmbed := indexOf([]string{"object", "embed", "iframe"}, tag) != -1
	rxVideoVilter := r.AllowedVideoRegex
	if rxVideoVilter == nil {
		rxVideoVilter = rxVideos
	}

	r.removeNodes(getElementsByTagName(node, tag), func(element *html.Node) bool {
		// Allow youtube and vimeo videos through as people usually want to see those.
		if isEmbed {
			// First, check the elements attributes to see if any of them contain
			// youtube or vimeo
			for _, attr := range element.Attr {
				if rxVideoVilter.MatchString(attr.Val) {
					return false
				}
			}

			// For embed with <object> tag, check inner HTML as well.
			if tagName(element) == "object" && rxVideoVilter.MatchString(innerHTML(element)) {
				return false
			}
		}
		return true
	})
}

// getRowAndColumnCount returns how many rows and columns this table has.
func (r *Readability) getRowAndColumnCount(table *html.Node) (int, int) {
	rows := 0
	columns := 0
	trs := getElementsByTagName(table, "tr")
	for i := 0; i < len(trs); i++ {
		strRowSpan := getAttribute(trs[i], "rowspan")
		rowSpan, _ := strconv.Atoi(strRowSpan)
		if rowSpan == 0 {
			rowSpan = 1
		}
		rows += rowSpan

		// Now look for column-related info
		columnsInThisRow := 0
		cells := getElementsByTagName(trs[i], "td")
		for j := 0; j < len(cells); j++ {
			strColSpan := getAttribute(cells[j], "colspan")
			colSpan, _ := strconv.Atoi(strColSpan)
			if colSpan == 0 {
				colSpan = 1
			}
			columnsInThisRow += colSpan
		}

		if columnsInThisRow > columns {
			columns = columnsInThisRow
		}
	}

	return rows, columns
}

// isReadabilityDataTable determines if a Node is a data table.
func (r *Readability) isReadabilityDataTable(node *html.Node) bool {
	return hasAttribute(node, "data-readability-table")
}

// setReadabilityDataTable marks whether a Node is data table or not.
func (r *Readability) setReadabilityDataTable(node *html.Node, isDataTable bool) {
	if isDataTable {
		setAttribute(node, "data-readability-table", "true")
		return
	}

	removeAttribute(node, "data-readability-table")
}

// markDataTables looks for 'data' (as opposed to 'layout') tables
// and mark it, which similar as used in Firefox:
// https://searchfox.org/mozilla-central/rev/f82d5c549f046cb64ce5602bfd894b7ae807c8f8/accessible/generic/TableAccessible.cpp#19
func (r *Readability) markDataTables(root *html.Node) {
	tables := getElementsByTagName(root, "table")
	for i := 0; i < len(tables); i++ {
		table := tables[i]

		role := getAttribute(table, "role")
		if role == "presentation" {
			r.setReadabilityDataTable(table, false)
			continue
		}

		datatable := getAttribute(table, "datatable")
		if datatable == "0" {
			r.setReadabilityDataTable(table, false)
			continue
		}

		if hasAttribute(table, "summary") {
			r.setReadabilityDataTable(table, true)
			continue
		}

		if captions := getElementsByTagName(table, "caption"); len(captions) > 0 {
			if caption := captions[0]; caption != nil && len(childNodes(caption)) > 0 {
				r.setReadabilityDataTable(table, true)
				continue
			}
		}

		// If the table has a descendant with any of these tags, consider a data table:
		hasDataTableDescendantTags := false
		for _, descendantTag := range []string{"col", "colgroup", "tfoot", "thead", "th"} {
			descendants := getElementsByTagName(table, descendantTag)
			if len(descendants) > 0 && descendants[0] != nil {
				hasDataTableDescendantTags = true
				break
			}
		}

		if hasDataTableDescendantTags {
			r.setReadabilityDataTable(table, true)
			continue
		}

		// Nested tables indicates a layout table:
		if len(getElementsByTagName(table, "table")) > 0 {
			r.setReadabilityDataTable(table, false)
			continue
		}

		rows, columns := r.getRowAndColumnCount(table)
		if rows >= 10 || columns > 4 {
			r.setReadabilityDataTable(table, true)
			continue
		}

		// Now just go by size entirely:
		if rows*columns > 10 {
			r.setReadabilityDataTable(table, true)
		}
	}
}

func (r *Readability) cleanSmallImg(element *html.Node) {
	r.removeNodes(getElementsByTagName(element, "img"), func(node *html.Node) bool {
		width := parseInt(getAttribute(node, "width"), 0)
		return width != 0 && width < 100
	})
}

func (r *Readability) getDynamicImageSrc(elem *html.Node) string {
	candidateAttrs := []string{
		"_src",
		"data-src",
		"data-original",
		"data-orig",
		"data-url",
		"data-orig-file",
		"data-large-file",
		"data-medium-file",
		"data-2000src",
		"data-1000src",
		"data-800src",
		"data-655src",
		"data-500src",
		"data-380src",
		"nitro-lazy-src",
		"data-ezsrc",
		"data-lazy",
		"data-lazy-src",
	}
	for _, attr := range candidateAttrs {
		val := getAttribute(elem, attr)
		if val != "" {
			return val
		}

	}
	return ""

}

// fixLazyImages convert images and figures that have properties like data-src into
// images that can be loaded without JS.
func (r *Readability) fixLazyImages(root *html.Node) {
	imageNodes := r.getAllNodesWithTag(root, "img", "picture", "figure", "svg")
	r.forEachNode(imageNodes, func(elem *html.Node, _ int) {
		src := getAttribute(elem, "src")
		srcset := getAttribute(elem, "srcset")
		nodeTag := tagName(elem)
		nodeClass := className(elem)
		dynamicsSrc := r.getDynamicImageSrc(elem)
		if dynamicsSrc != "" {
			setAttribute(elem, "src", dynamicsSrc)
		}
		/*
			lazySrc := getAttribute(elem, "data-lazy-src")
			if lazySrc != "" {
				setAttribute(elem, "src", lazySrc)
			}*/

		// In some sites (e.g. Kotaku), they put 1px square image as base64 data uri in
		// the src attribute. So, here we check if the data uri is too short, just might
		// as well remove it.
		if src != "" && rxB64DataURL.MatchString(src) {
			// Make sure it's not SVG, because SVG can have a meaningful image in
			// under 133 bytes.
			parts := rxB64DataURL.FindStringSubmatch(src)
			if parts[1] == "image/svg+xml" {
				return
			}

			// Make sure this element has other attributes which contains image.
			// If it doesn't, then this src is important and shouldn't be removed.
			srcCouldBeRemoved := false
			for _, attr := range elem.Attr {
				if attr.Key == "src" {
					continue
				}

				if rxImgExtensions.MatchString(attr.Val) && IsValidURL(attr.Val) {
					srcCouldBeRemoved = true
					break
				}
			}

			// Here we assume if image is less than 100 bytes (or 133B
			// after encoded to base64) it will be too small, therefore
			// it might be placeholder image.
			if srcCouldBeRemoved {
				b64starts := strings.Index(src, "base64") + 7
				b64length := len(src) - b64starts
				if b64length < 133 {
					src = ""
					removeAttribute(elem, "src")
				}
			}
		}

		if (src != "" || srcset != "") && !strings.Contains(strings.ToLower(nodeClass), "lazy") {
			// Removing image that is redundant loading placeholder
			// Example article: https://www.instyle.com/celebrity/gigi-hadid/gigi-hadid-bangs-2020 (image className: "loadingPlaceholder")
			if nodeClass != "" && rxLazyLoadingElements.MatchString(nodeClass) {
				elem.Parent.RemoveChild(elem)
			}
			return
		}

		for i := 0; i < len(elem.Attr); i++ {
			attr := elem.Attr[i]
			if attr.Key == "src" || attr.Key == "srcset" || attr.Key == "alt" {
				continue
			}

			copyTo := ""
			if rxLazyImageSrcset.MatchString(attr.Val) {
				copyTo = "srcset"
			} else if rxLazyImageSrc.MatchString(attr.Val) {
				copyTo = "src"
			}

			if copyTo == "" || !IsValidURL(attr.Val) {
				continue
			}

			if nodeTag == "img" || nodeTag == "picture" {
				// if this is an img or picture, set the attribute directly
				setAttribute(elem, copyTo, attr.Val)
			} else if nodeTag == "figure" && len(r.getAllNodesWithTag(elem, "img", "picture")) == 0 {
				// if the item is a <figure> that does not contain an image or picture,
				// create one and place it inside the figure see the nytimes-3
				// testcase for an example
				img := createElement("img")
				setAttribute(img, copyTo, attr.Val)
				appendChild(elem, img)
			}
		}
	})
}

func parseWidthFromStyle(style string) string {
	/*widthPattern := regexp.MustCompile(`(?i)\s*width:\s*(.*?);`)
	matches := widthPattern.FindStringSubmatch(style)
	if len(matches) > 1 {
		return matches[1]
	}*/
	items := strings.Split(style, ";")
	for _, item := range items {
		propArr := strings.Split(item, ":")
		if strings.TrimSpace(propArr[0]) == "width" {
			return strings.TrimSpace(propArr[1])
		}
	}

	return ""
}

func (r *Readability) wechatHandle(doc *html.Node) {
	sections := getElementsByTagName(doc, "section")
	r.forEachNode(sections, func(section *html.Node, _ int) {
		style := getAttribute(section, "style")
		if style != "" {
			styleWidth := parseWidthFromStyle(style)
			if strings.HasSuffix(styleWidth, "px") {
				styleWidth = styleWidth[0 : len(styleWidth)-2]
				w := parseInt(styleWidth, 100)
				if w < 60 {
					section.Parent.RemoveChild(section)
				}
			}
		}
	})

	imgs := getElementsByTagName(doc, "img")
	r.forEachNode(imgs, func(img *html.Node, _ int) {
		oriWidth := getAttribute(img, "width")
		if strings.HasSuffix(oriWidth, "%") {
			removeAttribute(img, "width")
			oriWidth = ""
		}
		style := getAttribute(img, "style")
		if style != "" {
			styleWidth := parseWidthFromStyle(style)
			if strings.HasSuffix(styleWidth, "px") {
				styleWidth = styleWidth[0 : len(styleWidth)-2]
				picWidth := parseInt(styleWidth, 0)
				if oriWidth == "" && picWidth < 100 {
					img.Parent.RemoveChild(img)
				}
			} else {
				width1 := 100
				if strings.HasSuffix(styleWidth, "%") {
					width1 = parseInt(styleWidth[0:len(styleWidth)-1], 100)
				}
				parentStyle := getAttribute(img.Parent, "style")
				parentStyleWidth := parseWidthFromStyle(parentStyle)
				if strings.HasSuffix(parentStyleWidth, "px") {
					parentWidth := parentStyleWidth[0 : len(parentStyleWidth)-2]
					picWidth := parseInt(parentWidth, 0) * width1 / 100
					if oriWidth == "" && picWidth < 100 {
						img.Parent.RemoveChild(img)
					}
				} else {
					if strings.HasSuffix(parentStyleWidth, "%") {
						width2 := parseInt(parentStyleWidth[0:len(parentStyleWidth)-1], 100)
						width1 = width2 * width1 / 100
					}
					if oriWidth == "" && width1 < 8 {
						img.Parent.RemoveChild(img)
					}
				}
			}

		}

	})
}
