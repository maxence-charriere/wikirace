# wikirace
[![Go Report Card](https://goreportcard.com/badge/github.com/maxence-charriere/wikirace)](https://goreportcard.com/report/github.com/maxence-charriere/wikirace)
[![GoDoc](https://godoc.org/github.com/maxence-charriere/wikirace?status.svg)](https://godoc.org/github.com/maxence-charriere/wikirace)

Program that performs a [wikiracing](https://en.wikipedia.org/wiki/Wikiracing).

## Install
```
go get -u github.com/maxence-charriere/wikirace
```

## Build
```
# go in $GOPATH/src/github.com/maxence-charriere/wikirace/bin/wikirace directory
go build
```

## Usage
```
./wikirace -start [TITLE TO START] -end [TITLE TO REACH]
```

## What the code does?
1. Generate a Search object that describe a page to search in from the command line args.
2. Create the different parts to make the program work:
    - An object that keep informations about the jobs (current and processed)
    - A queue to store the Search objects to process.
    - A thread pool that launch search job in paralel.
    - A result handler that analyses job results.
2. Put the Search object into the queue.
3. Searchs are dequeued and processed in a threadpool.
4. For each Search Object dequeued:
    - Page download
    - Page parsing
    - Each link found generate a new Search object which is sent to the result handler.
5. Result handler listen for new Search objects. Then Decides to validate a Search object or 
to put it in the queue for more searching.
