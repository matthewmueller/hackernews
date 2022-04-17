package hackernews_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/matryer/is"
	"github.com/matthewmueller/hackernews"
)

func TestSearch(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hn := hackernews.New()
	stories, err := hn.SearchRecent(ctx, &hackernews.Search{
		Points: "> 500",
	})
	is.NoErr(err)
	is.Equal(len(stories), 34) // 34 newest stories over 500 points
}

func ExampleClient() {
	ctx := context.Background()
	hn := hackernews.New()
	stories, _ := hn.FrontPage(ctx)
	fmt.Println(len(stories))
	// Output: 34
}

func TestShowHN(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hn := hackernews.New()
	stories, err := hn.ShowHN(ctx)
	is.NoErr(err)
	is.Equal(len(stories), 34) // 34 show stories
}

func TestAskHN(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hn := hackernews.New()
	stories, err := hn.AskHN(ctx)
	is.NoErr(err)
	is.Equal(len(stories), 34) // 34 ask stories
}

func TestNewest(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hn := hackernews.New()
	stories, err := hn.Newest(ctx)
	is.NoErr(err)
	is.Equal(len(stories), 34) // 34 newest stories
}

func TestFrontPage(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hn := hackernews.New()
	stories, err := hn.FrontPage(ctx)
	is.NoErr(err)
	is.Equal(len(stories), 34) // 34 front page stories
}

func TestFind(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hn := hackernews.New()
	post, err := hn.Find(ctx, 1)
	is.NoErr(err)
	is.Equal(post.Title, "Y Combinator") // title is not Y Combinator
	buf, err := json.MarshalIndent(post, "", "  ")
	is.NoErr(err)
	fmt.Println(string(buf))
}
