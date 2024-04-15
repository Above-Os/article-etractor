package processor

import (
	"net/url"
	"strings"
)

// domain => CSS selectors
var contentPredefinedRules = map[string]string{
	"blog.cloudflare.com":  "div.post-content",
	"cbc.ca":               ".story-content",
	"darkreading.com":      "#article-main:not(header)",
	"developpez.com":       "div[itemprop=articleBody]",
	"dilbert.com":          "span.comic-title-name, img.img-comic",
	"explosm.net":          "div#comic",
	"financialsamurai.com": "article",
	"francetvinfo.fr":      ".text",
	"github.com":           "article.entry-content",
	"heise.de":             "header .article-content__lead, header .article-image, div.article-layout__content.article-content",
	"igen.fr":              "section.corps",
	"ikiwiki.iki.fi":       ".page.group",
	"ilpost.it":            ".entry-content",
	"ing.dk":               "section.body",
	"lapresse.ca":          ".amorce, .entry",
	"lemonde.fr":           "article",
	"lepoint.fr":           ".art-text",
	"lesjoiesducode.fr":    ".blog-post-content img",
	"lesnumeriques.com":    ".text",
	"linux.com":            "div.content, div[property]",
	"mac4ever.com":         "div[itemprop=articleBody]",
	"monwindows.com":       ".blog-post-body",
	"npr.org":              "#storytext",
	"oneindia.com":         ".io-article-body",
	"opensource.com":       "div[property]",
	"openingsource.org":    "article.suxing-popup-gallery",
	"osnews.com":           "div.newscontent1",
	"phoronix.com":         "div.content",
	"pseudo-sciences.org":  "#art_main",
	"quantamagazine.org":   ".outer--content, figure, script",
	"raywenderlich.com":    "article",
	"royalroad.com":        ".author-note-portlet,.chapter-content",
	"slate.fr":             ".field-items",
	"smbc-comics.com":      "div#cc-comicbody, div#aftercomic",
	"swordscomic.com":      "img#comic-image, div#info-frame.tab-content-area",
	"theoatmeal.com":       "div#comic",
	"theregister.com":      "#top-col-story h2, #body",
	"theverge.com":         "h2.inline:nth-child(2),h2.duet--article--dangerously-set-cms-markup,figure.w-full,div.duet--article--article-body-component",
	"turnoff.us":           "article.post-content",
	"universfreebox.com":   "#corps_corps",
	"version2.dk":          "section.body",
	"wdwnt.com":            "div.entry-content",
	"wired.com":            "div.grid-layout__content",
	"zeit.de":              ".summary, .article-body",
	"zdnet.com":            "div.storyBody",
	"pbfcomics":            "div#comic",
	"yahoo.com":            "div.caas-body",
	"kyivindependent.com":  "div.c-content",
	"news.mit.edu":         "div.news-article--content--body--inner",
}

var contentPostExtractorTemplateRules = map[string]string{
	"weixin.qq.com": "WechatPostExtractor",
}

var contentTemplatePredefinedRules = map[string]string{
	"abcnews.go.com":                "AbcNewsScrapContent",
	"cnbc.com":                      "CNBCScrapContent",
	"bbc.co.uk":                     "BBCScrapContent",
	"bbc.com":                       "BBCScrapContent",
	"telegraph.co.uk":               "TelegraphScrapContent",
	"thestar.com":                   "TheStarScrapContent",
	"medium.com":                    "MediumScrapContent",
	"medium.datadriveninvestor.com": "MediumScrapContent",
	"cbsnews.com":                   "CbsNewsScrapContent",
	"news.sky.com":                  "SkyNewsScrapContent",
	"www.aljazeera.com":             "AljazeeraScrapContent",
	"themoscowtimes.com":            "MoscowTimesScrapContent",
	"themessenger.com":              "MessengerScrapContent",
	"euronews.com":                  "EuroNewsScrapContent",
	"huffpost.com":                  "HuffPostScrapContent",
	"dw.com":                        "DWScrapContent",
	"foxnews.com":                   "FoxNewsScrapContent",
	"pravda.com":                    "PravdaScrapContent",
	"time.com":                      "TimeScrapContent",
	"theguardian.com":               "TheguardianScrapContent",
	"reuters.com":                   "ReutersScrapContent",
	"abc.net.au":                    "AbcNetAUScrapContent",
	"yahoo.com":                     "YahoocrapContent",
	"nbcnews.com":                   "NbcNewsScrapContent",
}

