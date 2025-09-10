package controlador

import (
	"encoding/json"
	"net/http"

	"github.com/MetaDandy/Assistense-System/helper"
	"github.com/MetaDandy/Assistense-System/src/modelo"
	"github.com/MetaDandy/Assistense-System/src/vista"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type DocenteControlador struct {
	modelos   modelo.DocenteInterfaz
	vistaHTML *vista.DocenteVistaHTML
}

type DocenteControladorInterfaz interface {
	// API endpoints (JSON)
	RegistrarDocente(w http.ResponseWriter, r *http.Request)
	IniciarSesion(w http.ResponseWriter, r *http.Request)
	ObtenerDocentePorID(w http.ResponseWriter, r *http.Request)
	ActualizarDocente(w http.ResponseWriter, r *http.Request)

	// HTML Templates (MVC Clásico)
	MostrarInicio(w http.ResponseWriter, r *http.Request)
	MostrarRegistro(w http.ResponseWriter, r *http.Request)
	ProcesarRegistro(w http.ResponseWriter, r *http.Request)
	MostrarLogin(w http.ResponseWriter, r *http.Request)
	ProcesarLogin(w http.ResponseWriter, r *http.Request)
}

func NuevoDocenteControlador(modelos modelo.DocenteInterfaz) DocenteControladorInterfaz {
	return &DocenteControlador{
		modelos:   modelos,
		vistaHTML: vista.NewDocenteVistaHTML(),
	}
}

func (dc *DocenteControlador) RegistrarDocente(w http.ResponseWriter, r *http.Request) {
	var registro modelo.RegistrarDocenteDto

	if err := json.NewDecoder(r.Body).Decode(&registro); err != nil {
		helper.EnviarJson(w, http.StatusBadRequest, map[string]string{
			"error": "JSON inválido",
		})
		return
	}

	docente, token, err := dc.modelos.RegistrarDocente(&registro)
	if err != nil {
		helper.EnviarJson(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	respuesta := vista.DocenteARespuesta(docente)
	respuesta.Token = token

	helper.EnviarJson(w, http.StatusCreated, respuesta)
}

func (dc *DocenteControlador) IniciarSesion(w http.ResponseWriter, r *http.Request) {
	var login modelo.IniciarSesionDto

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		helper.EnviarJson(w, http.StatusBadRequest, map[string]string{
			"error": "JSON inválido",
		})
		return
	}

	docente, token, err := dc.modelos.IniciarSesion(login)
	if err != nil {
		helper.EnviarJson(w, http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
		return
	}

	respuesta := vista.DocenteARespuesta(docente)
	respuesta.Token = token

	helper.EnviarJson(w, http.StatusOK, respuesta)
}

func (dc *DocenteControlador) ObtenerDocentePorID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	docenteID, err := uuid.Parse(id)
	if err != nil {
		helper.EnviarJson(w, http.StatusBadRequest, map[string]string{
			"error": "ID inválido",
		})
		return
	}

	docente, err := dc.modelos.ObtenerDocentePorID(docenteID)
	if err != nil {
		helper.EnviarJson(w, http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
		return
	}

	respuesta := vista.DocenteARespuesta(docente)
	helper.EnviarJson(w, http.StatusOK, respuesta)
}

func (dc *DocenteControlador) ActualizarDocente(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var actualizar modelo.ActualizarDocente
	if err := json.NewDecoder(r.Body).Decode(&actualizar); err != nil {
		helper.EnviarJson(w, http.StatusBadRequest, map[string]string{
			"error": "JSON inválido",
		})
		return
	}

	docenteID, err := uuid.Parse(id)
	if err != nil {
		helper.EnviarJson(w, http.StatusBadRequest, map[string]string{
			"error": "ID inválido",
		})
		return
	}

	docente, err := dc.modelos.ActualizarDocente(docenteID, &actualizar)
	if err != nil {
		helper.EnviarJson(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	respuesta := vista.DocenteARespuesta(docente)
	helper.EnviarJson(w, http.StatusOK, respuesta)
}

// ============================================================================
// MÉTODOS PARA HTML TEMPLATES (MVC CLÁSICO)
// ============================================================================

// MostrarInicio muestra la página principal
func (dc *DocenteControlador) MostrarInicio(w http.ResponseWriter, r *http.Request) {
	dc.vistaHTML.RenderizarInicio(w, nil)
}

// MostrarRegistro muestra el formulario de registro HTML
func (dc *DocenteControlador) MostrarRegistro(w http.ResponseWriter, r *http.Request) {
	dc.vistaHTML.RenderizarRegistro(w, nil)
}

// ProcesarRegistro procesa el formulario de registro HTML
func (dc *DocenteControlador) ProcesarRegistro(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/registro", http.StatusSeeOther)
		return
	}

	// Parsear formulario HTML
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error procesando formulario", http.StatusBadRequest)
		return
	}

	// Crear DTO desde formulario
	registro := modelo.RegistrarDocenteDto{
		Correo:     r.FormValue("correo"),
		Nombre:     r.FormValue("nombre"),
		Apellidos:  r.FormValue("apellidos"),
		Contraseña: r.FormValue("contraseña"),
	}

	// Validar confirmación de contraseña
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

	// Procesar registro usando el modelo
	docente, token, err := dc.modelos.RegistrarDocente(&registro)

	data := map[string]interface{}{}
	if err != nil {
		data["Error"] = err.Error()
		data["Correo"] = registro.Correo
		data["Nombre"] = registro.Nombre
		data["Apellidos"] = registro.Apellidos
	} else {
		data["Exito"] = true
		data["Docente"] = vista.DocenteARespuesta(docente)
		data["Token"] = token
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
	docente, token, err := dc.modelos.IniciarSesion(login)

	data := map[string]interface{}{}
	if err != nil {
		data["Error"] = err.Error()
		data["Correo"] = login.Correo
	} else {
		data["Exito"] = true
		data["Docente"] = vista.DocenteARespuesta(docente)
		data["Token"] = token
	}

	dc.vistaHTML.RenderizarLogin(w, data)
}
