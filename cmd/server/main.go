package main

import (
	"html/template"
	"io"
	"os"
	"path/filepath"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func NewTemplate() *Template {
	return &Template{
		templates: template.Must(template.ParseFiles("./web/template.html")),
	}
}

// PageData represents the data that will be passed to the template
type PageData struct {
	Title   string
	Content template.HTML
}

func NewPage(file string, e *echo.Echo) PageData {
	html := markdownToHTML(file, e)
	e.Logger.Debug(string(html))

	return PageData{
		Title:   filepath.Base(file),
		Content: html,
	}
}

func markdownToHTML(file string, e *echo.Echo) template.HTML {
	file_content, err := os.ReadFile(file)
	if err != nil {
		e.Logger.Error("Failed to read file: %s", file)
	}

	// Create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs
	p := parser.NewWithExtensions(extensions)

	// Parse the markdown document
	doc := p.Parse(file_content)

	// Create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	// Render the HTML
	htmlContent := markdown.Render(doc, renderer)

	// Return as template.HTML to prevent auto-escaping in template
	return template.HTML(htmlContent)
}

func main() {
	// initialize server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)
	e.Renderer = NewTemplate()

	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		// create page data with a title based on the file name
		page_data := NewPage("./web/content/index.md", e)

		// Render the template with our data
		return c.Render(200, "template.html", page_data)
	})

	// start server
	e.Logger.Info(e.Start(":42069"))
}
