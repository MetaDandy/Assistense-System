package controlador

import (
	"net/http"
	"time"

	"github.com/MetaDandy/Assistense-System/helper"
	"github.com/MetaDandy/Assistense-System/src/modelo"
	"github.com/MetaDandy/Assistense-System/src/vista"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type SesionAsistenciaControlador struct {
	modelo modelo.SesionAsistenciaInterfaz
	vista  *vista.SesionAsistenciaVistaHTML
}

func NuevoSesionAsistenciaControlador(m modelo.SesionAsistenciaInterfaz, v *vista.SesionAsistenciaVistaHTML) *SesionAsistenciaControlador {
	return &SesionAsistenciaControlador{modelo: m, vista: v}
}

// GET /sesion-asistencia/registrar
func (c *SesionAsistenciaControlador) MostrarRegistrar(w http.ResponseWriter, r *http.Request) {
	c.vista.RenderizarRegistrar(w, nil)
}

// POST /sesion-asistencia/registrar
func (c *SesionAsistenciaControlador) ProcesarRegistrar(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.vista.RenderizarRegistrar(w, map[string]interface{}{"Error": "Error en el formulario"})
		return
	}

	// Obtener strings directamente del formulario SIN parsear
	fechaStr := r.FormValue("fecha")            // "2025-09-13"
	horaInicioStr := r.FormValue("hora_inicio") // "11:33"
	horaFinStr := r.FormValue("hora_fin")       // "12:33"

	// Obtener el DocenteID desde el JWT en la cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	claims, err := helper.ValidateJwt(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	docenteIDStr, ok := claims["id"].(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	docenteID, err := uuid.Parse(docenteIDStr)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Enviar strings directamente al modelo
	dto := &modelo.RegistrarSesionAsistenciaDto{
		Fecha:      fechaStr,      // "2025-09-13"
		HoraInicio: horaInicioStr, // "11:33"
		HoraFin:    horaFinStr,    // "12:33"
		DocenteID:  docenteID,
	}

	_, err = c.modelo.RegistrarSesionAsistencia(dto)
	if err != nil {
		c.vista.RenderizarRegistrar(w, map[string]interface{}{"Error": "No se pudo registrar la sesión"})
		return
	}
	c.vista.RenderizarRegistrar(w, map[string]interface{}{"Exito": true})
}

// GET /sesion-asistencia/listar
func (c *SesionAsistenciaControlador) ListarSesiones(w http.ResponseWriter, r *http.Request) {
	// Obtener el DocenteID desde el JWT en la cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	claims, err := helper.ValidateJwt(cookie.Value)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	docenteIDStr, ok := claims["id"].(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	docenteID, err := uuid.Parse(docenteIDStr)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	sesiones, _ := c.modelo.ObtenerSesionesAsistencia(docenteID)
	type SesionView struct {
		ID         string
		Fecha      string
		HoraInicio string
		HoraFin    string
		Activa     bool
	}
	var sesionesView []SesionView
	now := time.Now()

	for _, s := range sesiones {
		// Comparar hora actual con hora de sesión (simplificado)
		horaActual := now.Format("15:04")
		activa := horaActual >= s.HoraInicio && horaActual <= s.HoraFin

		sesionesView = append(sesionesView, SesionView{
			ID:         s.ID.String(),
			Fecha:      s.Fecha,
			HoraInicio: s.HoraInicio,
			HoraFin:    s.HoraFin,
			Activa:     activa,
		})
	}
	c.vista.RenderizarListar(w, map[string]interface{}{"Sesiones": sesionesView})
}

// GET /sesion-asistencia/{id}
func (c *SesionAsistenciaControlador) MostrarDetalle(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	sesion, err := c.modelo.ObtenerSesionAsistencia(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	now := time.Now()
	horaActual := now.Format("15:04")
	activa := horaActual >= sesion.HoraInicio && horaActual <= sesion.HoraFin
	if !activa {
		http.Error(w, "Sesión no activa", http.StatusForbidden)
		return
	}
	c.vista.RenderizarDetalle(w, map[string]interface{}{"Sesion": sesion})
}
