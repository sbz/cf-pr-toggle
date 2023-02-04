# cf-pr-toggle

[![Build Status](https://api.travis-ci.org/sbz/cf-pr-toggle.svg?branch=master)](https://travis-ci.org/sbz/cf-pr-toggle)

A *dead* simple [Golang][0] project tool in order to toggle my [Page Rules][1] on [Cloudflare][2].

Toggle means either `enable` or `disable` a given rule.

# Go dependencies

- [cloudflare/cloudflare-go](https://github.com/cloudflare/cloudflare-go)
- [setvargo/go-envconfig](https://github.com/sethvargo/go-envconfig)
- [smartystreets/goconvey](https://github.com/smartystreets/goconvey)

# Building

```bash
go build
```

# Testing

```bash
go test -v
```

# Usage

You need to pass the 2 following variables to authenticate against the
Cloudflare API:

```bash
$ export CLOUDFLARE_EMAIL=<youremail> CLOUDFLARE_TOKEN=<yourtoken>
$ ./cf-pr-toggle # list existing rules
$ ./cf-pr-toggle <rule-id> # toggle rule with id <rule-id>
```
[0]: https://go.dev/
[1]: https://www.cloudflare.com/features-page-rules
[2]: https://api.cloudflare.com/#page-rules-for-a-zone-edit-page-rule
