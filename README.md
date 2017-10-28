
# hermit

[![Platform](http://img.shields.io/badge/platform-macos-blue.svg?style=flat)](https://www.apple.com/macos/how-to-upgrade/)
[![Language](http://img.shields.io/badge/language-go-brightgreen.svg?style=flat)](https://golang.org/)

## Orgasmic URL Clipper That Drives Your Night

<div align="center">
<img src="https://github.com/creaaa/hermit/blob/master/image.png">
</div>

### Caution

Hermit assumes that you set `$GOPATH` in `~/go` as the first path.
If not, this does not work. This situation will be improved soon.
Sorry for the inconvenience.

### Install

As I showed above, make sure you set `$GOPATH` in `~/go` as the first path in advance. 

```sh
$ go get github.com/creaaa/hermit
```

### Usage

```
# add URL (make sure enclose URL in double quote if it includes `?`)
$ hermit add <"URL"> <alias> [description]

# open URL
$ hermit open <ID or alias>... # can designate multiple values by spacing

# shows list of URLs
$ hermit list

# fetch whether URL returns 404, then update database
$ hermit fetch

# delete URL
$ hermit delete <ID or alias>

# delete only URL that is already 404 (need `hermit fetch` in advance to set up a flag)
$ hermit delete -f

# delete all URLs
$ hermit deleteall
```

### Environment

- MacOS v10.13 High Sierra
- Go v1.9.1 or above
- SQLite v3.19.3 or above