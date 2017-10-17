
# Orgasm

[![Platform](http://img.shields.io/badge/platform-macos-blue.svg?style=flat)](https://www.apple.com/macos/how-to-upgrade/)
[![Language](http://img.shields.io/badge/language-go-brightgreen.svg?style=flat)](https://golang.org/)

## Zen Mode Bookmark List That Drives Your Night

<img src="https://github.com/creaaa/orgasm/blob/master/image.png">

### Caution

Orgasm assumes that you set `$GOPATH` in `~/go` as the first path.
If not, this does not work. This situation will be improved soon.
Sorry for the inconvenience. thanks.

### install

As I showed above, make sure you set `$GOPATH` in `~/go` as the first path in advance. 

```sh
$ go get github.com/creaaa/orgasm
```

### usage

```
# add URL (make sure enclose URL in double quote if it includes `?`)
$ orgasm add <"URL"> <alias> [description]

# open URL
$ orgasm open <ID or alias>... # can designate multiple values by spacing

# shows list of URLs
$ orgasm list

# fetch whether URL returns 404, then update database
$ orgasm fetch

# delete URL
$ orgasm delete <ID or alias>

# delete only URL that is already 404 (need `orgasm fetch` in advance)
$ orgasm delete -f

# delete all URLs
$ orgasm deleteall
```

### Environment

- MacOS v10.13 High Sierra
- Go v1.9.1 or above
- SQLite v3.19.3 or above