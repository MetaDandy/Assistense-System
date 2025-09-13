package vista

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed templates/*.html
var TemplatesFS embed.FS

type DocenteVistaHTML struct {
	tmpl *template.Template
}

func NuevoDocenteVistaHTML() *DocenteVistaHTML {
	t := template.Must(template.ParseFS(TemplatesFS, "templates/*.html"))
	return &DocenteVistaHTML{tmpl: t}
}

// RenderizarInicio renderiza la página principal
func (dv *DocenteVistaHTML) RenderizarInicio(w http.ResponseWriter, data interface{}) {
	if err := dv.tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RenderizarRegistro renderiza el formulario de registro
func (dv *DocenteVistaHTML) RenderizarRegistro(w http.ResponseWriter, data interface{}) {
	if err := dv.tmpl.ExecuteTemplate(w, "registro.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RenderizarLogin renderiza el formulario de login
func (dv *DocenteVistaHTML) RenderizarLogin(w http.ResponseWriter, data interface{}) {
	if err := dv.tmpl.ExecuteTemplate(w, "login.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RenderizarPanelDocente renderiza la vista del panel de docente
func (dv *DocenteVistaHTML) RenderizarPanelDocente(w http.ResponseWriter, data interface{}) {
	if err := dv.tmpl.ExecuteTemplate(w, "panel_docente.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
