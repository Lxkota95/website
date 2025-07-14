package main

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	Title       string
	Content     template.HTML
	Description string
	Author      string
}

func NewPage(file string, e *echo.Echo) PageData {
	html := markdownToHTML(file, e)
	e.Logger.Debug(string(html))

	// Extract title from filename or content
	title := extractTitle(file)

	return PageData{
		Title:       title,
		Content:     html,
		Description: "Jack Coleman - Staff Platform Engineer in FinTech, NYC",
		Author:      "Jack Coleman",
	}
}

func extractTitle(filename string) string {
	base := filepath.Base(filename)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	// Convert filename to title case
	if name == "index" {
		return "Jack Coleman - Software Engineer"
	}

	return strings.Title(strings.ReplaceAll(name, "_", " "))
}

func markdownToHTML(file string, e *echo.Echo) template.HTML {
	file_content, err := os.ReadFile(file)
	if err != nil {
		e.Logger.Error("Failed to read file: %s", file)
		return template.HTML("<p>Error loading content</p>")
	}

	// Create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.Tables | parser.Strikethrough
	p := parser.NewWithExtensions(extensions)

	// Parse the markdown document
	doc := p.Parse(file_content)

	// Create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{
		Flags: htmlFlags,
	}
	renderer := html.NewRenderer(opts)

	// Render the HTML
	htmlContent := markdown.Render(doc, renderer)

	// Return as template.HTML to prevent auto-escaping in template
	return template.HTML(htmlContent)
}

func main() {
	// initialize server
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Security headers
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "DENY",
		ContentSecurityPolicy: "default-src 'self'; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline' unpkg.com",
	}))

	e.Logger.SetLevel(log.DEBUG)
	e.Renderer = NewTemplate()

	// Static files
	e.Static("/static", "web/static")

	// Add favicon route
	e.GET("/favicon.ico", func(c echo.Context) error {
		return c.File("web/static/favicon.ico")
	})

	// Main route
	e.GET("/", func(c echo.Context) error {
		// create page data with a title based on the file name
		page_data := NewPage("./web/content/index.md", e)

		// Render the template with our data
		return c.Render(200, "template.html", page_data)
	})

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":  "healthy",
			"version": "1.0.0",
		})
	})

	// Graceful shutdown
	port := os.Getenv("PORT")
	if port == "" {
		port = "42069"
	}

	// start server
	e.Logger.Info("Starting server on port: " + port)
	e.Logger.Fatal(e.Start(":" + port))
}
