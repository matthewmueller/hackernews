# hackernews

[![GoDoc](https://godoc.org/github.com/matthewmueller/hackernews?status.svg)](https://godoc.org/github.com/matthewmueller/hackernews)

Simple Hacker News client for Go. Currently only supports reading from Hacker News.

```go
hn := hackernews.New()
ctx := context.Background()
stories, err := hn.FrontPage(ctx)
```

## Install

```sh
go get -u github.com/matthewmueller/gotext
```

## Authors

- Matt Mueller [@mattmueller](https://twitter.com/mattmueller)

## License

MIT
