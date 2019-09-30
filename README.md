# cf-pr-toggle

A *dead* simple tool in order to toggle my [Page Rules][1] on [Cloudflare][2].

Toggle means either `enable` or `disable` a given rule.

# Usage

You need to pass the 2 following variables to authenticate against the
Cloudflare API:

```bash
$ CLOUDFLARE_EMAIL=<youremail> CLOUDFLARE_TOKEN=<yourtoken> ./cf-pr-toggle
```

[1]: https://www.cloudflare.com/features-page-rules
[2]: https://api.cloudflare.com/#page-rules-for-a-zone-edit-page-rule
