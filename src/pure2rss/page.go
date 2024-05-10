package pure2rss

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Post struct {
	PostLink    PostLink
	Title       string
	Summary     string
	PublishDate time.Time
	Tags        []string
}

func FetchAndAndParsePost(c *http.Client, postLink PostLink) (Post, error) {
	req, err := http.NewRequest(http.MethodGet, postLink.Link.Loc, nil)
	if err != nil {
		return Post{}, fmt.Errorf("constructing post fetch request, %w", err)
	}

	response, err := c.Do(req)
	if err != nil {
		return Post{}, fmt.Errorf("fetching post, %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return Post{}, fmt.Errorf("unexpected http status when fetching post expected 200 got %d", response.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return Post{}, fmt.Errorf("parsing post page with goquery, %w", err)
	}

	post := Post{
		PostLink: postLink,
	}

	h1 := doc.Find("main h1")
	post.Title = h1.Text()

	summary := doc.Find("section.wpa-content-summary p")
	post.Summary = summary.Text()

	publishDateNode := doc.Find(".post-meta__date")
	post.PublishDate, err = time.Parse("Jan 02, 2006", publishDateNode.Text())
	if err != nil {
		return Post{}, fmt.Errorf("parsing post date, %w", err)
	}

	post.Tags = []string{}

	doc.Find("div.post-header ul.post-tags li a").Each(func(i int, s *goquery.Selection) {
		post.Tags = append(post.Tags, s.Text())
	})

	return post, nil
}

type PostLink struct {
	Link     Link
	Lang     string
	Category string
	Slug     string
}

func ParsePostLink(link Link) (PostLink, error) {
	parsed, err := url.Parse(link.Loc)
	if err != nil {
		return PostLink{}, fmt.Errorf("parsing post url, %w", err)
	}

	path := parsed.Path
	path = strings.Trim(path, "/")

	parts := strings.Split(path, "/")

	if len(parts) == 2 {
		return PostLink{
			Link:     link,
			Lang:     "en",
			Category: parts[0],
			Slug:     parts[1],
		}, nil
	}

	if len(parts) == 3 {
		return PostLink{
			Link:     link,
			Lang:     parts[0],
			Category: parts[1],
			Slug:     parts[2],
		}, nil
	}

	return PostLink{}, fmt.Errorf("the post url path has unexpected number of parts")
}
