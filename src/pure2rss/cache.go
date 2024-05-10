package pure2rss

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"sync"
)

var NoDataError = fmt.Errorf("no data found")

type Cache struct {
    mux   *sync.Mutex
	path  string
	posts map[string]Post
}

func NewCache(path string) (*Cache, error) {
    posts := make(map[string]Post)
	fInfo, err := os.Stat(path)
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {

			return nil, fmt.Errorf("cache file path, %w", err)
		}
	} else {
		if fInfo.IsDir() {
			return nil, fmt.Errorf("cache file is directory")
		}

        file, err := os.Open(path)
        if err != nil {
            return nil, fmt.Errorf("opening cache file, %w", err)
        }
        err = json.NewDecoder(file).Decode(&posts)
        if err != nil {
            return nil, fmt.Errorf("decoding cache file, %w", err)
        }
	}

	return &Cache{
        mux: &sync.Mutex{},
		path:  path,
		posts: posts,
	}, nil
}

func (c *Cache) Store(p Post) error {
    c.mux.Lock()
    defer c.mux.Unlock()

	c.posts[p.PostLink.Link.Loc] = p

	return nil
}

func (c *Cache) Load(url string) (Post, error) {
    c.mux.Lock()
    defer c.mux.Unlock()

	post, ok := c.posts[url]
	if !ok {
		return Post{}, NoDataError
	}
	return post, nil
}

func (c *Cache) LoadCategory(category string) ([]Post, error) {
    c.mux.Lock()
    defer c.mux.Unlock()

	posts := []Post{}

	for _, v := range c.posts {
		if v.PostLink.Category == category {
			posts = append(posts, v)
		}
	}

	return posts, nil
}

func (c *Cache) Persist() error {
    c.mux.Lock()
    defer c.mux.Unlock()

	tmpPath := fmt.Sprintf("%s.new", c.path)

	file, err := os.OpenFile(tmpPath, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		return fmt.Errorf("persisting post cache to %q, %w", c.path, err)
	}
	defer file.Close()
	defer os.Remove(tmpPath)

    encoder := json.NewEncoder(file)
    encoder.SetIndent("", " ")

	err = encoder.Encode(c.posts)
	if err != nil {
		return fmt.Errorf("encoding cache to file, %w", err)
	}

	err = os.Rename(tmpPath, c.path)
	if err != nil {
		return fmt.Errorf("renaming encoded file to cache file path, %w", err)
	}

	return nil
}
