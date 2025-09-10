package vista

import (
	"html/template"
	"net/http"
	"time"

	"github.com/MetaDandy/Assistense-System/src/modelo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type DocenteRespuesta struct {
	ID        uuid.UUID `json:"id"`
	Correo    string    `json:"correo"`
	Nombre    string    `json:"nombre"`
	Apellidos string    `json:"apellidos"`
	Token     string    `json:"token,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DocenteARespuesta(u *modelo.Docente) *DocenteRespuesta {
	var docente DocenteRespuesta

	copier.Copy(&docente, u)

	return &docente
}

type DocenteVistaHTML struct{}

func NewDocenteVistaHTML() *DocenteVistaHTML {
	return &DocenteVistaHTML{}
}

// RenderizarInicio renderiza la página principal
func (dv *DocenteVistaHTML) RenderizarInicio(w http.ResponseWriter, data interface{}) {
	tmpl := template.Must(template.ParseFiles("src/vista/templates/index.html"))
	tmpl.Execute(w, data)
}

// RenderizarRegistro renderiza el formulario de registro
func (dv *DocenteVistaHTML) RenderizarRegistro(w http.ResponseWriter, data interface{}) {
	tmpl := template.Must(template.ParseFiles("src/vista/templates/registro.html"))
	tmpl.Execute(w, data)
}

// RenderizarLogin renderiza el formulario de login
func (dv *DocenteVistaHTML) RenderizarLogin(w http.ResponseWriter, data interface{}) {
	tmpl := template.Must(template.ParseFiles("src/vista/templates/login.html"))
	tmpl.Execute(w, data)
}
