This site. A small Go program turns `site.toml` and markdown into static HTML, deployed to Cloudflare
Pages: every pull request gets a preview, and merging publishes it.

- Go standard library (`html/template`, `embed`) plus goldmark for markdown
- No external requests at all, fonts included, so the content security policy only allows the site
  itself
- Previously a Go web server on Fly.io; it never needed to be one
