# About
- Staff Platform Engineer in FinTech |📍NYC
- I built this site primarily to learn Go and HTMX
- Below you'll find information about my personal projects
- [LinkedIn](https://linkedin.com/in/jack-p-coleman)

## Projects

### 1. Website
#### A web server built on Go, HTMX, and Markdown to display my projects
- [Source Code Repository](https://github.com/Lxkota95/website)

Everything starts with markdown files - including this page you're reading 👀

Next, the markdown is converted to HTML with the following function:

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
And lastly, a simple back-end web server renders this HTML using Go's template library:
```go
func main() {
	// initialize server
	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)
	e.Renderer = NewTemplate()

	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		// create page data with a title based on the file name
		page_data := NewPage("./content/index.md", e)

		// Render the template with our data
		return c.Render(200, "template.html", page_data)
	})

	// start server
	e.Logger.Info(e.Start(":42069"))
}
```
I took this approach because I don't enjoy front-end development 😉

---
### 2. Ansible
#### A Rust crate that offers an API for Ansible
- [Source Code Repository](https://github.com/Lxkota95/ansible)
- Published on Rust's crate repo [here](https://crates.io/crates/ansible/)

An example of loading Ansible inventory data for a host
```rust
use ansible::{Inventory, Load};

let inventory = Inventory::load(PathBuf::from('/path/to/inventory'))?;
let host = inventory.get_host("<hostname>")?;
hostvars = host.get_vars()?;
```
---
### 3. Ranked Bot - `#TODO`
#### A Discord Bot written in Python to assign roles based on your rank in Rainbow Six: Siege
- [Source Code Repository](https://github.com/Lxkota95/ranked)
- Example:
```python
```

---
### 4. Infra - `#TODO`
#### My infrastructure repo (and tools), built primarily using Ansible
- [Source Code Repository](https://github.com/Lxkota95/infra)

Here's a breakdown of my CI/CD process:
```yaml
```