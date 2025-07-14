# About

I'm a Staff Platform Engineer working in FinTech here in NYC, focused on reliability, performance, and tooling. I built this site primarily to learn Go and HTMX, but it's also become a place where I can showcase my personal projects and experiments.

Below you'll find information about my personal projects, ranging from web servers to infrastructure tooling. Each project represents a different learning journey and technical challenge I've taken on.

**Connect with me:** [LinkedIn](https://linkedin.com/in/jack-p-coleman)

---

## Projects

### 1. Website
<div class="project-header">
  <h4 class="project-description">A web server built on Go, HTMX, and Markdown to display my projects</h4>
  <span class="project-status">Status: Active</span>
</div>

<div class="project-links">
  <a href="https://github.com/Lxkota95/website" class="project-link">
    📁 Source Code
  </a>
</div>

Everything starts with markdown files - including this page you're reading right now 👀

The architecture is simple: markdown gets converted to HTML, then rendered through Go's template system. Here's the heart of the markdown conversion:

```go
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
```

And the web server itself is straightforward:

```go
func main() {
	// initialize server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)
	e.Renderer = NewTemplate()

	e.Static("/static", "web/static")

	e.GET("/", func(c echo.Context) error {
		// create page data with a title based on the file name
		page_data := NewPage("./web/content/index.md", e)

		// Render the template with our data
		return c.Render(200, "template.html", page_data)
	})

	// start server
	e.Logger.Info(e.Start(":42069"))
}
```

> I took this approach because I don't enjoy front-end development 😉  
> Though I have to admit, working on this redesign has been pretty fun!

### 2. Ansible
<div class="project-header">
  <h4 class="project-description">A Rust crate that offers an API for Ansible</h4>
  <span class="project-status">Status: Active</span>
</div>

<div class="project-links">
  <a href="https://github.com/Lxkota95/ansible" class="project-link">
    📁 Source Code
  </a>
  <a href="https://crates.io/crates/ansible/" class="project-link">
    📦 Crates.io
  </a>
</div>

This crate provides a clean, idiomatic Rust API for working with Ansible inventory data. It's particularly useful when you need to programmatically access host variables and inventory information from within Rust applications.

**Example usage:**
```rust
use ansible::{Inventory, Load};

let inventory = Inventory::load(PathBuf::from('/path/to/inventory'))?;
let host = inventory.get_host("")?;
let hostvars = host.get_vars()?;
```

The crate handles all the complexity of parsing Ansible's various inventory formats and provides a consistent interface for accessing the data you need.

### 3. Ranked Bot
<div class="project-header">
  <h4 class="project-description">A Discord Bot written in Python to assign roles based on your rank in Rainbow Six: Siege</h4>
  <span class="project-status">Status: TODO</span>
</div>

<div class="project-links">
  <a href="https://github.com/Lxkota95/ranked" class="project-link">
    📁 Source Code
  </a>
</div>

This Discord bot integrates with the Rainbow Six: Siege API to automatically assign Discord roles based on players' competitive ranks. It's designed to help gaming communities organize their members by skill level.

**Features in development:**
- Automatic role assignment based on current rank
- Periodic rank updates
- Support for multiple game modes
- Customizable role mappings

```python
# Example implementation coming soon!
```

### 4. Infra
<div class="project-header">
  <h4 class="project-description">My infrastructure repo and tools, built primarily using Ansible</h4>
  <span class="project-status">Status: TODO</span>
</div>

<div class="project-links">
  <a href="https://github.com/Lxkota95/infra" class="project-link">
    📁 Source Code
  </a>
</div>

This repository contains my personal infrastructure setup, including CI/CD pipelines, server configurations, and deployment automation. It's built around Ansible for configuration management and includes various custom tools and scripts.

**Planned features:**
- Automated server provisioning
- CI/CD pipeline templates
- Monitoring and alerting setup
- Backup and disaster recovery procedures

```yaml
# CI/CD pipeline configuration coming soon!
```

---

## Technical Notes

This site is deployed on [Fly.io](https://fly.io) using their container deployment platform. The entire build and deployment process is automated through GitHub Actions, making updates as simple as pushing to the main branch.