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
	result, err := hn.SearchRecent(ctx, &hackernews.SearchRequest{
		Points: "> 500",
	})
	is.NoErr(err)
	is.True(len(result.Stories) >= 10) // 10+ newest stories over 500 points
}

func ExampleClient() {
	ctx := context.Background()
	hn := hackernews.New()
	stories, _ := hn.FrontPage(ctx)
	fmt.Println(len(stories) >= 10)
	// Output: true
}

func TestShowHN(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hn := hackernews.New()
	stories, err := hn.ShowHN(ctx)
	is.NoErr(err)
	is.True(len(stories) >= 10) // 10+ show stories
}

func TestAskHN(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hn := hackernews.New()
	stories, err := hn.AskHN(ctx)
	is.NoErr(err)
	is.True(len(stories) >= 10) // 10+ ask stories
}

func TestNewest(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hn := hackernews.New()
	stories, err := hn.Newest(ctx)
	is.NoErr(err)
	is.True(len(stories) >= 10) // 10+ newest stories
}

func TestFrontPage(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hn := hackernews.New()
	stories, err := hn.FrontPage(ctx)
	is.NoErr(err)
	is.True(len(stories) >= 10) // 10+ front page stories
	for _, story := range stories {
		is.True(story.ID != 0) // story has an ID
	}
}

func TestSecondPage(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	hn := hackernews.New()
	firstPage, err := hn.Search(ctx, &hackernews.SearchRequest{
		Tags: "front_page",
		Page: 1,
	})
	is.NoErr(err)
	is.True(len(firstPage.Stories) >= 10) // 10+ front page stories
	for _, story := range firstPage.Stories {
		is.True(story.ID != 0) // story has an ID
	}
	secondPage, err := hn.Search(ctx, &hackernews.SearchRequest{
		Tags: "front_page",
		Page: 2,
	})
	is.NoErr(err)
	is.True(len(secondPage.Stories) >= 10) // 10+ front page stories
	for _, story := range secondPage.Stories {
		is.True(story.ID != 0) // story has an ID
	}
	is.True(firstPage.Stories[0].ID != secondPage.Stories[0].ID) // first story on first page is not the same as first story on second page
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
