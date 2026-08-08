// Package watchproviders is used to interact with different providers that provide watch links for movies/tv shows.
package watchproviders

import (
	"context"

	"github.com/kakeetopius/flix/internal/feature"
)

type WatchProvider interface {
	Search(ctx context.Context, opts SearchOptions) (Link, error)
}

type SearchOptions struct {
	Query    string
	Type     feature.Type
	IMDBId   int
	TMDBId   int
	Season   int
	Episode  int
	Year     int
	Language string
}

type Link string
