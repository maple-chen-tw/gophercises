package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"golang.org/x/net/html"
)

var (
	link        string
	host        string
	xmlFileName string
)

func init() {
	flag.StringVar(&link, "link", "http://localhost:3030/", "a web link for testing")
	flag.StringVar(&xmlFileName, "sitemap.xml", "http://localhost:3030/", "xmlFileName")
}
func main() {
	flag.Parse()
	host = getHost(link)

	hrefs := bfs()
	sitemap := generateSitemap(hrefs)
	if err := writeSitemapToFile(xmlFileName, sitemap); err != nil {
		fmt.Println("Error writing sitemap to file:", err)
		return
	}

	fmt.Println("Sitemap written to", xmlFileName)
}

// 把http頁面抓下來
func getPage(link string) string {
	page, err := http.Get(link)
	if err != nil {
		log.Fatal(err)
	}
	content, err := io.ReadAll(page.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

// link解析
//type Link struct {
//	Href string
//	Text string
//}

func getLink(n *html.Node) string {

	var href string
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				href = a.Val
				break
			}
		}
		//link.Text = renderInnerContent(n)
	}

	return href
}

func traverse(n *html.Node) []string {
	var links []string
	if n.Type == html.ElementNode {
		href := getLink(n)
		if href != "" {
			links = append(links, href)
			//fmt.Printf("Link: %s\nText: %s\n", link.Href, link.Text)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		childLinks := traverse(c)
		links = append(links, childLinks...)
	}
	return links
}

// 拿到host
func getHost(href string) string {

	url, err := url.Parse(href)
	if err != nil {
		log.Fatal(err)
	}
	if url.Scheme == "" {
		return host
	}

	return url.Scheme + "://" + strings.TrimPrefix(url.Host, "www.") + "/"
}

// 判斷符合的domain link
func isUrlValid(href string) bool {
	hrefHost := getHost(href)
	return hrefHost == host
}

func toAbsoluteUrl(href string) string {
	parsedUrl, err := url.Parse(href)
	if err != nil {
		log.Fatal(err)
	}
	if parsedUrl.Scheme == "" {
		return host + strings.TrimPrefix(href, "/") // Trim leading slash if needed
	}
	return href
}

//對每個domain link抓頁面
//再處理link
//再加...

func bfs() []string {
	queue := []string{link}
	visited := map[string]bool{}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if visited[current] {
			continue
		}
		url := toAbsoluteUrl(current)
		visited[url] = true
		content := getPage(url)

		doc, err := html.Parse(strings.NewReader(content))
		if err != nil {
			log.Fatal(err)
		}
		links := traverse(doc)
		for _, val := range links {
			if isUrlValid(val) && !visited[val] {
				queue = append(queue, val)
			}
		}
	}
	var visitedUrls []string
	for u := range visited {
		visitedUrls = append(visitedUrls, u)
	}
	return visitedUrls
}

func generateSitemap(hrefs []string) string {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	builder.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")

	for _, href := range hrefs {
		builder.WriteString("  <url>\n")
		builder.WriteString(fmt.Sprintf("    <loc>%s</loc>\n", href))
		builder.WriteString("  </url>\n")

	}
	builder.WriteString(`</urlset>`)
	return builder.String()
}

func writeSitemapToFile(filename string, sitemap string) error {
	return os.WriteFile(filename, []byte(sitemap), 0644)
}
