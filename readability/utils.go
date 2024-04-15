package readability

import (
	"bytes"
	"log"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

// QuerySelectorAll returns array of document's elements that match
// the specified group of selectors.
func querySelectorAll(doc *html.Node, selectors string) []*html.Node {
	matcher, err := cascadia.ParseGroup(selectors)
	if err != nil {
		return nil
	}

	return cascadia.QueryAll(doc, matcher)
}

// QuerySelector returns the first document's element that match
// the specified group of selectors.
func querySelector(doc *html.Node, selectors string) *html.Node {
	matcher, err := cascadia.ParseGroup(selectors)
	if err != nil {
		return nil
	}

	return cascadia.Query(doc, matcher)
}

// IsVoidElement check whether a node can have any contents or not.
// Return true if element is void (can't have any children).
func isVoidElement(n *html.Node) bool {
	// If it's not element, it's void
	if n.Type != html.ElementNode {
		return true
	}

	// Check tag name
	switch n.Data {
	case "area", "base", "br", "col", "embed", "hr",
		"img", "input", "keygen", "link", "meta",
		"param", "source", "track", "wbr":
		return true
	default:
		return false
	}
}

func detachChild(child *html.Node) {
	if child.Parent != nil || child.PrevSibling != nil || child.NextSibling != nil {
		if child.Parent != nil {
			if child.Parent.FirstChild == child {
				child.Parent.FirstChild = child.NextSibling
			}

			if child.Parent.LastChild == child {
				child.Parent.LastChild = child.PrevSibling
			}
		}

		if child.PrevSibling != nil {
			child.PrevSibling.NextSibling = child.NextSibling
		}

		if child.NextSibling != nil {
			child.NextSibling.PrevSibling = child.PrevSibling
		}

		child.Parent = nil
		child.PrevSibling = nil
		child.NextSibling = nil
	}
}

// ReplaceChild replaces a child node within the given (parent) node.
// If the new child is already exist in document, ReplaceChild() will move it
// from its current position to replace old child. Returns both the new and old child.
//
// TODO: note sure but I *think* there are some issues here. Check later I guess.
func replaceChild(parent *html.Node, newChild *html.Node, oldChild *html.Node) (*html.Node, *html.Node) {
	// Make sure parent is specified and not void
	if parent == nil && !isVoidElement(parent) {
		return newChild, oldChild
	}

	// Make sure the specified parent IS the parent of the old child
	if oldChild.Parent != parent {
		return newChild, oldChild
	}

	// Detach the new child
	detachChild(newChild)
	parent.InsertBefore(newChild, oldChild)
	parent.RemoveChild(oldChild)
	return newChild, oldChild

}

func parseHtmlWithContext(h string, context *html.Node) []*html.Node {
	nodes, err := html.ParseFragment(strings.NewReader(h), context)
	if err != nil {
		panic("failed to parse HTML: " + err.Error())
	}
	return nodes
}

func replaceWithHtml(node *html.Node, htmlStr string) {
	parent := node.Parent
	if parent == nil {
		return
	}
	nodes := parseHtmlWithContext(htmlStr, parent)
	nextSibling := node.NextSibling
	for _, n := range nodes {
		if node.Parent != nil {
			node.Parent.InsertBefore(n, nextSibling)
		}
	}
	parent.RemoveChild(node)
}

// firstElementChild returns the object's first child Element, or nil if there
// are no child elements.
func firstElementChild(node *html.Node) *html.Node {
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			return child
		}
	}

	return nil
}

// nextElementSibling returns the Element immediately following the specified
// one in its parent's children list, or nil if the specified Element is the
// last one in the list.
func nextElementSibling(node *html.Node) *html.Node {
	for sibling := node.NextSibling; sibling != nil; sibling = sibling.NextSibling {
		if sibling.Type == html.ElementNode {
			return sibling
		}
	}

	return nil
}

// appendChild adds a node to the end of the list of children of a specified
// parent node. If the given child is a reference to an existing node in the
// document, appendChild moves it from its current position to the new position
// (there is no requirement to remove the node from its parent node before
// appending it to some other node).
//
// See: https://developer.mozilla.org/en-US/docs/Web/API/Node/appendChild
func appendChild(node *html.Node, child *html.Node) {
	if child.Parent != nil {
		temp := cloneNode(child)
		node.AppendChild(temp)
		child.Parent.RemoveChild(child)
		return
	}

	node.AppendChild(child)
}

