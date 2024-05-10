# Pure2RSS

Pure2RSS is a web scrapper for [blog.purestorage.com](https://choosealicense.com/licenses/mit/).
It once a day scrapes Pure Storage blog and generate a RSS2.0 feed.

The feed is for `purely-techincal` category of blog posts only and doesn't
include the content of the posts only brief summary.

## Usage

Add this url to your RSS reader.

```
https://raw.githubusercontent.com/encero/pure2rss/main/rss.xml
```

You can run the crawler yourself locally.
```
cd pure2rss

go run .
```

## Contributing

Pull requests are welcome.

Please make sure to update tests as appropriate.

## License

[BSD 3-Clause](https://choosealicense.com/licenses/bsd-3-clause/)
