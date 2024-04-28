package templates

import (
	"bytes"
	"context"
	"html/template"
	"io"
	"path/filepath"
	"sync"

	logutil "github.com/tangvis/erp/pkg/log"
)

// Container represents a set of templates that can be render
type Container struct {
	templates *template.Template
	mutex     sync.RWMutex
	//stop      chan struct{}
	//stopped   chan struct{}
	//watch     bool
}

// Data contains the data used to populate the template variables, it has Props
// that can be of any type and HTML that only can be `template.HTML` types.
type Data struct {
	Props map[string]interface{}
	HTML  map[string]template.HTML
}

// NewFromTemplate creates a new templates container using a
// `template.Template` object
func NewFromTemplate(templates *template.Template) *Container {
	return &Container{
		templates: templates,
		//mutex:     sync.RWMutex{},
	}
}

func NewDefaultTemplate() *Container {
	container, err := New("./templates")
	if err != nil {
		logutil.CtxErrorF(context.Background(), "NewDefaultTemplate failed: %+v", err)
	}
	return container
}

// New creates a new templates container scanning a directory.
func New(directory string) (*Container, error) {
	c := &Container{}

	htmlTemplates, err := template.ParseGlob(filepath.Join(directory, "*.html"))
	if err != nil {
		return nil, err
	}
	c.templates = htmlTemplates

	return c, nil
}

// RenderToString renders the template referenced with the template name using
// the data provided and return a string with the result
func (c *Container) RenderToString(templateName string, data Data) (string, error) {
	var text bytes.Buffer
	if err := c.Render(&text, templateName, data); err != nil {
		return "", err
	}
	return text.String(), nil
}

// Render renders the template referenced with the template name using
// the data provided and write it to the writer provided
func (c *Container) Render(w io.Writer, templateName string, data Data) error {
	c.mutex.RLock()
	htmlTemplates := c.templates
	c.mutex.RUnlock()

	if err := htmlTemplates.ExecuteTemplate(w, templateName, data); err != nil {
		return err
	}

	return nil
}
