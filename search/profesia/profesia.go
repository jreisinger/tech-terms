// https://github.com/jreisinger/profesia-jobs/blob/master/profesia-jobs
package profesia

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/gocolly/colly"
)

type SearchResult struct {
	term	string
	links	[]string
}

func GetJobOffers(term string, ch chan SearchResult) {
	nPages := getNumPages(term)
	chOffers := make(chan []string)

	for n := 1; n <= nPages; n++ {
		url := "https://www.profesia.sk/praca/?search_anywhere=" + term + "&page_num=" + strconv.Itoa(n)
		fmt.Println("Starting search for", url)
		go getJobOffersFromUrl(url, chOffers)
	}

	links := []string{}

	for n := 1; n <= nPages; n++ {
		moreLinks := <-chOffers
		links = append(links, moreLinks...)
	}

	result := SearchResult{
		term: term,
		links: links,
	}
	ch<- result
}

func getJobOffersFromUrl(url string, ch chan []string) {
    c := colly.NewCollector()
	links := []string{}
    c.OnHTML("a", func(e *colly.HTMLElement) {
		if e.Attr("class") == "title" {
			link := e.Attr("href")
			links = append(links, link)
		}
    })

    c.Visit(url)
	ch<- links
}

func getNumPages(term string) int {
    c := colly.NewCollector()

	nPages := 0

    // Find and visit all links
    c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if params := getLinkParams(`page_num=(?P<pagenum>\d+)$`, link); len(params) > 0 {
			pagenum, _ := strconv.Atoi(params["pagenum"])
			if pagenum > nPages {
				nPages = pagenum
			}
		}
    })

	url := "https://www.profesia.sk/search.php?search_anywhere=" + term
    c.Visit(url)

	return nPages
}

func getLinkParams(regEx, url string) (paramsMap map[string]string) {

    var compRegEx = regexp.MustCompile(regEx)
    match := compRegEx.FindStringSubmatch(url)

    paramsMap = make(map[string]string)
    for i, name := range compRegEx.SubexpNames() {
        if i > 0 && i <= len(match) {
            paramsMap[name] = match[i]
        }
    }

    return
}
