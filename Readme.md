# hackernews

[![GoDoc](https://godoc.org/github.com/matthewmueller/hackernews?status.svg)](https://godoc.org/github.com/matthewmueller/hackernews)

Simple Hacker News client for Go. Currently only supports reading from Hacker News.

## Example

### List stories from the front page

```go
hn := hackernews.New()
ctx := context.Background()
stories, err := hn.FrontPage(ctx)
```

### Get the newest stories

```go
hn := hackernews.New()
ctx := context.Background()
stories, err := hn.Newest(ctx)
```

### Find a story and its comments

```go
hn := hackernews.New()
ctx := context.Background()
story, err := hn.Find(ctx, 1)
```

## Install

```sh
go get -u github.com/matthewmueller/gotext
```

## Authors

- Matt Mueller [@mattmueller](https://twitter.com/mattmueller)

## License

MIT
