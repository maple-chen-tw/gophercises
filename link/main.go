package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

var htmlFile string

func init() {
	flag.StringVar(&htmlFile, "htmlFile", "ex4.html", "a example html file for testing")
}

func getLink(n *html.Node) Link {

	var link Link
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				link.Href = a.Val
				break
			}
		}
		link.Text = renderInnerContent(n)
	}

	return link
}

func renderInnerContent(n *html.Node) string {
	var result strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode {
			result.WriteString(c.Data)
		} else if c.Type == html.ElementNode {
			result.WriteString(renderInnerContent(c))
		}
	}
	return result.String()
}

func traverse(n *html.Node) {
	if n.Type == html.ElementNode {
		link := getLink(n)
		if link.Href != "" {
			fmt.Printf("Link: %s\nText: %s\n", link.Href, link.Text)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverse(c)
	}
}

func main() {
	flag.Parse()
	file, err := os.Open(htmlFile)
	if err != nil {
		log.Fatalf("Unable to read file: %s", err)
	}
	doc, err := html.Parse(file)
	if err != nil {
		log.Fatal(err)
	}
	traverse(doc)
}
