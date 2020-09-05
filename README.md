# tech-terms

## About

Search (technical) terms in job sites. Calculate number of job offers per term, store and plot (graph) the results. Job sites currently covered:

* https://profesia.sk -> https://github.com/jreisinger/profesia-jobs-per-term

## Usage

```
go install
tech-terms search perl python ruby
```

## Roadmap

* [x] search for many terms quickly (goroutines)
* [x] persist search results without duplicates (sqlite)
* [x] find place where to store the DB ([github](https://github.com/jreisinger/profesia-jobs-per-term), ~~local disk + backups~~)
* [x] show graphs on the Internet -> https://jreisinger.github.io/profesia-jobs-per-term/
* [x] setup regular graphs updates
* [ ] add more job sites
