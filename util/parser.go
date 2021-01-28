package util

import (
	"bytes"
	"log"
	"net/url"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// GetURLs returns the urls in an html string
func GetURLs(htmlStr []byte) ([][]byte, error) {
	r, err := regexp.Compile(`href="[^\"]*"`)
	if err != nil {
		return nil, err
	}
	return r.FindAll(htmlStr, -1), nil
}

// GetAnchorHref Method
func GetAnchorHref(htmlStr []byte) ([][]byte, error) {
	hm := make(map[string]struct{})
	doc, err := html.Parse(bytes.NewBuffer(htmlStr))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					hm[a.Val] = struct{}{}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	res := [][]byte{}
	for k := range hm {
		nk := strings.ReplaceAll(k, ",", "")
		res = append(res, []byte(nk))
	}
	return res, nil
}

// IsURLInSubdomain returns true if a url is in a subdomain
func IsURLInSubdomain(subdomain, link string) (string, bool) {
	u, err := url.Parse(link)
	if err != nil {
		log.Println(err)
		return "", false
	}
	if u.Hostname() == "" || u.Hostname() == subdomain {
		nu, err := url.Parse(u.Path)
		if err != nil {
			return "", false
		}
		nu.Scheme = "http"
		nu.Host = subdomain
		return nu.String(), true
	}
	return "", false
}
