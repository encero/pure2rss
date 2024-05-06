package pure2rss_test

import (
	"embed"
	"testing"
	"time"

	"github.com/matryer/is"

	"github.com/encero/pure2rss/src/pure2rss"
)

//go:embed test_data/*
var sitemaps embed.FS

func TestParseSitemapList(t *testing.T) {
	is := is.New(t)

	file, err := sitemaps.Open("test_data/sitemap_index.xml")
	is.NoErr(err)
	defer file.Close()

	list, err := pure2rss.ParseSiteMapList(file)
	is.NoErr(err)

	is.Equal(len(list), 6) // parsed sitemap has two items

	first := list[0]
	second := list[1]

	is.Equal(first.Loc, "#server#/post-sitemap.xml")
	is.Equal(second.Loc, "#server#/post-sitemap2.xml")
}


func TestParseSitemap(t *testing.T) {
	is := is.New(t)

	file, err := sitemaps.Open("test_data/post-sitemap.xml")
	is.NoErr(err)
	defer file.Close()

	posts, err := pure2rss.ParseSiteMap(file)
	is.NoErr(err)

	is.Equal(len(posts), 1000) // expected number of parsed urls

	first := posts[0]
	second := posts[1]

	is.Equal(first.Loc, "https://blog.purestorage.com/")
	is.True(first.LastMod.Sub(time.Date(2024, 5, 2, 2, 0, 5, 0, time.UTC)).Seconds() == 0) // expected and parsed last modification are equal

	is.Equal(second.Loc, "https://blog.purestorage.com/ko/perspectives-ko/dietz-on-the-day-delivering-the-data-platform-for-the-cloud-era/")
	is.True(second.LastMod.Sub(time.Date(2017, 9, 29, 3, 45, 7, 0, time.UTC)).Seconds() == 0) // expected and parsed last modification are equal
}
