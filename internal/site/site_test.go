package site

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const minimalConfig = `
name = "Ada"
role = "Engineer"
location = "NYC"

[[projects]]
name = "engine"
lang = "rust"
state = "v1.0.0"
status = "release"
what = "Analytical <engine>"
repo = "https://example.com/engine"
link_label = "docs"
link = "https://example.com/docs"

[[projects]]
name = "notes"
lang = "go"
state = "design"
status = "design"

[[contact]]
key = "email"
label = "ada@example.com"
url = "mailto:ada@example.com"
`

// writeContent creates a content directory with site.toml and content/about.md.
func writeContent(t *testing.T, config, about string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "site.toml"), []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "content"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "content", "about.md"), []byte(about), 0o644); err != nil {
		t.Fatal(err)
	}
	return dir
}

func build(t *testing.T, config, about string) string {
	t.Helper()
	out := filepath.Join(t.TempDir(), "public")
	err := Build(Options{
		ContentDir: writeContent(t, config, about),
		OutDir:     out,
		Now:        time.Date(2026, 10, 26, 0, 0, 0, 0, time.UTC),
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestBuildRendersTheHomePage(t *testing.T) {
	out := build(t, minimalConfig, "I build *things*.\n")
	html := readFile(t, filepath.Join(out, "index.html"))

	for _, want := range []string{
		"<title>Ada — lxkota</title>",
		"2026 · NYC",      // year from Options.Now
		"<em>things</em>", // markdown rendered
		`<a href="https://example.com/engine">engine</a>`,  // repo links the name
		`<span class="state release">v1.0.0</span>`,        // status drives the class
		`<a href="https://example.com/docs">docs ↗</a>`,    // optional link column
		`<span class="muted">soon</span>`,                  // no link: placeholder
		"Analytical &lt;engine&gt;",                        // text is escaped
		`<a href="mailto:ada@example.com">ada@example.com`, // contact
	} {
		if !strings.Contains(html, want) {
			t.Errorf("index.html is missing %q", want)
		}
	}
	if strings.Contains(html, "toolbox") {
		t.Error("an empty toolbox should leave out the whole section")
	}
}

func TestBuildShipsStaticFilesAndHeaders(t *testing.T) {
	out := build(t, minimalConfig, "hi\n")
	for _, f := range []string{"style.css", "theme.js", "fonts/JetBrainsMono-Regular.woff2", "fonts/OFL.txt", "_headers"} {
		if _, err := os.Stat(filepath.Join(out, f)); err != nil {
			t.Errorf("missing %s: %v", f, err)
		}
	}
	if !strings.Contains(readFile(t, filepath.Join(out, "_headers")), "default-src 'self'") {
		t.Error("_headers should set a self-only content security policy")
	}
}

func TestMarkdownCannotInjectHTML(t *testing.T) {
	out := build(t, minimalConfig, "hello <script>alert(1)</script>\n")
	if strings.Contains(readFile(t, filepath.Join(out, "index.html")), "<script>alert(1)") {
		t.Error("raw HTML in markdown must not reach the page")
	}
}

func TestRefusesToReplaceAFolderItDidNotCreate(t *testing.T) {
	out := t.TempDir() // exists, has no marker
	keep := filepath.Join(out, "keep.txt")
	if err := os.WriteFile(keep, []byte("mine"), 0o644); err != nil {
		t.Fatal(err)
	}
	err := Build(Options{ContentDir: writeContent(t, minimalConfig, "hi\n"), OutDir: out})
	if err == nil || !strings.Contains(err.Error(), "refusing to replace") {
		t.Fatalf("expected a refusal, got %v", err)
	}
	if _, err := os.Stat(keep); err != nil {
		t.Error("the existing folder must be left untouched")
	}
}

func TestConfigErrors(t *testing.T) {
	tests := map[string]struct{ config, wantErr string }{
		"unknown key": {
			// Top level: appended at the end it would land inside the last [[contact]] table
			config:  "titel = \"typo\"\n" + minimalConfig,
			wantErr: "unknown keys: titel",
		},
		"bad status": {
			config:  strings.Replace(minimalConfig, `status = "design"`, `status = "wip"`, 1),
			wantErr: `status "wip" is not one of`,
		},
		"link without label": {
			config:  strings.Replace(minimalConfig, `link_label = "docs"`, "", 1),
			wantErr: "set both link_label and link",
		},
		"missing name": {
			config:  strings.Replace(minimalConfig, `name = "Ada"`, "", 1),
			wantErr: "name is required",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			dir := writeContent(t, tc.config, "hi\n")
			_, err := LoadConfig(filepath.Join(dir, "site.toml"))
			if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
				t.Errorf("got error %v, want one containing %q", err, tc.wantErr)
			}
		})
	}
}

func TestTheRealSiteBuilds(t *testing.T) {
	// The repo's own site.toml and content must always build
	out := filepath.Join(t.TempDir(), "public")
	if err := Build(Options{ContentDir: filepath.Join("..", ".."), OutDir: out}); err != nil {
		t.Fatal(err)
	}
}
