# website
Source of [lxkota.com](https://lxkota.com): a static site built by a small Go generator from
`site.toml` and markdown, styled as a monospace grid with terminal colour themes.

## Editing
- `site.toml`: name, role, projects table, toolbox, contact. Unknown keys fail the build.
- `content/about.md`: the about paragraph (markdown; raw HTML is not rendered).
- `internal/site/assets/`: the page template, CSS, theme script, fonts and `_headers`.

## Building
```
go run ./cmd/site serve          # preview on http://localhost:8080, rebuilt on each page load
go run ./cmd/site build          # write the site to ./public
go test ./...
```
`serve` reads templates and CSS from disk (edit and refresh) and sends the same headers as
production, including the content security policy.

## How it works
- `cmd/site`: the `build` and `serve` commands.
- `internal/site`: loads and validates `site.toml`, renders markdown with
  [goldmark](https://github.com/yuin/goldmark), executes `html/template`, and writes `public/`.
  Assets are embedded in the binary (`//go:embed all:assets`).
- Builds go to a temporary folder that's swapped into place, and an existing output folder is only
  replaced if a previous build created it.
- Everything is self-hosted, including the [JetBrains Mono](https://github.com/JetBrains/JetBrainsMono)
  font (SIL Open Font License), so `_headers` can allow `'self'` only.

## Hosting
Cloudflare Pages, deployed from GitHub Actions. The earlier version, a Go/Echo server on Fly.io, is
tagged [`go-server-final`](https://github.com/Lxkota95/website/tree/go-server-final).

## License
MIT