// childNodes returns list of a node's direct children.
func childNodes(node *html.Node) []*html.Node {
	var list []*html.Node

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		list = append(list, c)
	}

	return list
}

// includeNode determines if node is included inside nodeList.
func includeNode(nodeList []*html.Node, node *html.Node) bool {
	for i := 0; i < len(nodeList); i++ {
		if nodeList[i] == node {
			return true
		}
	}

	return false
}

// cloneNode returns a duplicate of the node on which this method was called.
//
// See: https://developer.mozilla.org/en-US/docs/Web/API/Node/cloneNode
func cloneNode(node *html.Node) *html.Node {
	clone := &html.Node{
		Type:     node.Type,
		DataAtom: node.DataAtom,
		Data:     node.Data,
		Attr:     make([]html.Attribute, len(node.Attr)),
	}

	copy(clone.Attr, node.Attr)

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		clone.AppendChild(cloneNode(c))
	}

	return clone
}

// createElement creates the HTML element specified by tagName.
//
// See: https://developer.mozilla.org/en-US/docs/Web/API/Document/createElement
func createElement(tagName string) *html.Node {
	return &html.Node{Type: html.ElementNode, Data: tagName}
}

// createTextNode creates a new Text node.
func createTextNode(data string) *html.Node {
	return &html.Node{Type: html.TextNode, Data: data}
}

// getElementsByTagName returns a collection of HTML elements with the given
// tag name. If tag name is an asterisk, a list of all the available HTML nodes
// will be returned instead.
//
// See: https://developer.mozilla.org/en-US/docs/Web/API/Document/getElementsByTagName
func getElementsByTagName(node *html.Node, tag string) []*html.Node {
	var lst []*html.Node
	var fun func(*html.Node)

	fun = func(n *html.Node) {
		if n.Type == html.ElementNode && (tag == "*" || n.Data == tag) {
			lst = append(lst, n)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			fun(c)
		}
	}

	fun(node)

	return lst
}

// getAttribute returns the value of a specified attribute on the element. If
// the given attribute does not exist, the function returns an empty string.
func getAttribute(node *html.Node, attrName string) string {
	for i := 0; i < len(node.Attr); i++ {
		if node.Attr[i].Key == attrName {
			return node.Attr[i].Val
		}
	}

	return ""
}

// setAttribute sets attribute for node. If attribute already exists, it will
// be replaced.
func setAttribute(node *html.Node, attrName string, attrValue string) {
	attrIdx := -1
	for i := 0; i < len(node.Attr); i++ {
		if node.Attr[i].Key == attrName {
			attrIdx = i
			break
		}
	}

	if attrIdx >= 0 {
		node.Attr[attrIdx].Val = attrValue
	} else {
		node.Attr = append(node.Attr, html.Attribute{
			Key: attrName,
			Val: attrValue,
		})
	}
}

// removeAttribute removes attribute with given name.
func removeAttribute(node *html.Node, attrName string) {
	attrIdx := -1

	for i := 0; i < len(node.Attr); i++ {
		if node.Attr[i].Key == attrName {
			attrIdx = i
			break
		}
	}

	if attrIdx >= 0 {
		a := node.Attr
		a = append(a[:attrIdx], a[attrIdx+1:]...)
		node.Attr = a
	}
}

// hasAttribute returns a Boolean value indicating whether the specified node
// has the specified attribute or not.
func hasAttribute(node *html.Node, attrName string) bool {
	for i := 0; i < len(node.Attr); i++ {
		if node.Attr[i].Key == attrName {
			return true
		}
	}

	return false
}

// outerHTML returns an HTML serialization of the element and its descendants.
func outerHTML(node *html.Node) string {
	var buffer bytes.Buffer

	if err := html.Render(&buffer, node); err != nil {
		return ""
	}

	return buffer.String()
}

// innerHTML returns the HTML content (inner HTML) of an element.
func innerHTML(node *html.Node) string {
	var err error
	var buffer bytes.Buffer

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if err = html.Render(&buffer, child); err != nil {
			return ""
		}
	}

	return strings.TrimSpace(buffer.String())
}

// documentElement returns the root element of the document.
func documentElement(doc *html.Node) *html.Node {
	nodes := getElementsByTagName(doc, "html")

	if len(nodes) > 0 {
		return nodes[0]
	}

	return nil
}

// className returns the value of the class attribute of the element.
func className(node *html.Node) string {
	className := getAttribute(node, "class")
	className = strings.TrimSpace(className)
	className = rxNormalize.ReplaceAllString(className, "\x20")
	return className
}

