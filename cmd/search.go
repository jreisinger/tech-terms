package cmd

import (
	"fmt"
	"log"

	"github.com/jreisinger/tech-terms/search/profesia"
	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search the specified terms",
	Run: func(cmd *cobra.Command, searchTerms []string) {
        ch := make(chan profesia.SearchResult)

        for _, term := range searchTerms {
			log.Println("Starting a goroutine to search for", term)
            go profesia.GetJobOffers(term, ch)
        }

        for range searchTerms {
            //fmt.Printf("%v\n", <-ch)
			result := <-ch
			fmt.Println(result.Term, result.LinksCount)
        }
	},
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// searchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// searchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
