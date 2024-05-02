package pure2rss

import (
	"fmt"
	"net/http"
)

type Crawler struct {
	sitemapURL string
	done       chan error
	onIndexLink    func(Link) bool
	onPostLink func(Link) bool
}

func NewCrawler(sitemapURL string) *Crawler {
	return &Crawler{
		sitemapURL: sitemapURL,
		done:       make(chan error),
		onIndexLink:    func(l Link) bool {return true},
		onPostLink: func(l Link) bool { return true },
	}
}

func (c *Crawler) Run() {
	response, err := http.Get(c.sitemapURL)
	if err != nil {
		c.done <- fmt.Errorf("loading sitemap index, %w", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		c.done <- fmt.Errorf("unexpected status code %q when fetching sitemap index", response.Status)
		return
	}

	links, err := ParseSiteMapList(response.Body)
	if err != nil {
		c.done <- fmt.Errorf("parsing sitemap list, %w", err)
		return
	}

	for _, sitemapLink := range links {
	    if !c.onIndexLink(sitemapLink) {
            continue
        }

		sitemapResponse, err := http.Get(sitemapLink.Loc)
		if err != nil {
			c.done <- fmt.Errorf("fetching sitemap %q failed, %w", sitemapLink.Loc, err)
			return
		}
		defer sitemapResponse.Body.Close()

		postLinks, err := ParseSiteMap(sitemapResponse.Body)
		if err != nil {
			c.done <- fmt.Errorf("parsing posts sitemap %q, %w", sitemapLink.Loc, err)
			return
		}

		for _, link := range postLinks {
			c.onPostLink(link)
		}
	}

	c.done <- nil
}

func (c *Crawler) OnIndexLink(f func(Link) bool) {
	c.onIndexLink = f
}

func (c *Crawler) OnPostLink(f func(Link) bool) {
	c.onPostLink = f
}

func (c *Crawler) Done() <-chan error {
	return c.done
}
