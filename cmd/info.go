package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func InfoCommand() *cobra.Command {
	var (
		season      int
		episode     int
		releaseYear int
		tmdbID      int

		isMovie bool
		isSerie bool
	)
	watchCmd := cobra.Command{
		Use:   "info",
		Short: "Get Information about a movie or a tv show",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Infoing")
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

	return &watchCmd
}
