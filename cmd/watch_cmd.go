package cmd

import (
	"fmt"

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
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Searching")
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
