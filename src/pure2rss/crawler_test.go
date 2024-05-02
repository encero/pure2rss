package pure2rss_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
	crawler.OnIndexLink(func(link pure2rss.Link) bool {
		onIndexCalled += 1

		return filterOutNonPostSiteMaps(link)
	})

	onPostCalled := 0
	crawler.OnPostLink(func(l pure2rss.Link) bool {
		onPostCalled += 1

		return true
	})

	go crawler.Run()

	select {
	case err := <-crawler.Done():
		is.NoErr(err) // crawler error
	case <-time.NewTimer(time.Second).C:
		is.Fail() // crawler didn't fininsh in time
	}

	is.Equal(onIndexCalled, 6)   // OnIndex was called expected times
	is.Equal(onPostCalled, 2438) // OnPost was called expected times
}

func newTestServer() (string, func()) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		filePath := fmt.Sprintf("sitemaps%s", r.URL.Path)

		file, err := sitemaps.Open(filePath)
		if err != nil {
			fmt.Println("Can't open file", filePath)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		defer file.Close()


        data, err := io.ReadAll(file)
        if err != nil {
			fmt.Println("Can't read the file", filePath, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
        }


        serverURL := fmt.Sprintf("http://%s", r.Host)

        fmt.Println(serverURL)

        data = []byte(strings.ReplaceAll(string(data), "#server#", serverURL))

		w.WriteHeader(http.StatusOK)
        w.Write(data)
	}))

	return s.URL, func() {
		s.Close()
	}
}

func filterOutNonPostSiteMaps(link pure2rss.Link) bool {
    return strings.Contains(link.Loc, "post-")
}
