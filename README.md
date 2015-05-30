Description
===========

TwittBid Search is the search client for Twitter platform.

Requirements
============

To use this, you must have the following items installed and working:

* [Go](https://golang.org/)
* [GoDep](https://github.com/tools/godep)

Usage
=====

To start, create a file twittbid/credentials.go with the contents:
```
package twittbid

var key = "yourKey"
var secret = "yourSecret"
var accessTokens = "yourToken
var accessTokenSecret = "YourTokenSecret"

```


Test
=====

```
$ go test
Interesting Take on generic-ish functionality in #golang https://t.co/PDmCeZ5sBDRT @Tutor
(...)

#golangRT @enneff: Go 1.5 will likely set GOMAXPROCS=runtime.NumCPU()
(...)
ok  	twittbid	0.651s

```

RUN
====

```
$ export PORT=3006 && go run main.go
[martini] listening on :3006 (development)
```

SEARCH
======

```
$ curl "http://localhost:3006/search/golang"
```

License
=======
The MIT License (MIT)
Copyright (c) 2015 Thomas Modeneis <thomas.modeneis@gmail.com>
