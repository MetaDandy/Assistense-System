package controlador

import (
	"net/http"

	"github.com/MetaDandy/Assistense-System/src/modelo"
	"github.com/MetaDandy/Assistense-System/src/vista"
)

type DocenteControlador struct {
	modelos   modelo.DocenteInterfaz
	vistaHTML *vista.DocenteVistaHTML
}

type DocenteControladorInterfaz interface {
	MostrarInicio(w http.ResponseWriter, r *http.Request)
	MostrarRegistro(w http.ResponseWriter, r *http.Request)
	ProcesarRegistro(w http.ResponseWriter, r *http.Request)
	MostrarLogin(w http.ResponseWriter, r *http.Request)
	ProcesarLogin(w http.ResponseWriter, r *http.Request)
	MostrarPanelDocente(w http.ResponseWriter, r *http.Request)
}

func NuevoDocenteControlador(modelos modelo.DocenteInterfaz) DocenteControladorInterfaz {
	return &DocenteControlador{
		modelos:   modelos,
		vistaHTML: vista.NewDocenteVistaHTML(),
	}
}

func (dc *DocenteControlador) MostrarInicio(w http.ResponseWriter, r *http.Request) {
	dc.vistaHTML.RenderizarInicio(w, nil)
}

func (dc *DocenteControlador) MostrarRegistro(w http.ResponseWriter, r *http.Request) {
	dc.vistaHTML.RenderizarRegistro(w, nil)
}

// ProcesarRegistro procesa el formulario de registro HTML
func (dc *DocenteControlador) ProcesarRegistro(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/registro", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error procesando formulario", http.StatusBadRequest)
		return
	}

	registro := modelo.RegistrarDocenteDto{
		Correo:     r.FormValue("correo"),
		Nombre:     r.FormValue("nombre"),
		Apellidos:  r.FormValue("apellidos"),
		Contraseña: r.FormValue("contraseña"),
	}

	confirmarContraseña := r.FormValue("confirmar_contraseña")
	if registro.Contraseña != confirmarContraseña {
		data := map[string]interface{}{
			"Error":     "Las contraseñas no coinciden",
			"Correo":    registro.Correo,
			"Nombre":    registro.Nombre,
			"Apellidos": registro.Apellidos,
		}
		dc.vistaHTML.RenderizarRegistro(w, data)
		return
	}

	docente, token, err := dc.modelos.RegistrarDocente(&registro)
	if err != nil {
		data := map[string]interface{}{
			"Error": err.Error(),
		}
		dc.vistaHTML.RenderizarRegistro(w, data)
		return
	}

	// Enviar el ID del docente y el token al cliente
	data := map[string]interface{}{
		"Exito":   true,
		"Docente": docente,
		"Token":   token,
	}
	dc.vistaHTML.RenderizarRegistro(w, data)
}

// MostrarLogin muestra el formulario de login HTML
func (dc *DocenteControlador) MostrarLogin(w http.ResponseWriter, r *http.Request) {
	dc.vistaHTML.RenderizarLogin(w, nil)
}

// ProcesarLogin procesa el formulario de login HTML
func (dc *DocenteControlador) ProcesarLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Parsear formulario HTML
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error procesando formulario", http.StatusBadRequest)
		return
	}

	// Crear DTO desde formulario
	login := modelo.IniciarSesionDto{
		Correo:     r.FormValue("correo"),
		Contraseña: r.FormValue("contraseña"),
	}

	// Procesar login usando el modelo
	_, _, err := dc.modelos.IniciarSesion(login)
	if err != nil {
		data := map[string]interface{}{
			"Error":  err.Error(),
			"Correo": login.Correo,
		}
		dc.vistaHTML.RenderizarLogin(w, data)
		return
	}

	// Redirigir al panel del docente
	http.Redirect(w, r, "/panel-docente", http.StatusSeeOther)
}

func (dc *DocenteControlador) MostrarPanelDocente(w http.ResponseWriter, r *http.Request) {
	dc.vistaHTML.RenderizarPanelDocente(w, nil)
}
