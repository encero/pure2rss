package pure2rss

import (
	"encoding/xml"
	"fmt"
	"io"
)

type SiteMapLink struct {
	Loc string
}

func ParseSiteMapList(data io.Reader) ([]SiteMapLink, error) {
	decoder := xml.NewDecoder(data)

	var sitemapList = struct {
		XMLName xml.Name `xml:"sitemapindex"`
		SiteMap []struct {
			Loc string `xml:"loc"`
		} `xml:"sitemap"`
	}{}

	err := decoder.Decode(&sitemapList)

	if err != nil {
		return nil, fmt.Errorf("decoding sitemap list from xml: %w", err)
	}

	output := make([]SiteMapLink, 0, len(sitemapList.SiteMap))

	for _, sitemap := range sitemapList.SiteMap {
		output = append(output, SiteMapLink{Loc: sitemap.Loc})
	}

	return output, nil
}

type Link struct{
    Loc string
    LastMod string
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
		output = append(output, Link{Loc: url.Loc, LastMod: url.LastMod})
	}

	return output, nil
}
