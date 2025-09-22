package vista

import (
	"html/template"
	"net/http"
)

type EstudianteVistaHTML struct {
	tmpl *template.Template
}

func NuevaEstudianteVistaHTML() *EstudianteVistaHTML {
	t := template.Must(template.ParseFS(TemplatesFS, "templates/*.html"))
	return &EstudianteVistaHTML{tmpl: t}
}

// RenderizarGestionarEstudiantes renderiza la vista para gestionar estudiantes
func (ev *EstudianteVistaHTML) RenderizarGestionarEstudiantes(w http.ResponseWriter, data interface{}) {
	if err := ev.tmpl.ExecuteTemplate(w, "gestionar_estudiantes.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// RenderizarEditarEstudiante renderiza la vista para editar un estudiante
func (ev *EstudianteVistaHTML) RenderizarEditarEstudiante(w http.ResponseWriter, data interface{}) {
	if err := ev.tmpl.ExecuteTemplate(w, "editar_estudiante.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
