package cmd

import (
	"context"
	"fmt"

	"github.com/kakeetopius/flix/internal/feature"
	"github.com/kakeetopius/flix/internal/tmdb"
	watchproviders "github.com/kakeetopius/flix/internal/watch_providers"
	"github.com/kakeetopius/flix/internal/watch_providers/vidsrc"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

func WatchCommand() *cobra.Command {
	var (
		season      int
		episode     int
		releaseYear int
		tmdbID      int

		isMovie bool
		isSerie bool
		link    bool
		trailer bool
	)
	watchCmd := cobra.Command{
		Use:   "watch",
		Short: "Find and watch a movie or a tv show episode from your browser.",
		Long: `Watch a movie or a tv show episode from your browser.
	
This command automatically finds a link to any movie or tv show episode and opens the link using a broswer on the system for viewing.
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			qtype := feature.FeatureTypeMovie
			if isSerie || episode != 0 || season != 0 {
				qtype = feature.FeatureTypeSerie
			}
			tmdbOpts := tmdb.QueryOptions{
				Query: args[0],
				Year:  releaseYear,
				Type:  qtype,
			}

			if trailer {
				trailerLink, err := tmdb.GetTrailerLink(tmdbOpts)
				if err != nil {
					return err
				}
				browser.OpenURL(trailerLink)
				fmt.Println("Trailer opened in browser")
				fmt.Println("Link: ", trailerLink)
				return nil
			}

			queryOpts := watchproviders.SearchOptions{
				Query:   args[0],
				Year:    releaseYear,
				Type:    qtype,
				Season:  season,
				Episode: episode,
			}

			link, err := vidsrc.NewWatchProvider().Search(context.Background(), queryOpts)
			if err != nil {
				return err
			}

			browser.OpenURL(string(link))
			fmt.Println("Video opened in browser")
			fmt.Println("Link: ", link)

			return nil
		},
	}

	watchCmd.Flags().SortFlags = false
	watchCmd.Flags().IntVarP(&season, "season", "s", 0, "The season if searching for a tv show.")
	watchCmd.Flags().IntVarP(&episode, "episode", "e", 0, "The episode number in a serie's season.")
	watchCmd.Flags().IntVarP(&releaseYear, "year", "y", 0, "The release year of the movie or show")
	watchCmd.Flags().IntVar(&tmdbID, "imdb-id", 0, "Search for a show or movie using the TMDB ID.")
	watchCmd.Flags().BoolVar(&isMovie, "movie", false, "Specifies that the search is for a movie to reduce ambiguity")
	watchCmd.Flags().BoolVar(&isSerie, "serie", false, "Specifies that the search is for a serie to reduce ambiguity")
	watchCmd.Flags().BoolVar(&link, "link", false, "Just returns the link to the video. No browser is opened for viewing.")
	watchCmd.Flags().BoolVarP(&trailer, "trailer", "t", false, "Watch the trailer instead of the actual movie/tv show episode.")

	return &watchCmd
}
