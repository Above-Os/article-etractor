# Article-Extractor

It is a Go package that find the main readable content and the metadata from a HTML page. It works by removing clutter like buttons, ads, background images, script, etc.

This package is based from [Readability.js] by [Mozilla] and [omnivore]. 

For some websites, specific configuration templates are used to improve the accuracy of extractor.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)


## Installation

To install this package, just run `go get` :

```
go get github.com/Above-Os/article-extractor
```

## Usage

To get the readable content from an URL, you can use `processor.ArticleReadabilityExtractor`. It will fetch the web page from specified url, check if it's readable, then parses the response to find the readable content.  

| Input parameters                     | describe                                                   |
|--------------------------------------|------------------------------------------------------------|
| rawContent                           | raw content of the page                                    |
| entryUrl                             | url of the entry                                           |
| feedUrl                              | feed url， it can be "" if don’t have the value            |
| rules                                | custom parsing rules                                       |
| isrecommend                          | reserved parameters ,not used yet                          |


| Out parameters                       | describe                                                   |
|--------------------------------------|------------------------------------------------------------|
| content                              | content of the page                                        |
| pureContent                          | pure content                                               |
| publishedDate                        | published date,parsed by readability                       |
| image                                | cover image of the page                                    |
| title                                | title of the page                                          |
| author                               | author of the page,parsed by templates                     |
| byline                               | byline , parsed by readability                             |
| publishedAtTimeStamp                 | published timeStamp,parsed by templates                    |


To get the published date, publishedAtTimeStamp field can be used first, if the value is not empty. 
To get the author of article, author field can be used first, if the value is not empty. 
