package sender

import (
	"html/template"
	"io"
	"io/fs"
)

type Template struct {
	debug    bool
	fs       fs.FS
	patterns []string
	name     string
	template *template.Template
}

func NewTemplate(fs fs.FS, patterns ...string) *Template {
	return &Template{
		fs:       fs,
		patterns: patterns,
		name:     patterns[len(patterns)-1],
		template: refreshTemplate(fs, patterns...),
	}
}

func (t *Template) EnableDebug() {
	t.debug = true
}

func (t *Template) DisableDebug() {
	t.debug = false
}

func (t *Template) Template() *template.Template {
	if t.debug {
		return refreshTemplate(t.fs, t.patterns...)
	}

	return t.template
}

func (t *Template) ExecuteTemplate(wr io.Writer, data any) error {
	return t.Template().ExecuteTemplate(wr, t.name, data)
}

func refreshTemplate(fs fs.FS, patterns ...string) *template.Template {
	return template.Must(template.ParseFS(fs, patterns...))
}
