package pure2rss

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

func ParseSiteMapList(data io.Reader) ([]Link, error) {
	decoder := xml.NewDecoder(data)

	var sitemapList = struct {
		XMLName xml.Name `xml:"sitemapindex"`
		SiteMap []struct {
			Loc string `xml:"loc"`
            LastMod string `xml:"lastmod"`
		} `xml:"sitemap"`
	}{}

	err := decoder.Decode(&sitemapList)

	if err != nil {
		return nil, fmt.Errorf("decoding sitemap list from xml: %w", err)
	}

    output := make([]Link, 0, len(sitemapList.SiteMap))

	for _, sitemap := range sitemapList.SiteMap {
        lastMod, err := time.Parse(time.RFC3339, sitemap.LastMod)
        if err != nil {
            return nil, fmt.Errorf("parsing last modification time of a page, %w", err)
        }
		output = append(output, Link{Loc: sitemap.Loc})
	}

	return output, nil
}

type Link struct{
    Loc string
    LastMod time.Time
}

func ParseSiteMap(reader io.Reader) ([]Link, error) {
	decoder := xml.NewDecoder(reader)

	var sitemap = struct {
		XMLName xml.Name `xml:"urlset"`
		Urls    []struct {
			Loc     string `xml:"loc"`
			LastMod string `xml:"lastmod"`
		} `xml:"url"`
	}{}

	err := decoder.Decode(&sitemap)

	if err != nil {
		return nil, fmt.Errorf("decoding sitemap list from xml: %w", err)
	}

	output := make([]Link, 0, len(sitemap.Urls))

	for _, url := range sitemap.Urls {
        lastMod, err := time.Parse(time.RFC3339, url.LastMod)
        if err != nil {
            return nil, fmt.Errorf("parsing last modification time of a page, %w", err)
        }

		output = append(output, Link{Loc: url.Loc, LastMod: lastMod})
	}

	return output, nil
}
