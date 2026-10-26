// Command site builds lxkota.com from site.toml and content/*.md.
//
//	site build [-content DIR] [-out DIR]     write the static site (default: ./public)
//	site serve [-content DIR] [-addr ADDR]   preview locally, rebuilding on every page load
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Lxkota95/website/internal/site"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		usage()
	}
	switch os.Args[1] {
	case "build":
		build(os.Args[2:])
	case "serve":
		serve(os.Args[2:])
	default:
		usage()
	}
}

func usage() {
	fmt.Fprintln(os.Stderr, "usage: site build [-content DIR] [-out DIR]\n       site serve [-content DIR] [-addr ADDR]")
	os.Exit(2)
}

func build(args []string) {
	flags := flag.NewFlagSet("build", flag.ExitOnError)
	content := flags.String("content", ".", "directory with site.toml and content/")
	out := flags.String("out", "public", "output directory (replaced on each build)")
	flags.Parse(args)

	if err := site.Build(site.Options{ContentDir: *content, OutDir: *out}); err != nil {
		log.Fatalf("build failed: %v", err)
	}
	log.Printf("built %s", *out)
}

// serve rebuilds the whole site (it takes milliseconds) before each page request, so edits to
// site.toml and markdown show on refresh. When run from the repo, templates and CSS are read from
// disk too instead of the copy embedded at compile time.
func serve(args []string) {
	flags := flag.NewFlagSet("serve", flag.ExitOnError)
	content := flags.String("content", ".", "directory with site.toml and content/")
	addr := flags.String("addr", "localhost:8080", "listen address")
	flags.Parse(args)

	var assets fs.FS
	if dir := filepath.Join(*content, "internal", "site", "assets"); isDir(dir) {
		assets = os.DirFS(dir)
		log.Printf("using templates and static files from %s", dir)
	}

	tmp, err := os.MkdirTemp("", "site-serve-*")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmp)
	out := filepath.Join(tmp, "public")

	var mu sync.Mutex
	files := http.FileServer(http.Dir(out))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || filepath.Ext(r.URL.Path) == ".html" {
			mu.Lock()
			err := site.Build(site.Options{ContentDir: *content, OutDir: out, Assets: assets})
			mu.Unlock()
			if err != nil {
				http.Error(w, "build failed:\n\n"+err.Error(), http.StatusInternalServerError)
				log.Printf("build failed: %v", err)
				return
			}
		}
		// Send the same headers Cloudflare Pages will, so CSP problems show up locally
		for name, value := range rootHeaders(filepath.Join(out, "_headers")) {
			w.Header().Set(name, value)
		}
		files.ServeHTTP(w, r)
	})

	log.Printf("serving on http://%s (Ctrl-C to stop)", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// rootHeaders reads the "/*" block of a Cloudflare Pages _headers file: indented "Name: value"
// lines under the "/*" line. Other path blocks only tune caching, so the preview ignores them.
func rootHeaders(path string) map[string]string {
	headers := map[string]string{}
	data, err := os.ReadFile(path)
	if err != nil {
		return headers
	}
	inRoot := false
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "" || strings.HasPrefix(trimmed, "#"):
			continue
		case !strings.HasPrefix(line, " ") && !strings.HasPrefix(line, "\t"):
			inRoot = trimmed == "/*"
		case inRoot:
			if name, value, ok := strings.Cut(trimmed, ":"); ok {
				headers[strings.TrimSpace(name)] = strings.TrimSpace(value)
			}
		}
	}
	return headers
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
