package pure2rss

import (
	"fmt"
	"net/http"
	"time"
)

type Crawler struct {
	client *http.Client

	sitemapURL  string
	done        chan error
	onIndexLink func(Link) bool
	onPostLink  func(Link)
}

type crawlerOption func (c *Crawler)

func CrawlerTimeout(d time.Duration) crawlerOption {
    return func(c *Crawler) {
        c.client.Timeout = d
    }
}

func NewCrawler(sitemapURL string, options ...crawlerOption) *Crawler {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	return &Crawler{
		client:      client,
		sitemapURL:  sitemapURL,
		done:        make(chan error),
		onIndexLink: func(l Link) bool { return true },
		onPostLink:  func(l Link) { return },
	}
}

func (c *Crawler) fetchAndParseSitemapIndex() ([]Link, error) {
	req, err := http.NewRequest(http.MethodGet, c.sitemapURL, nil)
	if err != nil {
		return nil, fmt.Errorf("constructing site map index request, %w", err)
	}

	response, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("loading sitemap index, %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %q when fetching sitemap index", response.Status)
	}

	links, err := ParseSiteMapList(response.Body)
	if err != nil {
		return nil, fmt.Errorf("parsing sitemap list, %w", err)
	}

	return links, nil
}

func (c *Crawler) Run() {
	links, err := c.fetchAndParseSitemapIndex()
	if err != nil {
		c.done <- fmt.Errorf("fetching and parsing site map index, %w", err)
		return
	}

	for _, sitemapLink := range links {
		if !c.onIndexLink(sitemapLink) {
			continue
		}

		req, err := http.NewRequest(http.MethodGet, sitemapLink.Loc, nil)
		if err != nil {
			c.done <- fmt.Errorf("constructing post sitemap request, %w", err)
		}

		sitemapResponse, err := c.client.Do(req)
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
            if link.Loc == "https://blog.purestorage.com/" {
                continue
            }

			c.onPostLink(link)
		}
	}

	c.done <- nil
}

func (c *Crawler) OnIndexLink(f func(Link) bool) {
	c.onIndexLink = f
}

func (c *Crawler) OnPostLink(f func(Link)) {
	c.onPostLink = f
}

func (c *Crawler) Done() <-chan error {
	return c.done
}
