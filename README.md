
# Orgasm

## Zen Mode Bookmark List That Drives Your Night

### Caution

Orgasm assumes that you set $GOPATH in `~/go` as the first path.
If not, this doesn't work. This situation will be improved soon.
Sorry for the inconvinence. thanks.

### install

make sure you set $GOPATH in `~/go` as the first path in advance. 

```sh
$ go get github.com/creaaa/orgasm
```

### usage

```
# add URL (make sure enclose URL in double quote if it includes `?`)
$ orgasm add <"URL"> <alias> [description]
<br>
# open URL
$ orgasm open <ID or alias>... # can designate multiple values by spacing
<br>
# shows list of URLs
$ orgasm list
<br>
# fetch whether URL returns 404, then update database
$ orgasm fetch
<br>
# delete URL
$ orgasm delete <ID or alias>... # can designate multiple values by spacing
<br>
# delete only URL that is already 404 (need `orgasm fetch` in advance)
$ orgasm delete -f
<br>
# delete all URLs
$ orgasm deleteall
```

### Environment

- MacOS v10.13 High Sierra
- Go v1.9.1 or above
- SQLite v3.19.3 or above