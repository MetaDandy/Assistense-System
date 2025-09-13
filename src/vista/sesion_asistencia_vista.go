package vista

import (
	"html/template"
	"net/http"
)

type SesionAsistenciaVistaHTML struct {
	tmpl *template.Template
}

func NuevaSesionAsistenciaVistaHTML() *SesionAsistenciaVistaHTML {
	t := template.Must(template.ParseFS(TemplatesFS, "templates/*.html"))
	return &SesionAsistenciaVistaHTML{tmpl: t}
}

func (v *SesionAsistenciaVistaHTML) RenderizarRegistrar(w http.ResponseWriter, data interface{}) {
	v.tmpl.ExecuteTemplate(w, "registrar_sesion_asistencia.html", data)
}

func (v *SesionAsistenciaVistaHTML) RenderizarListar(w http.ResponseWriter, data interface{}) {
	v.tmpl.ExecuteTemplate(w, "listar_sesiones_asistencia.html", data)
}

func (v *SesionAsistenciaVistaHTML) RenderizarDetalle(w http.ResponseWriter, data interface{}) {
	v.tmpl.ExecuteTemplate(w, "detalle_sesion_asistencia.html", data)
}
