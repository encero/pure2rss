package pure2rss_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/encero/pure2rss/src/pure2rss"
	"github.com/matryer/is"
)

func TestFetchAndParsePost(t *testing.T) {
	is := is.New(t)

	serverURL, cleanup := newTestServer()
	defer cleanup()

	postURL := fmt.Sprintf("%s/post.html", serverURL)

	client := http.DefaultClient
	post, err := pure2rss.FetchAndAndParsePost(client, postURL)
	is.NoErr(err)

	is.Equal(post.Title, "Streamlining Azure VMware Solution: Automating Pure Cloud Block Store Expansion")                                                                           // title
	is.Equal(post.Summary, "A new feature is now available for Pure Cloud Block Store on Azure that enables you to automate capacity upgrades using Azure Functions and PowerShell.") // summary
	is.Equal(post.Tags, []string{"Azure", "Cloud Migration", "Featured", "Pure Cloud Block Store", "VMware"})                                                                         // tags
}

func TestParsePostLink(t *testing.T) {
	is := is.New(t)

	link, err := pure2rss.ParsePostLink(pure2rss.Link{
        Loc: "https://blog.purestorage.com/perspectives/randomware-shakes-up-the-ransomware-game/",
    })
	is.NoErr(err)
	is.Equal(link.Lang, "en")                                       // link lang
	is.Equal(link.Category, "perspectives")                         // link category
	is.Equal(link.Slug, "randomware-shakes-up-the-ransomware-game") // link slug

	link, err = pure2rss.ParsePostLink(pure2rss.Link{
        Loc:"https://blog.purestorage.com/ko/news-events-ko/goat-of-the-year-siriusxm/",
    })
	is.NoErr(err)
	is.Equal(link.Lang, "ko")                        // link lang
	is.Equal(link.Category, "news-events-ko")        //link category
	is.Equal(link.Slug, "goat-of-the-year-siriusxm") // link slug
}
