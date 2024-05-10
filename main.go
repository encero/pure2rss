package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"log/slog"

	"github.com/encero/pure2rss/src/pure2rss"
	"github.com/gorilla/feeds"
)

func main() {
	cache, err := pure2rss.NewCache("./posts.json")
	if err != nil {
		panic(err)
	}

	err = downloadPostInfo(cache)
	if err != nil {
		panic(err)
	}

	err = cache.Persist()
	if err != nil {
		panic(err)
	}

	posts, err := cache.LoadCategory("purely-technical")
    if err != nil {
        panic(err)
    }

    slices.SortFunc(posts, func(a,b pure2rss.Post) int {
        return b.PostLink.Link.LastMod.Compare(a.PostLink.Link.LastMod)
    })

    lastPost := posts[len(posts) - 1]



    feed := feeds.Feed{
    	Title:       "Purely Technical",
        Link:        &feeds.Link{Href:"https://blog.purestorage.com/purely-technical/"},
    	Description: "A wealth of technical detail for working with Pure products.",
    	Updated:     lastPost.PostLink.Link.LastMod,
        Items: make([]*feeds.Item, 0, len(posts)),
    }

    for _, post := range posts {
        feed.Items = append(feed.Items, &feeds.Item{
        	Title:       post.Title,
            Link:        &feeds.Link{Href: post.PostLink.Link.Loc},
        	Description: post.Summary,
        	Updated:     post.PostLink.Link.LastMod,
        })
    }

    file, err := os.Create("./rss.xml")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    feed.WriteRss(file)
}

func downloadPostInfo(cache *pure2rss.Cache) error {
	crawler := pure2rss.NewCrawler("https://blog.purestorage.com/sitemap_index.xml", pure2rss.CrawlerTimeout(time.Second*30))

	crawler.OnIndexLink(func(l pure2rss.Link) bool {
		return strings.Contains(l.Loc, "post-sitemap")
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	linkCh := make(chan pure2rss.PostLink)

	wg := &sync.WaitGroup{}
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go postDownloader(ctx, wg, client, linkCh, cache)
	}

	crawler.OnPostLink(func(l pure2rss.Link) {
		postLink, err := pure2rss.ParsePostLink(l)
		if err != nil {
			slog.Error("Parsing link info", slog.String("url", l.Loc), slog.Any("err", err))
			return
		}

		if postLink.Lang != "en" {
			return
		}

		if postLink.Category != "purely-technical" {
			return
		}
		cached, err := cache.Load(l.Loc)
		if err != nil && !errors.Is(err, pure2rss.NoDataError) {
			return
		}

		if err == nil && (cached.PostLink.Link.LastMod.After(l.LastMod) || cached.PostLink.Link.LastMod.Equal(l.LastMod)) {
			slog.Info("skipping post, already in cache", slog.String("url", l.Loc))
			return
		}

		linkCh <- postLink
	})

	go crawler.Run()

	select {
	case err := <-crawler.Done():
		if err != nil {
			return fmt.Errorf("Crawler failed, %w", err)
		}
	}

	cancel()
	wg.Wait()

	return nil
}

func postDownloader(ctx context.Context, wg *sync.WaitGroup, client *http.Client, linkCh chan pure2rss.PostLink, cache *pure2rss.Cache) {
	defer wg.Done()

	for {
		var work pure2rss.PostLink
		select {
		case work = <-linkCh:
		case <-ctx.Done():
			return
		}

		post, err := pure2rss.FetchAndAndParsePost(client, work)
		if err != nil {
			slog.Error("parsing post content", slog.Any("err", err), slog.String("url", work.Link.Loc))
		}

		slog.Info("storing post", slog.String("title", post.Title))

		cache.Store(post)
	}
}
