// Package site builds lxkota.com: site.toml and markdown in, a static folder out.
package site

import (
	"bytes"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/yuin/goldmark"
)

// The templates, CSS, script, fonts and _headers are compiled into the binary, so `site` is one
// self-contained tool. `all:` because a plain `//go:embed assets` skips names starting with "_" or
// ".", which would silently drop _headers.
//
//go:embed all:assets
var embedded embed.FS

// marker is written into every output folder; Build only ever deletes a folder that has it.
const marker = ".site-build"

// Options says where to read content from and where to write the site.
type Options struct {
	ContentDir string    // holds site.toml and content/*.md
	OutDir     string    // replaced by the built site
	Assets     fs.FS     // templates and static files; nil means the embedded copy
	Now        time.Time // for the year in the masthead; zero means time.Now()
}

// page is what the HTML templates see.
type page struct {
	*Config
	About       template.HTML
	ProjectRows []projectRow
	Year        int
}

// projectRow is a project plus its optional description from content/projects/<name>.md; rows with
// a description render as expandable.
type projectRow struct {
	Project
	Description template.HTML
}

// pages are rendered from templates/layout.html plus templates/<page>. Cloudflare Pages serves
// 404.html (with a 404 status) for any path that doesn't exist.
var pages = []string{"index.html", "404.html"}

// Build renders the site into a temporary folder next to OutDir, then swaps it into place, so a
// failed build never leaves a half-written site behind.
func Build(opts Options) error {
	assets := opts.Assets
	if assets == nil {
		sub, err := fs.Sub(embedded, "assets")
		if err != nil {
			return err
		}
		assets = sub
	}
	now := opts.Now
	if now.IsZero() {
		now = time.Now()
	}

	cfg, err := LoadConfig(filepath.Join(opts.ContentDir, "site.toml"))
	if err != nil {
		return err
	}
	about, err := renderMarkdown(filepath.Join(opts.ContentDir, "content", "about.md"))
	if err != nil {
		return err
	}
	rows, err := projectRows(cfg.Projects, filepath.Join(opts.ContentDir, "content", "projects"))
	if err != nil {
		return err
	}
	data := page{Config: cfg, About: about, ProjectRows: rows, Year: now.Year()}

	parent := filepath.Dir(filepath.Clean(opts.OutDir))
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return err
	}
	tmp, err := os.MkdirTemp(parent, ".site-build-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp) // no-op once renamed into place

	for _, name := range pages {
		html, err := renderPage(assets, name, data)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(tmp, name), html, 0o644); err != nil {
			return err
		}
	}
	if err := os.WriteFile(filepath.Join(tmp, marker), nil, 0o644); err != nil {
		return err
	}
	if err := copyTree(assets, "static", tmp); err != nil {
		return err
	}
	headers, err := fs.ReadFile(assets, "_headers")
	if err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(tmp, "_headers"), headers, 0o644); err != nil {
		return err
	}

	if err := replaceDir(tmp, opts.OutDir); err != nil {
		return err
	}
	return nil
}

// renderPage executes the shared layout with one page's blocks. Each page gets its own template set,
// because every page defines the same block names ("title", "main").
func renderPage(assets fs.FS, name string, data page) ([]byte, error) {
	tmpl, err := template.ParseFS(assets, "templates/layout.html", "templates/"+name)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", name, err)
	}
	var out bytes.Buffer
	if err := tmpl.ExecuteTemplate(&out, "layout", data); err != nil {
		return nil, fmt.Errorf("rendering %s: %w", name, err)
	}
	return out.Bytes(), nil
}

// projectRows pairs each project with content/projects/<name>.md when that file exists.
func projectRows(projects []Project, dir string) ([]projectRow, error) {
	rows := make([]projectRow, len(projects))
	for i, p := range projects {
		rows[i].Project = p
		path := filepath.Join(dir, p.Name+".md")
		if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
			continue
		}
		desc, err := renderMarkdown(path)
		if err != nil {
			return nil, err
		}
		rows[i].Description = desc
	}
	return rows, nil
}

// renderMarkdown converts a markdown file to HTML. goldmark's default renderer drops raw HTML in
// the markdown, so content can't inject markup.
func renderMarkdown(path string) (template.HTML, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	var out bytes.Buffer
	if err := goldmark.Convert(src, &out); err != nil {
		return "", fmt.Errorf("%s: %w", path, err)
	}
	return template.HTML(out.String()), nil // #nosec: goldmark escapes text and omits raw HTML
}

// copyTree copies the dir subtree of src into dst, keeping relative paths.
func copyTree(src fs.FS, dir, dst string) error {
	return fs.WalkDir(src, dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := fs.ReadFile(src, path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}

// replaceDir moves the freshly built tmp into place as out. An existing out is only removed when it
// carries the marker, i.e. was made by Build: a mistyped -out can never delete someone's folder.
func replaceDir(tmp, out string) error {
	if _, err := os.Stat(out); err == nil {
		if _, err := os.Stat(filepath.Join(out, marker)); err != nil {
			return fmt.Errorf("refusing to replace %s: it wasn't created by this tool (no %s file)",
				out, marker)
		}
		if err := os.RemoveAll(out); err != nil {
			return err
		}
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	return os.Rename(tmp, out)
}
