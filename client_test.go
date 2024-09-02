package hackernews_test

import (
	"context"
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
	is.True(len(stories) >= 30) // 30+ newest stories over 500 points
}

func ExampleClient() {
	ctx := context.Background()
	hn := hackernews.New()
	stories, _ := hn.FrontPage(ctx)
	fmt.Println(len(stories) >= 30)
	// Output: true
}

func TestShowHN(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hn := hackernews.New()
	stories, err := hn.ShowHN(ctx)
	is.NoErr(err)
	is.True(len(stories) >= 30) // 30+ show stories
}

func TestAskHN(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hn := hackernews.New()
	stories, err := hn.AskHN(ctx)
	is.NoErr(err)
	is.True(len(stories) >= 30) // 30+ ask stories
}

func TestNewest(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hn := hackernews.New()
	stories, err := hn.Newest(ctx)
	is.NoErr(err)
	is.True(len(stories) >= 30) // 30+ newest stories
}

func TestFrontPage(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hn := hackernews.New()
	stories, err := hn.FrontPage(ctx)
	is.NoErr(err)
	is.True(len(stories) >= 30) // 30+ front page stories
	for _, story := range stories {
		is.True(story.ID != 0) // story has an ID
	}
}

func TestFind(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hn := hackernews.New()
	story, err := hn.Find(ctx, 1)
	is.NoErr(err)
	is.Equal(story.ID, 1)
	is.Equal(story.Title, "Y Combinator") // title is not Y Combinator
}
