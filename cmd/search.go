package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jreisinger/tech-terms/search/profesia"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

var store bool

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search (and store) the specified terms",
	Run: func(cmd *cobra.Command, searchTerms []string) {
		ch := make(chan profesia.SearchResult)

		for _, term := range searchTerms {
			if Debug {
				log.Println("Starting a goroutine to search for", term)
			}
			time.Sleep(time.Millisecond * 500)
			go profesia.GetJobOffers(term, ch, Debug)
		}

		for range searchTerms {
			//fmt.Printf("%v\n", <-ch)
			result := <-ch
			fmt.Println(result.Term, result.LinksCount)

			if store {
				storeToSqlite(result.Term, result.LinksCount)
			}
		}
	},
}

func storeToSqlite(term string, count int) {
	database, err := sql.Open("sqlite3", "data/jobs-count.db")
	if err != nil {
		log.Fatalf("opening sqlite3 file data/jobs-count.db: %v", err)
	}
	if err := database.Ping(); err != nil {
		log.Fatalf("pinging sqlite3 DB in data/jobs-count.db: %v", err)
	}

	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS jobs (date TEXT, term TEXT, count INTEGER, PRIMARY KEY (date, term))")
	if err != nil {
		log.Fatalf("preparing DB statement: %v", err)
	}
	statement.Exec()

	statement, err = database.Prepare("INSERT INTO jobs (date, term, count) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatalf("preparing DB statement: %v", err)
	}
	dt := time.Now()
	statement.Exec(dt.Format("2006-01-02"), strings.ToLower(term), count)
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
	searchCmd.Flags().BoolVarP(&store, "store", "s", false, "store results into database")
}
