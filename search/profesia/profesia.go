// https://github.com/jreisinger/profesia-jobs/blob/master/profesia-jobs
package profesia

import (
	"regexp"
	"strconv"

	"github.com/gocolly/colly"
)

type SearchResult struct {
	term	string
	offers	int
}

func GetNumJobOffers(term string, ch chan SearchResult) {
	result := SearchResult{
		term: term,
		offers: getNumPages(term),
	}
	ch <- result
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
