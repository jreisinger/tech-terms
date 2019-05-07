package cmd

import (
	"bufio"
	"encoding/csv"
	_ "fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"

	"github.com/spf13/cobra"
)

type Record struct {
	Date  string
	Term  string
	Count int
}

// plotCmd represents the plot command
var plotCmd = &cobra.Command{
	Use:   "plot",
	Short: "Generate a graph containing the specified terms",
	Run: func(cmd *cobra.Command, args []string) {
		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}

		p.Title.Text = "Scripting languages"
		//p.X.Label.Text = "X"
		// define how we convert and display time.Time values.
		p.X.Tick.Marker = plot.TimeTicks{Format: "2006-01-02\n15:04"}
		p.Y.Label.Text = "Job offers"

		err = plotutil.AddLines(p,
			"bash", jobsPoints("bash"),
			"perl", jobsPoints("perl"),
			"python", jobsPoints("python"),
			"ruby", jobsPoints("ruby"),
			"shell", jobsPoints("shell"),
		)
		if err != nil {
			log.Fatal(err)
		}

		// Save the plot to a PNG file.
		if err := p.Save(40*vg.Centimeter, 20*vg.Centimeter, "points.png"); err != nil {
			log.Fatal(err)
		}
	},
}

// Read in CSV records (lines) from a file
func readCSV(csvFileName string) []Record {
	var records []Record

	csvFile, err := os.Open(csvFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	r := csv.NewReader(bufio.NewReader(csvFile))
	r.Comma = ';'
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("%T: %v\n", record, record)
		count, err := strconv.Atoi(record[2])
		if err != nil {
			log.Fatal(err)
		}
		records = append(records, Record{
			Date:  record[0],
			Term:  record[1],
			Count: count,
		})
	}

	return records
}

// Generate points to be graphed for a given term
func jobsPoints(term string) plotter.XYs {
	records := readCSV("jobs-count.csv")

	var newRecords []Record

	for _, record := range records {
		if !(record.Term == term) {
			continue
		}
		newRecords = append(newRecords, record)
	}

	pts := make(plotter.XYs, len(newRecords))
	for i, record := range newRecords {
		date, err := time.Parse("2006-01-02", record.Date)
		if err != nil {
			log.Fatal(err)
		}
		pts[i].X = float64(date.Unix())
		pts[i].Y = float64(record.Count)
	}

	return pts
}

func init() {
	rootCmd.AddCommand(plotCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// plotCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// plotCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
