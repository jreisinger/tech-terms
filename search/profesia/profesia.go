// https://github.com/jreisinger/profesia-jobs/blob/master/profesia-jobs
package profesia

import (
	"log"
	u "net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

var urlBase string = "https://www.profesia.sk/praca/?search_anywhere="

type SearchResult struct {
	Term       string
	Links      []string
	LinksCount int
}

func GetJobOffers(term string, ch chan SearchResult, debug bool) {
	nPages := getNumPages(term)
	chOffers := make(chan []string)
	links := []string{}

	if nPages == 0 { // just one page?
		url := (urlBase + u.QueryEscape(term))
		if debug {
			log.Println("Starting a goroutine to scrape", url)
		}
		go getJobOffersFromUrl(url, chOffers)

		moreLinks := <-chOffers
		links = append(links, moreLinks...)
	}

	if nPages > 0 {
		for n := 1; n <= nPages; n++ {
			url := (urlBase + u.QueryEscape(term) + "&page_num=" + strconv.Itoa(n))
			if debug {
				log.Println("Starting a goroutine to scrape", url)
			}
			go getJobOffersFromUrl(url, chOffers)
		}

		for n := 1; n <= nPages; n++ {
			moreLinks := <-chOffers
			links = append(links, moreLinks...)
		}
	}

	result := SearchResult{
		Term:       term,
		Links:      links,
		LinksCount: len(links),
	}
	ch <- result
}

func getJobOffersFromUrl(url string, ch chan []string) {
	c := colly.NewCollector()
	links := []string{}
	c.OnHTML("a", func(e *colly.HTMLElement) {
		if strings.HasPrefix(e.Attr("id"), "offer") {
			link := e.Attr("href")
			links = append(links, link)
		}
	})

	c.Visit(url)
	ch <- links
}

func getNumPages(term string) int {
	c := colly.NewCollector()

	nPages := 0

	// Find and visit all links
	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if params := getLinkParams(`page_num=(?P<pagenum>\d+)$`, link); len(params) > 0 {
			pagenum, err := strconv.Atoi(params["pagenum"])
			if err != nil {
				panic(err)
			}
			if pagenum > nPages {
				nPages = pagenum
			}
		}
	})

	url := urlBase + u.QueryEscape(term)
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
