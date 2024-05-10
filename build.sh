# Build binary of the crawler for use in CI
#!/bin/sh

GOOS=linux GOARCH=amd64 go build -o pure2rss .
