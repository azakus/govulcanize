package htmlutils

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
)

// Attr returns the value of the attribute. ok indicates if the attribute
// exists for the given node.
func Attr(n *html.Node, attrKey string) (val string, ok bool) {
	for _, attr := range n.Attr {
		if attr.Key == attrKey {
			return attr.Val, true
		}
	}
	return "", false
}

// SetAttr sets the value of the attribute for the given node
func SetAttr(n *html.Node, attrKey string, attrValue string) {
	for i, attr := range n.Attr {
		if attr.Key == attrKey {
			attr.Val = attrValue
			n.Attr[i] = attr
			return
		}
	}
	n.Attr = append(n.Attr, html.Attribute{
		Key: attrKey,
		Val: attrValue,
	})
}

func RemoveAttr(n *html.Node, attrKey string) {
	for i := 0; i < len(n.Attr); i++ {
		attr := n.Attr[i]
		if attr.Key == attrKey {
			newAttrs := make([]html.Attribute, len(n.Attr)-1)
			copy(newAttrs, n.Attr[:i])
			copy(newAttrs[i:], n.Attr[i+1:])
			n.Attr = newAttrs
			i--
		}
	}
}

// TextContent returns the text within the given node
func TextContent(n *html.Node) string {
	child := n.FirstChild
	if child.Type == html.TextNode {
		return child.Data
	}
	return ""
}

// SetTextContent sets the text within the given node
func SetTextContent(n *html.Node, text string) {
	child := n.FirstChild
	if child.Type == html.TextNode {
		child.Data = text
	}
}

// GetElementByID returns the element with the given id, if one exists
func GetElementByID(doc *html.Node, id string) *html.Node {
	matches := Search(doc, func(n *html.Node) bool {
		nodeid, ok := Attr(n, "id")
		return ok && nodeid == id
	})
	if len(matches) == 1 {
		return matches[0]
	}
	return nil
}

func RemoveNode(doc *Fragment, n *html.Node) {
	if n.PrevSibling != nil {
		n.PrevSibling.NextSibling = n.NextSibling
	}
	if n.NextSibling != nil {
		n.NextSibling.PrevSibling = n.PrevSibling
	}
	if doc.FirstNode == n {
		doc.FirstNode = n.NextSibling
	}
	if doc.LastNode == n {
		doc.LastNode = n.PrevSibling
	}
}

func ReplaceNodeWithNode(doc *Fragment, origNode *html.Node, newNode *html.Node) {
	ReplaceNodeWithFragment(doc, origNode, FromNode(newNode))
}

func ReplaceNodeWithFragment(doc *Fragment, node *html.Node, fragment *Fragment) {
	if fragment.Length() == 0 {
		RemoveNode(doc, node)
		return
	}

	// Set all new nodes' parent
	fragment.eachNode(func(n *html.Node) {
		n.Parent = node.Parent
	})

	// Insert into linked list
	fragment.FirstNode.PrevSibling = node.PrevSibling
	fragment.LastNode.NextSibling = node.NextSibling
	if node.PrevSibling != nil {
		node.PrevSibling.NextSibling = fragment.FirstNode
	}
	if node.NextSibling != nil {
		node.NextSibling.PrevSibling = fragment.LastNode
	}

	// Update parent
	if node.Parent != nil {
		if node.Parent.FirstChild == node {
			node.Parent.FirstChild = fragment.FirstNode
		}
		if node.Parent.LastChild == node {
			node.Parent.LastChild = fragment.LastNode
		}
	}

	// Update doc
	if doc.FirstNode == node {
		doc.FirstNode = fragment.FirstNode
	}
	if doc.LastNode == node {
		doc.LastNode = fragment.LastNode
	}
}

func CreateScript(content string) *html.Node {
	script := &html.Node{
		Type:     html.ElementNode,
		Data:     "script",
		DataAtom: atom.Script,
	}
	textnode := &html.Node{
		Type: html.TextNode,
		Data: content,
	}
	script.AppendChild(textnode)
	return script
}

func CreateExternalScript(src string) *html.Node {
	script := &html.Node{
		Type:     html.ElementNode,
		Data:     "script",
		DataAtom: atom.Script,
		Attr: []html.Attribute{
			html.Attribute{Key: "src", Val: src},
		},
	}
	return script
}

func CreateStyle(content string) *html.Node {
	style := &html.Node{
		Type:     html.ElementNode,
		Data:     "style",
		DataAtom: atom.Style,
	}
	textnode := &html.Node{
		Type: html.TextNode,
		Data: content,
	}
	style.AppendChild(textnode)
	return style
}
