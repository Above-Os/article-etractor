package sanitizer

import (
	"bytes"
	"io"

	"golang.org/x/net/html"
)

// StripTags removes all HTML/XML tags from the input string.
func StripTags(input string) string {
	tokenizer := html.NewTokenizer(bytes.NewBufferString(input))
	var buffer bytes.Buffer

	for {
		if tokenizer.Next() == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				return buffer.String()
			}

			return ""
		}

		token := tokenizer.Token()
		switch token.Type {
		case html.TextToken:
			buffer.WriteString(token.Data)
		}
	}
}
