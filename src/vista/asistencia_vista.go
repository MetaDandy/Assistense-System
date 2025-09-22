package vista

import (
	"html/template"
	"net/http"
)

type AsistenciaVistaHTML struct {
	tmpl *template.Template
}

func NuevaAsistenciaVistaHTML() *AsistenciaVistaHTML {
	t := template.Must(template.ParseFS(TemplatesFS, "templates/*.html"))
	return &AsistenciaVistaHTML{tmpl: t}
}

func (v *AsistenciaVistaHTML) RenderizarConfirmarAsistencia(w http.ResponseWriter, data interface{}) {
	if err := v.tmpl.ExecuteTemplate(w, "confirmar_asistencia.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (v *AsistenciaVistaHTML) RenderizarCapturarFoto(w http.ResponseWriter, data interface{}) {
	if err := v.tmpl.ExecuteTemplate(w, "capturar_foto.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (v *AsistenciaVistaHTML) RenderizarListarAsistencias(w http.ResponseWriter, data interface{}) {
	if err := v.tmpl.ExecuteTemplate(w, "listar_asistencias.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
