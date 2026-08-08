// Package tmdb provides functions to search for details of a feature from themoviedb.org
// A feature is a movie or a tv show on themoviedb.org,
package tmdb

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kakeetopius/flix/internal/feature"
	"github.com/kakeetopius/net-tools/pkg/httpclient"
)

const (
	baseURL         = "https://www.themoviedb.org"
	searchURL       = baseURL + "/search"
	movieSearchPath = "/movie"
	serieSearchPath = "/tv"
)

// searchPath returns the `themoviedb.org` search path for the query type.
func searchPathOf(fType feature.Type) string {
	switch fType {
	case feature.FeatureTypeMovie:
		return movieSearchPath
	case feature.FeatureTypeSerie:
		return serieSearchPath
	default:
		return ""
	}
}

type QueryOptions struct {
	Query string       `url:"query"`
	Type  feature.Type `url:"-"`
	Year  int          `url:"-"`
}

// GetID accepts a query and returns the tmdb id of the first result of the search. The query can be a movie or a tv show.
func GetID(query QueryOptions) (int, error) {
	resultsPage, err := getSearchResultsPage(query)
	if err != nil {
		return 0, err
	}

	result, err := getFirstResult(resultsPage, query.Type)
	if err != nil {
		return 0, err
	}
	link, err := getFeaturePagePathFromResult(result)
	if err != nil {
		return 0, err
	}

	return getIDFromFeaturePath(link)
}

// GetTrailerLink accepts a query and returns the link to the trailer of the first result of the search. The query can be a movie or a tv show.
func GetTrailerLink(query QueryOptions) (string, error) {
	featurePage, err := getFeaturePage(query)
	if err != nil {
		return "", err
	}

	return getTrailerLinkFromFeaturePage(featurePage)
}

// getSearchResultsPage accepts a query and returns the search results page from themoviedb.org of the query. The query can be a movie or a tv show.
func getSearchResultsPage(query QueryOptions) (*goquery.Document, error) {
	client := httpclient.New().WithBaseURL(searchURL)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if query.Year != 0 {
		query.Query = fmt.Sprintf("%v y:%v", query.Query, query.Year)
	}

	resp, err := client.Get(ctx, searchPathOf(query.Type), query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return goquery.NewDocumentFromReader(resp.Body)
}

// getFeaturePage accepts a query and returns the feature page of the first result of the search. The query can be a movie or a tv show.
// a feature is a movie or a tv show on themoviedb.org
func getFeaturePage(query QueryOptions) (*goquery.Document, error) {
	id, err := GetID(query)
	if err != nil {
		return nil, err
	}

	client := httpclient.New().WithBaseURL(baseURL)
	resp, err := client.Get(context.Background(), getFeaturePagePath(id, query.Type), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return goquery.NewDocumentFromReader(resp.Body)
}

// getTrailerLinkFromFeaturePage accepts a feature page and returns the link to the trailer of the feature. The feature can be a movie or a tv show.
func getTrailerLinkFromFeaturePage(doc *goquery.Document) (string, error) {
	allLinks := doc.Find("a")
	var trailerID *string
	var trailerSite *string

	allLinks.Each(func(i int, s *goquery.Selection) {
		title, _ := s.Attr("data-title")
		if title != "Play Trailer" {
			return
		}
		id, found := s.Attr("data-id")
		if !found {
			return
		}
		trailerID = &id

		site, found := s.Attr("data-site")
		if !found {
			return
		}
		trailerSite = &site
	})

	if trailerID == nil || trailerSite == nil {
		return "", fmt.Errorf("could not get link for the trailer")
	}

	return getWatchLinkOfSite(*trailerSite, *trailerID)
}

// getWatchLinkOfSite accepts a site and an id and returns the link to watch the trailer on the site.
func getWatchLinkOfSite(site string, id string) (string, error) {
	switch site {
	case "YouTube":
		return "https://youtube.com/watch?v=" + id, nil
	default:
		return "", fmt.Errorf("error getting trailer link: unknown site %v", site)
	}
}

// getFirstResult accepts a search results page from themoviedb.org and returns the first result of the search. The query can be a movie or a tv show.
func getFirstResult(doc *goquery.Document, qtype feature.Type) (*goquery.Selection, error) {
	divID := "#movie_results"
	if qtype == feature.FeatureTypeSerie {
		divID = "#tv_results"
	}

	return doc.Find(divID + " > div > div").Children().First(), nil
}

// getFeaturePagePathFromResult accepts a single search result and returns the path to the feature page of the result.
func getFeaturePagePathFromResult(result *goquery.Selection) (string, error) {
	links := result.Find("a")
	var link *string

	links.Each(func(i int, s *goquery.Selection) {
		_, found := s.Attr("data-media-type")
		if !found {
			return
		}
		href, _ := s.Attr("href")
		link = &href
	})

	if link == nil {
		return "", fmt.Errorf("could not get the feature path")
	}

	return *link, nil
}

// getFeaturePagePath accepts a tmdb id and a query type and returns the path to the feature page of the feature on themoviedb.org. The feature can be a movie or a tv show.
func getFeaturePagePath(tmdbID int, qtype feature.Type) string {
	return fmt.Sprintf("%v/%v", searchPathOf(qtype), tmdbID)
}

// getIDFromFeaturePath accepts a feature path and returns the tmdb id of the feature. The feature can be a movie or a tv show.
func getIDFromFeaturePath(link string) (int, error) {
	parts := strings.Split(link, "/")
	if len(parts) != 3 {
		return 0, fmt.Errorf("could not get the tmdb id for the query")
	}

	tmdbID := parts[2]
	dashIndex := strings.Index(tmdbID, "-")

	if dashIndex != -1 {
		tmdbID = tmdbID[:dashIndex]
	}

	return strconv.Atoi(tmdbID)
}
