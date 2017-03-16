# bandcamp-go
A Go program for downloading Bandcamp songs, but moreso a rewrite of [bandcamp-dl](https://github.com/iheanyi/bandcamp-dl).

# Installation
## Linux
* [Go version 1.7](https://github.com/golang/go/releases/tag/go1.7.3)

Clone Git repo:
```
$ go get github.com/PuerkitoBio/goquery/...
$ go get github.com/robertkrimen/otto...
$ go build bandcamp.go
```
TODO (Package in Homebrew)

# Getting Started
```
$ ./bandcamp <album/track_url>
```
it will download the tracks to `music/<album_or_track_title>`.
Will create a binary soon for OS X.

# Running Tests
TODO

# Notes
TODO