func isProbablyNavigation(node *html.Node) bool {
	navRelatedClasses := regexp.MustCompile(`next|prev|previous`)
	name := tagName(node)
	if name != "OL" && name != "UL" {
		return false
	}
	children := getElementsByTagName(node, "li")
	for _, child := range children {
		if navRelatedClasses.MatchString(className(child)) && getElementsByTagName(node, "a") != nil {
			return true
		}
	}
	return false
}

// id returns the value of the id attribute of the specified element.
func id(node *html.Node) string {
	id := getAttribute(node, "id")
	id = strings.TrimSpace(id)
	return id
}

// children returns an HTMLCollection of the child elements of Node.
func children(node *html.Node) []*html.Node {
	var children []*html.Node

	if node == nil {
		return nil
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode {
			children = append(children, child)
		}
	}

	return children
}

// wordCount returns number of word in str.
func wordCount(str string) int {
	return len(strings.Fields(str))
}

// indexOf returns the first index at which a given element can be found in the
// array, or -1 if it is not present.
func indexOf(array []string, key string) int {
	for idx, val := range array {
		if val == key {
			return idx
		}
	}

	return -1
}

// tagName returns the tag name of the element on which it’s called.
//
// For example, if the element is an <img>, its tagName property is “IMG” (for
// HTML documents; it may be cased differently for XML/XHTML documents).
//
// See: https://developer.mozilla.org/en-US/docs/Web/API/Element/tagName
func tagName(node *html.Node) string {
	if node.Type != html.ElementNode {
		return ""
	}

	return node.Data
}

// textContent returns text content of a Node and its descendants.
//
// See: https://developer.mozilla.org/en-US/docs/Web/API/Node/textContent
func textContent(node *html.Node) string {
	var buffer bytes.Buffer
	var finder func(*html.Node)

	finder = func(n *html.Node) {
		if n.Type == html.TextNode {
			buffer.WriteString(n.Data)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			finder(c)
		}
	}

	finder(node)

	return buffer.String()
}

// PreviousElementSibling returns the the Element immediately prior
// to the specified one in its parent's children list, or null if
// the specified element is the first one in the list.
func previousElementSibling(node *html.Node) *html.Node {
	for sibling := node.PrevSibling; sibling != nil; sibling = sibling.PrevSibling {
		if sibling.Type == html.ElementNode {
			return sibling
		}
	}
	return nil
}

// toAbsoluteURI convert uri to absolute path based on base.
// However, if uri is prefixed with hash (#), the uri won't be changed.
func toAbsoluteURI(uri string, base *url.URL) string {
	if uri == "" || base == nil {
		return ""
	}

	// If it is hash tag, return as it is
	if uri[:1] == "#" {
		return uri
	}

	// If it is already an absolute URL, return as it is
	tmp, err := url.ParseRequestURI(uri)
	if err == nil && tmp.Scheme != "" && tmp.Hostname() != "" {
		return uri
	}

	// Otherwise, resolve against base URI.
	tmp, err = url.Parse(uri)
	if err != nil {
		return uri
	}

	return base.ResolveReference(tmp).String()
}

func sliceToMap(strings ...string) map[string]struct{} {
	result := make(map[string]struct{})
	for _, s := range strings {
		result[s] = struct{}{}
	}
	return result
}

// charCount returns number of char in str.
func charCount(str string) int {
	return utf8.RuneCountInString(str)
}

// isValidURL checks if URL is valid.
func IsValidURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

// strOr returns the first not empty string in args.
func strOr(args ...string) string {
	for i := 0; i < len(args); i++ {
		if args[i] != "" {
			return args[i]
		}
	}
	return ""
}

func strFilter(strs []string, filter func(string) bool) []string {
	var result []string
	for _, s := range strs {
		if filter(s) {
			result = append(result, s)
		}
	}
	return result
}

func trim(s string) string {
	s = strings.Join(strings.Fields(s), " ")
	return strings.TrimSpace(s)
}

func InsertToFile(filePath string, content string) {

	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal("open file error", filePath, ": ", err)
	}
	defer file.Close()

	dataToWrite := []byte(content)
	_, err = file.Write(dataToWrite)
	if err != nil {
		log.Fatal("write data: ", err)
	}
}

func parseInt(value string, defaultV int) int {
	if value == "" {
		return defaultV
	}

	v, err := strconv.Atoi(value)
	if err != nil {
		return defaultV
	}
	return v
}
