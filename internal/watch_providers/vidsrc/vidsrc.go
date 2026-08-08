// Package vidsrc is used to get watch links from vidsrc.ru
package vidsrc

import (
	"context"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/kakeetopius/flix/internal/feature"
	"github.com/kakeetopius/flix/internal/tmdb"
	watchproviders "github.com/kakeetopius/flix/internal/watch_providers"
)

const (
	baseURL   = "https://vidsrcme.ru/embed"
	moviePath = "/movie"
	seriePath = "/tv"
)

type VidSrc struct{}

type queryParams struct {
	TMDBId  int `url:"tmdb"`
	Season  int `url:"season,omitempty"`
	Episode int `url:"episode,omitempty"`
}

func NewWatchProvider() *VidSrc {
	return new(VidSrc)
}

func (p *VidSrc) Search(ctx context.Context, opts watchproviders.SearchOptions) (watchproviders.Link, error) {
	if opts.Season != 0 || opts.Episode != 0 {
		opts.Type = feature.FeatureTypeSerie
	}

	if opts.TMDBId == 0 {
		tmdbOpts := tmdb.QueryOptions{
			Query: opts.Query,
			Year:  opts.Year,
			Type:  opts.Type,
		}

		id, err := tmdb.GetID(tmdbOpts)
		if err != nil {
			return "", err
		}
		opts.TMDBId = id
	}

	qParams := queryParams{
		TMDBId: opts.TMDBId,
	}
	if opts.Season != 0 {
		qParams.Season = opts.Season
	}
	if opts.Episode != 0 {
		qParams.Episode = opts.Episode
	}

	url, err := url.Parse(baseURL + endPointOf(opts.Type))
	if err != nil {
		return "", err
	}
	values, err := query.Values(qParams)
	if err != nil {
		return "", err
	}
	url.RawQuery = values.Encode()

	return watchproviders.Link(url.String()), nil
}

func endPointOf(fType feature.Type) string {
	switch fType {
	case feature.FeatureTypeMovie:
		return moviePath
	case feature.FeatureTypeSerie:
		return seriePath
	}

	return ""
}
