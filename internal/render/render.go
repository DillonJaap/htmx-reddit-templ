package render

import (
	"html/template"
	"htmx-reddit/ui"
	"io/fs"
	"net/http"
	"path"

	"github.com/charmbracelet/log"
)

// TODO make this a closure?
type Renderer map[string]*template.Template

func New() Renderer {
	var renderer Renderer = make(Renderer)
	pages, err := fs.Glob(ui.Files, "html/pages/**.html")
	if err != nil {
		log.Fatal("couldn't glob files", "error", err)
	}

	morePages, err := fs.Glob(ui.Files, "html/pages/**/*.html")
	if err != nil {
		log.Fatal("couldn't glob files", "error", err)
	}
	pages = append(pages, morePages...)

	tmpl, err := template.ParseFS(
		ui.Files,
		"html/components/*.html",
		"html/base.html",
	)
	if err != nil {
		log.Fatal("couldn't glob files", "error", err)
	}
	renderer["components"] = tmpl

	for _, page := range pages {
		name := path.Base(page)
		log.Print("page", "page", page)
		tmpl, err := template.ParseFS(
			ui.Files,
			"html/components/*.html",
			"html/base.html",
			page,
		)
		if err != nil {
			log.Fatal("couldn't glob files", "error", err)
		}
		renderer[name] = tmpl
	}

	return renderer
}

// TODO should this return an error
func (r Renderer) RenderPage(w http.ResponseWriter, status int, name string, data interface{}) {
	if r[name] == nil {
		w.WriteHeader(http.StatusNotFound)
		log.Error("failed to get template", "page", name)
		return
	}
	if err := r[name].ExecuteTemplate(w, "base", data); err != nil {
		log.Error("failed to execute template", "error", err)
	}

	w.Header().Add("content-Type", "text/html; charset=UTF-8")
	if status != http.StatusOK {
		w.WriteHeader(status)
	}
}

func (r Renderer) RenderComponent(w http.ResponseWriter, status int, name string, data interface{}) {
	if _, ok := r["components"]; !ok {
		w.WriteHeader(http.StatusNotFound)
		log.Error("failed to get template", "page", name)
		return
	}
	if err := r["components"].ExecuteTemplate(w, name, data); err != nil {
		log.Error("failed to execute template", "error", err)
	}

	w.Header().Add("content-Type", "text/html; charset=UTF-8")
	if status != http.StatusOK {
		w.WriteHeader(status)
	}
}
