package main

import (
	"context"
	"net/http"
	"strings"
	"sync"
	"time"

	"log/slog"

	"github.com/encero/pure2rss/src/pure2rss"
)

func main() {
	crawler := pure2rss.NewCrawler("https://blog.purestorage.com/sitemap_index.xml", pure2rss.CrawlerTimeout(time.Second*30))

	crawler.OnIndexLink(func(l pure2rss.Link) bool {
		return strings.Contains(l.Loc, "post-sitemap")
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	linkCh := make(chan pure2rss.PostLink)

	wg := sync.WaitGroup{}
    client := &http.Client{
        Timeout: time.Second * 5,
    }

	for i := 0; i < 5; i++ {
		wg.Add(1)

		go func(ctx context.Context, linkCh chan pure2rss.PostLink) {
			defer wg.Done()

			for {
				var work pure2rss.PostLink
				select {
				case work = <-linkCh:
				case <-ctx.Done():
					return
				}

                post, err := pure2rss.FetchAndAndParsePost(client, work.Link.Loc)
                if err != nil {
                    slog.Error("parsing post content", slog.Any("err", err), slog.String("url", work.Link.Loc))
                }

				slog.Info("parsed", slog.String("title", post.Title))
			}
		}(ctx, linkCh)
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

		linkCh <- postLink
	})

	go crawler.Run()

	select {
	case err := <-crawler.Done():
		if err != nil {
			panic(err)
		}
	}

	cancel()
	wg.Wait()
}
