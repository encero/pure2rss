package pure2rss_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/encero/pure2rss/src/pure2rss"
	"github.com/matryer/is"
)

func TestCrawler(t *testing.T) {
	is := is.NewRelaxed(t)

    serverURL, cleanup := newTestServer()
    defer cleanup()

	sitemapURL := fmt.Sprintf("%s/sitemap_index.xml", serverURL)

	crawler := pure2rss.NewCrawler(sitemapURL)

    onIndexCalled := 0
	crawler.OnIndex(func(links []pure2rss.Link) []pure2rss.Link {
        onIndexCalled += 1

        is.Equal(len(links), 6) // index map url count
		return links
	})

	go crawler.Run()

	select {
	case <-crawler.Done():
	case <-time.NewTimer(time.Millisecond * 100).C:
		is.Fail() // crawler didn't fininsh in time
	}

    is.Equal(onIndexCalled, 1) // OnIndex was called once
}

func newTestServer() (string, func()) {
    s := httptest.NewServer(http.HandlerFunc(func (w http.ResponseWriter, r *http.Request){
        filePath := fmt.Sprintf("sitemaps/%s", r.URL.Path)

        file, err := sitemaps.Open(filePath)
        if err != nil {
            fmt.Println("Can't open file", filePath)
            w.WriteHeader(http.StatusNotFound)
            return
        }
        defer file.Close()

        w.WriteHeader(http.StatusOK)
        io.Copy(w, file)
    }))

    return s.URL, func () {
        s.Close()
    }
}
