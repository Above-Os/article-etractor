package postExtractor

import (
	"log"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

var (
	maxDepth               = 3
	articleEndExactMatch   = []string{"完", "推荐阅读", "推荐阅读", "参考资料", "更多阅读", "也许你还想看:", "铁粉推荐", "原创不易，感谢有你！", "咨询合作，请联系微信", "<本篇完>", "往期优质文章：", "延伸阅读", "往期精选", "点击阅读原文", "加入交流群&内容转载&合作相关", "商业财经 国际视角"}
	articleEndSimilarMatch = []string{"本文图片仅用于图片介绍", "识别下图二维码", "推荐阅读", "文中观点仅为作者观点", "点击下方阅读原文", "华丽的分割线", "推荐你关注", "您的在看，是我创作的动力", "对本文亦有贡献", "添加微信烦劳", "添加时请注明", "本文仅代表作者本人观点", "本文仅代表作者观点", "点击阅读原文", "欢迎按指引星标tech星球"}
	rxarticleEndElements   = regexp.MustCompile(`(?i)\S*关注\S*微信公众号\S*`)

	//1海豚投研
	removeTextMatch = []string{"一般披露提示"}
	//1 晚点 4\5钛媒体
	removeImgSrcMatch = []string{"https://mmbiz.qpic.cn/mmbiz_png/VWpZENjIo5v8Z0EK9CBxR8L5nABo6f5qdtW2IzuRp67mhlLxo0WzY4QJzYcXKicm1tqmLcyworGzUmEYYhoia0Ig/640?wx_fmt=png&from=appmsg",
		"https://mmbiz.qpic.cn/mmbiz_png/VWpZENjIo5uFjORrmxUNN3iaHGFZqC8oCBsicHA3cujRQYNHLShjicag526TXFmymiczG9ncbj5atv07CmQDan6NKw/640?wx_fmt=png&from=appmsg",
		"https://mmbiz.qpic.cn/mmbiz_png/VWpZENjIo5vibTpO6sKfyAG1KYEcfNMBhgAgq5hpAWM70Xoq8A4kGmlXRFWH2DIQPKiaSXCKWoMk9aYcDWKicb1Yg/640?wx_fmt=png&from=appmsg",
		"https://mmbiz.qpic.cn/mmbiz_gif/OaFsUa11r0CALGfmcHTGJLo2T5jniaTo8Kg7s8ZwlUX6ostHVrZnl2d96wA7AkaIoTdFlfJqurGVV2ve19A2eog/640?wx_fmt=gif&from=appmsg",
		"https://mmbiz.qpic.cn/mmbiz_gif/OaFsUa11r0CALGfmcHTGJLo2T5jniaTo8Kg7s8ZwlUX6ostHVrZnl2d96wA7AkaIoTdFlfJqurGVV2ve19A2eog/640?wx_fmt=gif&from=appmsg"}

	//1极客公园 2 36客 3财经杂志 4财经无忌 5一点财经 6投资界 7 创业帮
	endImgSrcMatch = []string{"https://mmbiz.qpic.cn/mmbiz_png/8cu01Kavc5ZEzbKBVM2xq7iavNBh3BS8pRQ8mFL6zhL2hxID7Co8G29MQHWKr4aFxH4zQv64PkgTEZlxOKVH0nw/640?wx_fmt=png",
		"https://mmbiz.qpic.cn/mmbiz_png/QicyPhNHD5vZ3Txm1k8cuNfHWPCVgS2F5kKTm9MalhVfvXiaTBy8ia1rH39Jicc03tCXCzzMz4Hico0xLRVLmibQiaEYA/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1",
		"https://mmbiz.qpic.cn/mmbiz_png/ia1nxOhDj7ASqxWCrRia780M4Jxo6iaJ531IhXrbwxnPOJVcgpkc5yndJqanqIWJuDlINrZHCQMdVezacao0HNsfg/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1",
		"https://mmbiz.qpic.cn/sz_mmbiz_gif/QPAn9ic1pIlic1zbMQeia4icF0b3oLC2wcg6RbIPiaxeQ4flibxKSGjFFHHCbeG5guq1If6TUNrDKkWajgaUCwfoickGA/640?wx_fmt=gif&wxfrom=5&wx_lazy=1&wx_co=1",
		"https://mmbiz.qpic.cn/mmbiz_png/nyYLSss2B3unolCbzURQaNTKh6WpI0rRPOf9OKHWqRtq99g8QKT2Km84pdf6wMdLaIn8k9qf9XQl0XM1qm4cvQ/640?wx_fmt=png",
		"https://mmbiz.qpic.cn/mmbiz_gif/PIKqwNU2vsfvSic60uj6Hich3TFjdZooDhhbib0tTE4yiaUEGWAtgBBLyChG3cNqhjedeqJ1IbxE8zMic63HI5PvRFg/640?wx_fmt=gif&from=appmsg",
		"https://mmbiz.qpic.cn/mmbiz_png/03KNO9Fib2w4vwgRlwq2uFGLSGoGz6vrxSDBwlvePGicDdY5wNIbIJpHx2MDNRHyibiaCsmSCbUusybibPytLRrTtDQ/640?wx_fmt=png&wxfrom=5&wx_lazy=1&wx_co=1"}
)

type PostExtractorTemplate struct{}

func urlMatch(array []string, key string) int {
	for idx, val := range array {

		if val == strings.TrimSpace(key) {
			log.Printf("wechat  url match....%s,%s", val, key)
			return idx
		}
	}
	return -1
}
func textMatch(array []string, key, matchType string) int {
	for idx, val := range array {

		if matchType == "exact" && val == strings.ToLower(strings.TrimSpace(key)) {
			log.Printf("wechat match....%s,%s", val, key)
			return idx
		}
		if matchType == "similar" && strings.Contains(strings.ToLower(key), val) {
			log.Printf("wechat match....%s,%s", val, key)
			return idx
		}
	}
	return -1
}

func countENNonEmpty(arr []string) int {
	var count int
	for _, value := range arr {
		if len(value) > 0 {
			count++
		}
	}
	return count
}
func checkArticleEnd(depth int, nodes []*goquery.Selection) (bool, int) {
	isArticleEnd := false
	nodeIndex := 0
	articleValidIndex := len(nodes)
	for ; nodeIndex < len(nodes); nodeIndex++ {
		curNode := nodes[nodeIndex]
		htmlNode := curNode.Get(0)
		if depth < maxDepth && len(curNode.Children().Nodes) > 1 {
			subNodes := make([]*goquery.Selection, 0)
			curNode.Children().Each(func(i int, child *goquery.Selection) {
				subNodes = append(subNodes, child)
			})
			isArticleEnd, _ = checkArticleEnd(depth+1, subNodes)
			if isArticleEnd {
				break
			}
		} else {
			if htmlNode.Data == "img" {
				src, _ := curNode.Attr("src")
				if urlMatch(endImgSrcMatch, src) != -1 {
					log.Printf("wechat check articleend  by img  ....%s", src)
					isArticleEnd = true
					break
				}
			} else {
				nodeText := curNode.Text()
				if len(nodeText) > 0 {
					nodeText = strings.ReplaceAll(nodeText, "\u00a0", " ")
					strArr := strings.Split(nodeText, " ")
					if countENNonEmpty(strArr) < 5 {
						for _, str := range strArr {
							lowerStr := strings.ToLower(str)
							if lowerStr == "end" || lowerStr == "完" || (len(lowerStr) < 6 && strings.Contains(lowerStr, "end")) {
								isArticleEnd = true
								break
							}
						}
					}
					if isArticleEnd {
						break
					}

					text := strings.TrimSpace(curNode.Text())

					if textMatch(articleEndExactMatch, text, "exact") != -1 || textMatch(articleEndSimilarMatch, text, "similar") != -1 || rxarticleEndElements.MatchString(text) {
						//log.Printf("wechat check article end ....%s", text)
						isArticleEnd = true
						break
					}
				}

			}
		}
	}
	if isArticleEnd {
		articleValidIndex = nodeIndex
		for ; nodeIndex < len(nodes); nodeIndex++ {
			nodes[nodeIndex].Remove()
		}

	}
	return isArticleEnd, articleValidIndex
}

func checkIsEmptyNode(node *goquery.Selection) bool {
	item := node.Get(0)
	if item.Data == "a" {
		text := node.Text()
		if len(text) == 0 && item.FirstChild == nil {
			return true
		}
	}
	if item.Data == "p" || item.Data == "section" || item.Data == "strong" {
		text := node.Text()
		if len(text) == 0 && len(node.Children().Nodes) < 1 {
			return true
		}
	}

	return false
}
func checkIsImgNode(node *goquery.Selection) bool {
	//https://mp.weixin.qq.com/s?__biz=MjM5MjAyNDUyMA==&mid=2650979519&idx=2&sn=813b43bc467f6b469c59e6a07e0101bc&chksm=bcd031c9106ef1bc5663187aed2955ccea4c0812f1920254af896b5e2e97450beeae65df3a47&scene=0&xtrack=1#rd
	fatherBro := node.Parent().Next()
	if fatherBro != nil {
		text := fatherBro.Text()
		if len(text) > 0 {
			return false
		}
	}
	htmlNode := node.Get(0)
	if htmlNode.Data == "img" || htmlNode.Data == "br" {
		return true
	}
	if checkIsEmptyNode(node) {
		return true
	}
	if len(node.Children().Nodes) < 1 {
		return false
	}
	if len(node.Children().Nodes) == 1 {
		if node.Children().Nodes[0].Data == "img" {
			return true
		}
	}
	if len(node.Children().Nodes) > 1 {
		isImg := true
		node.Children().Each(func(i int, child *goquery.Selection) {
			htmlNode := child.Get(0)
			grandChildIsImg := false
			if len(child.Children().Nodes) == 1 {
				if child.Children().Nodes[0].Data == "img" {
					grandChildIsImg = true
				}
			}
			if htmlNode.Data != "img" && htmlNode.Data != "br" && !grandChildIsImg {
				isImg = false
			}
		})
		return isImg
	}
	return false

}

func removeEmptyTagsFirst(document *goquery.Document) {
	document.Find("span,h1,h2,hr").Each(func(i int, s *goquery.Selection) {
		textContent := strings.TrimSpace(s.Text())
		if textContent == "" && len(s.Children().Nodes) < 1 {
			s.Remove()
		}
	})
	document.Find("img").Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr("src")
		if textMatch(removeImgSrcMatch, src, "exact") != -1 {
			s.Remove()
		}

	})

}

