package pure2rss

import (
	"fmt"
	"net/http"
)

type Crawler struct {
	sitemapURL string
	done       chan error
	onIndex    func([]Link) []Link
}

func NewCrawler(sitemapURL string) *Crawler {
	return &Crawler{
		sitemapURL: sitemapURL,
		done:       make(chan error),
	}
}

func (c *Crawler) Run() {
    response, err := http.Get(c.sitemapURL)
    if err != nil {
        c.done <- fmt.Errorf("loading sitemap index, %w", err)
        return
    }
    defer response.Body.Close()

    ParseSiteMapList(response.Body)

    c.onIndex([]Link{})

	c.done <- nil
}

func (c *Crawler) OnIndex(f func([]Link) []Link) {
	c.onIndex = f
}

func (c *Crawler) Done() <-chan struct{} {
	return c.done
}
