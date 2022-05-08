package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"link"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Loc string `xml:"loc"`
}

type urlset struct {
	Xmlns string `xml:"xmlns,attr"`
	Urls  []loc  `xml:"url"`
}

func main() {
	urlFlag := flag.String("url", "https://gophercises.com", "the url that you want to build the sitemap for")
	maxDepth := flag.Int("depth", 3, "max number of links deep to traverse")
	flag.Parse()
	fmt.Println(*maxDepth)
	toXML := urlset{
		Xmlns: xmlns,
	}
	pages := bfs(*urlFlag, *maxDepth)
	for _, page := range pages {
		toXML.Urls = append(toXML.Urls, loc{page})
	}
	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", " ")
	if err := enc.Encode(toXML); err != nil {
		panic(err)
	}
	fmt.Println()
}

func bfs(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: {},
	}
	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})
		for url := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			for _, link := range get(url) {
				nq[link] = struct{}{}
			}
		}
	}
	ret := make([]string, 0, len(seen))
	for url := range seen {
		ret = append(ret, url)
	}
	return ret
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parse(r)
	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()
	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()
	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func filter(links []string, keepFn func(link string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(pfx string) func(link string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}
