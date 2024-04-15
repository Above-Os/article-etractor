package templates

import (

	"github.com/PuerkitoBio/goquery"
)

func (t *Template) ApnNewsScrapMetaData(document *goquery.Document) (string, string) {
	author := ""
	published_at := ""
	author,published_at= t.AuthorExtractFromScriptMetadata(document)
	if author == "" {
		author = "AP News"
    }

	return author, published_at
}