func removeEmptyTagsSecond(document *goquery.Document) {
	document.Find("a").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		item := s.Get(0)
		if len(text) == 0 && item.FirstChild == nil {
			s.Remove()
		}
	})
	document.Find("p,section,strong").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if len(text) == 0 && len(s.Children().Nodes) < 1 {
			s.Remove()
		}
	})
}

func removeExtractText(document *goquery.Document) {

	document.Find("section").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if textMatch(removeTextMatch, text, "exact") != -1 {
			s.Remove()
		}

	})

}

func (t PostExtractorTemplate) WechatPostExtractor(content, feedUrl string) string {
	//content = strings.ReplaceAll(content, "&nbsp;", "")
	templateData := strings.NewReader(content)

	doc, _ := goquery.NewDocumentFromReader(templateData)
	doc.Find("svg,mp-common-profile").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	removeEmptyTagsFirst(doc)
	nodes := make([]*goquery.Selection, 0)
	parentSelection := doc.Find("div").First()
	if len(parentSelection.Nodes) == 1 {
		parentSelection = parentSelection.Children()
	}
	if len(parentSelection.Nodes) == 1 {
		parentSelection = parentSelection.Children()
	}
	if parentSelection != nil {
		parentSelection.Children().Each(func(i int, child *goquery.Selection) {
			nodes = append(nodes, child)
		})
	}

	if parentSelection.Parent() != nil && len(parentSelection.Parent().Nodes) > 0 {
		parentNode := parentSelection.Parent().Get(0)
		parentNode.Attr = append(parentNode.Attr, html.Attribute{
			Key: "style",
			Val: "line-height: 1.8em",
		})
	}
	/*for _, node := range parentSelection.Nodes {
		if node.Data == "section" {
			node.Attr = append(node.Attr, html.Attribute{
				Key: "style",
				Val: "margin-bottom:14px",
			})
		}
	}*/

	_, validNum := checkArticleEnd(0, nodes)
	checkImgNo := validNum - 1

	if len(nodes) > 0 {
		removeEmptyTagsSecond(doc)
		for ; checkImgNo >= 0; checkImgNo-- {
			node := nodes[checkImgNo]
			if checkIsImgNode(node) {
				node.Remove()
			} else {
				break
			}
		}
		removeEmptyTagsSecond(doc)
		removeExtractText(doc)

	} else {
		log.Printf(`wechat conent empty %s`, content)
	}
	newContent, _ := doc.Html()
	return newContent
}