var metadataTemplatePredefinedRules = map[string]string{
	"abcnews.go.com":     "AbcNewsScrapMetaData",
	"apnews.com":         "ApnNewsScrapMetaData",
	"www.aljazeera.com":  "AljazeeraScrapMetaData",
	"news.sky.com":       "SkyNewsScrapMetaData",
	"yahoo.com":          "YahooNewsScrapMetaData",
	"abc.net.au":         "AbcNetAUScrapMetaData",
	"cbsnews.com":        "CbsNewsScrapMetaData",
	"cnbc.com":           "CnbcScrapMetaData",
	"dw.com":             "DWScrapMetaData",
	"euronews.com":       "EuroNewsScrapMetaData",
	"foxnews.com":        "FoxNewsScrapMetaData",
	"huffpost.com":       "HuffPostScrapMetaData",
	"nbcnews.com":        "NbcNewsScrapMetaData",
	"ndtv.com":           "NdtvNewsScrapMetaData",
	"pravda.com":         "PravdaScrapMetaData",
	"themoscowtimes.com": "ThemoscowtimesScrapMetaData",
	"themessenger.com":   "ThemessengerScrapMetaData",
	"theguardian.com":    "TheguardianScrapMetaData",
	"bbc.com":            "BBCScrapMetaData",
	".bbc.co.":           "BBCScrapMetaData",
	"time.com":           "TimesScrapMetaData",
}

var publishedAtTimeStampTemplatePredefinedRules = map[string]string {
	"abcnews.go.com":     "CommonGetPublishedAtTimestampSingleJson",
	"apnews.com":         "CommonGetPublishedAtTimestampSingleJson",
	"www.aljazeera.com":  "CommonGetPublishedAtTimestampSingleJson",
	"news.sky.com":       "CommonGetPublishedAtTimestampSingleJson",
	"yahoo.com":          "CommonGetPublishedAtTimestampSingleJson",
	"abc.net.au":         "CommonGetPublishedAtTimestampSingleJson",
	"cbsnews.com":        "CbsnewsWorldGetPublishedAtTimestampSingleJson",
	"cnbc.com":           "CnbcWorldGetPublishedAtTimestampSingleJson",
	"dw.com":             "CommonGetPublishedAtTimestampSingleJson",
	"euronews.com":       "EuroNewsGetPublishedAtTimeStamp",
	"foxnews.com":        "CommonGetPublishedAtTimestampSingleJson",
	"huffpost.com":       "CommonGetPublishedAtTimestampSingleJson",
	"nbcnews.com":        "CommonGetPublishedAtTimestampSingleJson",
	"ndtv.com":           "NdtvGetPublishedAtTimestamp",
	"pravda.com":         "CommonGetPublishedAtTimestampSingleJson",
	"themoscowtimes.com": "CommonGetPublishedAtTimestampSingleJson",
	"themessenger.com":   "TheMessengerGetPublishedAtTimestampSingleJson",
	"theguardian.com":    "CommonGetPublishedAtTimestampMultipleJson",
	"bbc.com":            "CommonGetPublishedAtTimestampSingleJson",
	".bbc.co.":           "CommonGetPublishedAtTimestampSingleJson",
	"time.com":           "CommonGetPublishedAtTimestampMultipleJson",

}


func getPredefinedPublishedAtTimestampTemplateRules(websiteURL string) (string, string) {
	urlDomain := domain(websiteURL)

	for domain, rules := range publishedAtTimeStampTemplatePredefinedRules {
		if strings.Contains(urlDomain, domain) {
			return domain, rules
		}
	}
	return "", ""
}


func getContentPostExtractorTemplateRules(websiteURL string) string {
	urlDomain := domain(websiteURL)
	for url, rules := range contentPostExtractorTemplateRules {
		if strings.Contains(urlDomain, url) {
			return rules
		}
	}
	return ""
}

func getPredefinedScraperRules(websiteURL string) (string, string) {
	urlDomain := domain(websiteURL)

	for domain, rules := range contentPredefinedRules {
		if strings.Contains(urlDomain, domain) {
			return domain, rules
		}
	}
	return "", ""
}

func getPredefinedContentTemplateRules(websiteURL string) (string, string) {
	urlDomain := domain(websiteURL)

	for domain, rules := range contentTemplatePredefinedRules {
		if strings.Contains(urlDomain, domain) {
			return domain, rules
		}
	}
	return "", ""
}

func getPredefinedMetaDataTemplateRules(websiteURL string) (string, string) {
	urlDomain := domain(websiteURL)

	for domain, rules := range metadataTemplatePredefinedRules {
		if strings.Contains(urlDomain, domain) {
			return domain, rules
		}
	}
	return "", ""
}

func domain(websiteURL string) string {
	parsedURL, err := url.Parse(websiteURL)
	if err != nil {
		return websiteURL
	}

	return parsedURL.Host
}
