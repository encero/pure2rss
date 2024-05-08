package pure2rss_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/encero/pure2rss/src/pure2rss"
	"github.com/matryer/is"
)

func TestCache(t *testing.T) {
	is := is.New(t)
	tempDirPath, err := os.MkdirTemp(os.TempDir(), "pure2rss_")
	is.NoErr(err)
	defer os.RemoveAll(tempDirPath)

	cacheFilePath := fmt.Sprintf("%s/posts.json", tempDirPath)

	cache, err := pure2rss.NewCache(cacheFilePath)
    is.NoErr(err)

	post := pure2rss.Post{
		PostLink: pure2rss.PostLink{
			Link: pure2rss.Link{
				Loc:     "https://blog.purestorage.com/purely-technical/slug/",
				LastMod: time.Date(2024, 1, 2, 3, 4, 5, 6, time.UTC),
			},
			Lang:     "en",
			Category: "purely-technical",
			Slug:     "slug",
		},
		Title:   "post title",
		Summary: "post summary",
		Tags:    []string{"tag1", "tag2"},
	}

    err = cache.Store(post)
    is.NoErr(err) // store post in cache

    loaded, err := cache.Load(post.PostLink.Link.Loc)
    is.NoErr(err) // load existing post from cache
    is.Equal(loaded, post) // loaded post is what we stored

    loaded, err = cache.Load("https://google.com")
    is.Equal(err, pure2rss.NoDataError) // load non existing post from cache

    posts, err :=cache.LoadCategory(post.PostLink.Category)
    is.NoErr(err) // loading post in category
    is.Equal(posts, []pure2rss.Post{post})

    err = cache.Persist()
    is.NoErr(err)

	cache, err = pure2rss.NewCache(cacheFilePath)
    is.NoErr(err)

    loaded, err = cache.Load(post.PostLink.Link.Loc)
    is.NoErr(err) // load existing post from persisted cache
    is.Equal(loaded, post) // loaded post is what we stored and persisted
}
