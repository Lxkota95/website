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
		CSS:   "highlight", // Add CSS class for syntax highlighting
	}
	renderer := html.NewRenderer(opts)

	// Render the HTML
	htmlContent := markdown.Render(doc, renderer)

	// Post-process HTML to add language classes to code blocks
	htmlString := string(htmlContent)
	htmlString = enhanceCodeBlocks(htmlString)

	// Return as template.HTML to prevent auto-escaping in template
	return template.HTML(htmlString)
}

// enhanceCodeBlocks adds language classes and other enhancements to code blocks
func enhanceCodeBlocks(html string) string {
	// This is a simple enhancement - in production you might want to use a proper HTML parser
	// Add language-specific classes for syntax highlighting
	html = strings.ReplaceAll(html, `<pre><code class="language-go">`, `<pre><code class="language-go highlight-go">`)
	html = strings.ReplaceAll(html, `<pre><code class="language-rust">`, `<pre><code class="language-rust highlight-rust">`)
	html = strings.ReplaceAll(html, `<pre><code class="language-python">`, `<pre><code class="language-python highlight-python">`)
	html = strings.ReplaceAll(html, `<pre><code class="language-yaml">`, `<pre><code class="language-yaml highlight-yaml">`)
	html = strings.ReplaceAll(html, `<pre><code class="language-javascript">`, `<pre><code class="language-javascript highlight-js">`)
	html = strings.ReplaceAll(html, `<pre><code class="language-html">`, `<pre><code class="language-html highlight-html">`)
	html = strings.ReplaceAll(html, `<pre><code class="language-css">`, `<pre><code class="language-css highlight-css">`)

	return html
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

	// API route for potential HTMX interactions
	e.GET("/api/projects", func(c echo.Context) error {
		// This could return project data as JSON or HTML fragments
		// Useful for dynamic loading with HTMX
		return c.JSON(200, map[string]interface{}{
			"projects": []map[string]interface{}{
				{
					"name":        "Website",
					"description": "A web server built on Go, HTMX, and Markdown",
					"status":      "active",
					"languages":   []string{"Go", "HTML", "CSS"},
				},
				{
					"name":        "Ansible",
					"description": "A Rust crate that offers an API for Ansible",
					"status":      "active",
					"languages":   []string{"Rust"},
				},
				{
					"name":        "Ranked Bot",
					"description": "A Discord Bot for Rainbow Six: Siege ranks",
					"status":      "todo",
					"languages":   []string{"Python"},
				},
				{
					"name":        "Infra",
					"description": "Infrastructure repo built with Ansible",
					"status":      "todo",
					"languages":   []string{"YAML", "Ansible"},
				},
			},
		})
	})

	// Project details handler for HTMX
	e.GET("/api/project-details/:project", func(c echo.Context) error {
		projectName := c.Param("project")

		var details string
		switch projectName {
		case "website":
			details = `
        <div class="project-timeline">
            <div class="timeline-item">
                <div class="timeline-date">Jan 2025</div>
                <div class="timeline-content">
                    <strong>Frontend Redesign</strong> - Complete visual overhaul with modern CSS and HTMX interactions
                </div>
            </div>
            <div class="timeline-item">
                <div class="timeline-date">Dec 2024</div>
                <div class="timeline-content">
                    <strong>Initial Development</strong> - Built core Go server with Echo framework and markdown processing
                </div>
            </div>
            <div class="timeline-item">
                <div class="timeline-date">Nov 2024</div>
                <div class="timeline-content">
                    <strong>Planning & Design</strong> - Decided on Go + HTMX + Markdown architecture
                </div>
            </div>
        </div>
        
        <div class="project-stats">
            <div class="project-stat">
                <span>📊</span>
                <span class="project-stat-value">~500</span>
                <span>lines of code</span>
            </div>
            <div class="project-stat">
                <span>⚡</span>
                <span class="project-stat-value">&lt;50ms</span>
                <span>response time</span>
            </div>
            <div class="project-stat">
                <span>🚀</span>
                <span class="project-stat-value">Fly.io</span>
                <span>deployment</span>
            </div>
        </div>
        
        <div style="margin-top: var(--space-lg);">
            <h4>Technologies Used</h4>
            <div style="display: flex; flex-wrap: wrap; gap: var(--space-xs); margin-top: var(--space-sm);">
                <span class="skill-tag">Go</span>
                <span class="skill-tag">Echo Framework</span>
                <span class="skill-tag">HTMX</span>
                <span class="skill-tag">Markdown</span>
                <span class="skill-tag">Docker</span>
                <span class="skill-tag">GitHub Actions</span>
                <span class="skill-tag">Fly.io</span>
            </div>
        </div>`

		case "ansible":
			details = `
        <div class="project-timeline">
            <div class="timeline-item">
                <div class="timeline-date">Dec 2024</div>
                <div class="timeline-content">
                    <strong>v0.2.0 Release</strong> - Added support for dynamic inventory and improved error handling
                </div>
            </div>
            <div class="timeline-item">
                <div class="timeline-date">Nov 2024</div>
                <div class="timeline-content">
                    <strong>v0.1.0 Release</strong> - Initial crate with basic inventory loading functionality
                </div>
            </div>
            <div class="timeline-item">
                <div class="timeline-date">Oct 2024</div>
                <div class="timeline-content">
                    <strong>Development Started</strong> - Need arose from lack of good Rust libraries for Ansible
                </div>
            </div>
        </div>
        
        <div class="project-stats">
            <div class="project-stat">
                <span>📦</span>
                <span class="project-stat-value">50+</span>
                <span>downloads</span>
            </div>
            <div class="project-stat">
                <span>⭐</span>
                <span class="project-stat-value">5</span>
                <span>GitHub stars</span>
            </div>
            <div class="project-stat">
                <span>🔧</span>
                <span class="project-stat-value">0</span>
                <span>open issues</span>
            </div>
        </div>
        
        <div style="margin-top: var(--space-lg);">
            <h4>Key Features</h4>
            <ul style="margin-top: var(--space-sm);">
                <li>Type-safe inventory loading</li>
                <li>Support for multiple inventory formats</li>
                <li>Host variable resolution</li>
                <li>Group membership queries</li>
                <li>Comprehensive error handling</li>
            </ul>
        </div>`

		case "ranked-bot":
			details = `
        <div class="project-timeline">
            <div class="timeline-item">
                <div class="timeline-date">Planned</div>
                <div class="timeline-content">
                    <strong>Beta Release</strong> - Initial bot with basic rank tracking
                </div>
            </div>
            <div class="timeline-item">
                <div class="timeline-date">In Progress</div>
                <div class="timeline-content">
                    <strong>API Integration</strong> - Working on Rainbow Six Siege API integration
                </div>
            </div>
        </div>
        
        <div style="margin-top: var(--space-lg);">
            <h4>Planned Features</h4>
            <ul style="margin-top: var(--space-sm);">
                <li>Automatic role assignment based on rank</li>
                <li>Periodic rank updates</li>
                <li>Multiple game mode support</li>
                <li>Customizable role mappings</li>
                <li>Statistics tracking</li>
            </ul>
        </div>`

		case "infra":
			details = `
        <div class="project-timeline">
            <div class="timeline-item">
                <div class="timeline-date">Planned</div>
                <div class="timeline-content">
                    <strong>Production Ready</strong> - Complete infrastructure automation
                </div>
            </div>
            <div class="timeline-item">
                <div class="timeline-date">In Progress</div>
                <div class="timeline-content">
                    <strong>Initial Setup</strong> - Basic Ansible playbooks and CI/CD
                </div>
            </div>
        </div>
        
        <div style="margin-top: var(--space-lg);">
            <h4>Infrastructure Components</h4>
            <ul style="margin-top: var(--space-sm);">
                <li>Automated server provisioning</li>
                <li>CI/CD pipeline templates</li>
                <li>Monitoring and alerting</li>
                <li>Backup and disaster recovery</li>
                <li>Security hardening</li>
            </ul>
        </div>`

		default:
			details = `<p>Project details not found.</p>`
		}

		return c.HTML(200, details)
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
