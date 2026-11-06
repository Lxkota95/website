package site

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config is site.toml: the structured parts of the site. Prose lives in markdown files.
type Config struct {
	Name        string    `toml:"name"`
	Role        string    `toml:"role"`
	Location    string    `toml:"location"`
	Description string    `toml:"description"`
	Projects    []Project `toml:"projects"`
	Toolbox     []Pair    `toml:"toolbox"`
	Contact     []Link    `toml:"contact"`
	Learning    []Step    `toml:"learning"` // the learning path, in order
}

// Project is one row of the projects table.
type Project struct {
	Name      string `toml:"name"`
	Lang      string `toml:"lang"`
	State     string `toml:"state"`  // shown text, e.g. "v0.3.0" or "active"
	Status    string `toml:"status"` // drives the colour: one of validStatuses
	What      string `toml:"what"`
	Repo      string `toml:"repo"`       // optional: links the project name
	LinkLabel string `toml:"link_label"` // optional: the right-hand link
	Link      string `toml:"link"`
	Next      string `toml:"next"` // optional: what's planned next, shown when the row is opened
}

// Step is one entry of the learning path.
type Step struct {
	Name    string `toml:"name"`
	Status  string `toml:"status"` // one of validStepStatuses
	What    string `toml:"what"`
	Project string `toml:"project"` // optional: the project it's built in; links to its row
}

// Pair is a key/value line, as in the toolbox section.
type Pair struct {
	Key   string `toml:"key"`
	Value string `toml:"value"`
}

// Link is a labelled URL, as in the contact section.
type Link struct {
	Key   string `toml:"key"`
	Label string `toml:"label"`
	URL   string `toml:"url"`
}

var (
	validStatuses     = []string{"release", "active", "learning", "design", "done"}
	validStepStatuses = []string{"now", "planned", "idea"}
)

// LoadConfig reads and validates site.toml. Unknown keys are errors, so a typo (say `staus`)
// fails the build instead of silently dropping the value.
func LoadConfig(path string) (*Config, error) {
	var c Config
	md, err := toml.DecodeFile(path, &c)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}
	if undecoded := md.Undecoded(); len(undecoded) > 0 {
		keys := make([]string, len(undecoded))
		for i, k := range undecoded {
			keys[i] = k.String()
		}
		return nil, fmt.Errorf("%s: unknown keys: %s", path, strings.Join(keys, ", "))
	}
	if err := c.validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return &c, nil
}

func (c *Config) validate() error {
	var errs []error
	if c.Name == "" {
		errs = append(errs, errors.New("name is required"))
	}
	for i, p := range c.Projects {
		if p.Name == "" || p.Lang == "" || p.State == "" {
			errs = append(errs, fmt.Errorf("projects[%d]: name, lang and state are required", i))
		}
		if !slices.Contains(validStatuses, p.Status) {
			errs = append(errs, fmt.Errorf("projects[%d] (%s): status %q is not one of %s",
				i, p.Name, p.Status, strings.Join(validStatuses, ", ")))
		}
		if (p.LinkLabel == "") != (p.Link == "") {
			errs = append(errs, fmt.Errorf("projects[%d] (%s): set both link_label and link, or neither",
				i, p.Name))
		}
	}
	for i, s := range c.Learning {
		if s.Name == "" || s.What == "" {
			errs = append(errs, fmt.Errorf("learning[%d]: name and what are required", i))
		}
		if !slices.Contains(validStepStatuses, s.Status) {
			errs = append(errs, fmt.Errorf("learning[%d] (%s): status %q is not one of %s",
				i, s.Name, s.Status, strings.Join(validStepStatuses, ", ")))
		}
		if s.Project != "" && !slices.ContainsFunc(c.Projects, func(p Project) bool {
			return p.Name == s.Project
		}) {
			errs = append(errs, fmt.Errorf("learning[%d] (%s): project %q is not in projects",
				i, s.Name, s.Project))
		}
	}
	for i, l := range c.Contact {
		if l.Key == "" || l.Label == "" || l.URL == "" {
			errs = append(errs, fmt.Errorf("contact[%d]: key, label and url are required", i))
		}
	}
	return errors.Join(errs...)
}
