package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/monzo/webcrawler/service"
	"github.com/monzo/webcrawler/store"
	"github.com/monzo/webcrawler/util"
)

// Crawler struct
type Crawler struct {
	p     service.Producer
	cache *store.Redis
}

var upgrader = websocket.Upgrader{}

// NewCrawler returns a Crawler struct
func NewCrawler(p service.Producer, cache *store.Redis) *Crawler {
	return &Crawler{p: p, cache: cache}
}

// Crawl a given URL
func (c *Crawler) Crawl(w http.ResponseWriter, r *http.Request) {

	// Upgrade connection to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	messageType, p, err := conn.ReadMessage()
	if err != nil {
		log.Println(err)
		return
	}
	// publish to queue for crawling
	err = c.p.Publish(context.Background(), p)
	if err != nil {
		log.Println(err)
		conn.WriteMessage(messageType, []byte(err.Error()))
		return
	}

	c.bfs(string(p), conn, messageType)
}

// BFS method
func (c *Crawler) bfs(urlStr string, conn *websocket.Conn, messageType int) {
	hm := map[string]string{}
	visited := map[string]bool{}
	domain, err := url.Parse(urlStr)
	if err != nil {
		log.Println(err)
		return
	}
	subdomain := domain.Hostname()
	if subdomain == "" {
		log.Println("unknown domain")
		return
	}
	nurl, err := normalizeURL(subdomain, domain.Path)
	if err != nil {
		log.Println("normalize error", err)
		return
	}
	queue := util.NewFifo()
	queue.Push(nurl)
	fmt.Println("BFS method called...", nurl, queue.Length())
	for {
		if queue.IsEmpty() {
			hm["key"] = "Done"
			hm["value"] = "Done"
			conn.WriteJSON(hm)
			return
		}
		p := queue.Front()
		if visited[p] {
			continue
		}
		v, err := c.cache.Fetch(p)
		if err != nil {
			log.Println("fetched nothing from cache", err)
		}
		if err == nil && v != "" {
			visited[p] = true
			hm["key"] = p
			hm["value"] = v
			conn.WriteJSON(hm)
			time.Sleep(3 * time.Second)
			newUrls := strings.Split(v, ",")
			for _, v := range newUrls {
				nu, b := util.IsURLInSubdomain(subdomain, v)
				if b {
					queue.Push(nu)
				}
			}
		} else {
			queue.Push(p)
			time.Sleep(1 * time.Second)
		}
	}
}

func normalizeURL(domain, path string) (string, error) {
	nurl, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	nurl.Host = domain
	nurl.Scheme = "http"
	return nurl.String(), nil
}
