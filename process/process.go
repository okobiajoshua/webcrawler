package process

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/monzo/webcrawler/data"
	"github.com/monzo/webcrawler/service"
	"github.com/monzo/webcrawler/store"
	"github.com/monzo/webcrawler/util"
)

// Process struct
type Process struct {
	data  data.Data
	pub   service.Producer
	cache store.DataStore
}

// NewProcess function
func NewProcess(data data.Data, pub service.Producer, cache store.DataStore) *Process {
	return &Process{data: data, pub: pub, cache: cache}
}

// Process function
func (p *Process) Process(url []byte) error {
	v, err := p.cache.Fetch(string(url))
	if err != nil {
		return err
	}
	if v != "" {
		return nil
	}
	subdomain, err := getSubDomain(url)
	if err != nil {
		return err
	}
	htmlStr, err := p.data.GetHTML(url)
	if err != nil {
		log.Println("fetching page content", err)
		return err
	}
	urls, err := util.GetAnchorHref(htmlStr)
	if err != nil {
		log.Println("Error: ", err)
		return err
	}
	err = p.cache.Save(string(url), byteArrayToString(urls))
	if err != nil {
		log.Println("error saving to cache", err)
	}
	for _, v := range urls {
		if newURL, ok := util.IsURLInSubdomain(subdomain, string(v)); ok {
			v, err := p.cache.Fetch(newURL)
			if err != nil {
				log.Println("error fetching from cache", err)
			} else if v == "" {
				p.pub.Publish(context.TODO(), []byte(newURL))
			}
		}
	}
	return nil
}

func byteArrayToString(inp [][]byte) string {
	res := []string{}
	for _, v := range inp {
		res = append(res, string(v))
	}

	return strings.Join(res, ",")
}

func getSubDomain(uri []byte) (string, error) {
	xURI, err := url.Parse(string(uri))
	if err != nil {
		return "", err
	}
	if strings.Trim(xURI.Hostname(), " ") == "" {
		return "", fmt.Errorf("unknown domain")
	}
	return strings.Trim(xURI.Hostname(), " "), nil
}
