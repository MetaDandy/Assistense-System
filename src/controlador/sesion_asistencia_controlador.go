package controlador

import (
	"net/http"

	"github.com/MetaDandy/Assistense-System/helper"
	"github.com/MetaDandy/Assistense-System/src/controlador/sesion_estado"
	"github.com/MetaDandy/Assistense-System/src/modelo"
	"github.com/MetaDandy/Assistense-System/src/vista"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type SesionAsistenciaControladorInterfaz interface {
	MostrarRegistrar(w http.ResponseWriter, r *http.Request)
	ProcesarRegistrar(w http.ResponseWriter, r *http.Request)
	ListarSesiones(w http.ResponseWriter, r *http.Request)
	MostrarDetalle(w http.ResponseWriter, r *http.Request)
	MostrarGestionarSesiones(w http.ResponseWriter, r *http.Request)
	ProcesarGestionarSesiones(w http.ResponseWriter, r *http.Request)
	MostrarRegistrarAsistencias(w http.ResponseWriter, r *http.Request)
	ProcesarSeleccionEstudiante(w http.ResponseWriter, r *http.Request)
	MostrarFormularioFoto(w http.ResponseWriter, r *http.Request)
}

type SesionAsistenciaControlador struct {
	modelo           modelo.SesionAsistenciaInterfaz
	estudianteModelo modelo.EstudianteModeloInterfaz
	vista            *vista.SesionAsistenciaVistaHTML
}

func NuevoSesionAsistenciaControlador(m modelo.SesionAsistenciaInterfaz, em modelo.EstudianteModeloInterfaz, v *vista.SesionAsistenciaVistaHTML) SesionAsistenciaControladorInterfaz {
	return &SesionAsistenciaControlador{
		modelo:           m,
		estudianteModelo: em,
		vista:            v,
	}
}

func (c *SesionAsistenciaControlador) MostrarRegistrar(w http.ResponseWriter, r *http.Request) {
	c.vista.RenderizarRegistrar(w, nil)
}

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

	for _, s := range sesiones {
		// Verificar si la sesión está activa usando el patrón State
		ctx := &sesion_estado.Sesion{
			Fecha:      s.Fecha,
			HoraInicio: s.HoraInicio,
			HoraFin:    s.HoraFin,
		}
		activa := ctx.CanRegistrarAsistencia()

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

	// Verificar si la sesión está activa usando el patrón State
	ctx := &sesion_estado.Sesion{
		Fecha:      sesion.Fecha,
		HoraInicio: sesion.HoraInicio,
		HoraFin:    sesion.HoraFin,
	}
	activa := ctx.CanRegistrarAsistencia()

	data := map[string]interface{}{
		"Sesion": sesion,
		"Activa": activa,
	}

	c.vista.RenderizarDetalle(w, data)
}

func (c *SesionAsistenciaControlador) MostrarGestionarSesiones(w http.ResponseWriter, r *http.Request) {
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

	for _, s := range sesiones {
		// Verificar si la sesión está activa usando el patrón State
		ctx := &sesion_estado.Sesion{
			Fecha:      s.Fecha,
			HoraInicio: s.HoraInicio,
			HoraFin:    s.HoraFin,
		}
		activa := ctx.CanRegistrarAsistencia()

		sesionesView = append(sesionesView, SesionView{
			ID:         s.ID.String(),
			Fecha:      s.Fecha,
			HoraInicio: s.HoraInicio,
			HoraFin:    s.HoraFin,
			Activa:     activa,
		})
	}
	c.vista.RenderizarGestionarSesiones(w, map[string]interface{}{"Sesiones": sesionesView})
}

func (c *SesionAsistenciaControlador) ProcesarGestionarSesiones(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		c.renderGestionarConError(w, r, "Error en el formulario")
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
		c.renderGestionarConError(w, r, "No se pudo registrar la sesión")
		return
	}
	c.renderGestionarConExito(w, r)
}

func (c *SesionAsistenciaControlador) renderGestionarConError(w http.ResponseWriter, r *http.Request, mensaje string) {
	// Recargar las sesiones para mostrar la lista actualizada
	cookie, _ := r.Cookie("token")
	claims, _ := helper.ValidateJwt(cookie.Value)
	docenteIDStr, _ := claims["id"].(string)
	docenteID, _ := uuid.Parse(docenteIDStr)

	sesiones, _ := c.modelo.ObtenerSesionesAsistencia(docenteID)
	type SesionView struct {
		ID         string
		Fecha      string
		HoraInicio string
		HoraFin    string
		Activa     bool
	}
	var sesionesView []SesionView

	for _, s := range sesiones {
		// Verificar si la sesión está activa usando el patrón State
		ctx := &sesion_estado.Sesion{
			Fecha:      s.Fecha,
			HoraInicio: s.HoraInicio,
			HoraFin:    s.HoraFin,
		}
		activa := ctx.CanRegistrarAsistencia()
		sesionesView = append(sesionesView, SesionView{
			ID:         s.ID.String(),
			Fecha:      s.Fecha,
			HoraInicio: s.HoraInicio,
			HoraFin:    s.HoraFin,
			Activa:     activa,
		})
	}
	c.vista.RenderizarGestionarSesiones(w, map[string]interface{}{
		"Sesiones": sesionesView,
		"Error":    mensaje,
	})
}

func (c *SesionAsistenciaControlador) renderGestionarConExito(w http.ResponseWriter, r *http.Request) {
	// Recargar las sesiones para mostrar la lista actualizada
	cookie, _ := r.Cookie("token")
	claims, _ := helper.ValidateJwt(cookie.Value)
	docenteIDStr, _ := claims["id"].(string)
	docenteID, _ := uuid.Parse(docenteIDStr)

	sesiones, _ := c.modelo.ObtenerSesionesAsistencia(docenteID)
	type SesionView struct {
		ID         string
		Fecha      string
		HoraInicio string
		HoraFin    string
		Activa     bool
	}
	var sesionesView []SesionView

	for _, s := range sesiones {
		// Verificar si la sesión está activa usando el patrón State
		ctx := &sesion_estado.Sesion{
			Fecha:      s.Fecha,
			HoraInicio: s.HoraInicio,
			HoraFin:    s.HoraFin,
		}
		activa := ctx.CanRegistrarAsistencia()
		sesionesView = append(sesionesView, SesionView{
			ID:         s.ID.String(),
			Fecha:      s.Fecha,
			HoraInicio: s.HoraInicio,
			HoraFin:    s.HoraFin,
			Activa:     activa,
		})
	}
	c.vista.RenderizarGestionarSesiones(w, map[string]interface{}{
		"Sesiones": sesionesView,
		"Exito":    true,
	})
}

func (c *SesionAsistenciaControlador) MostrarRegistrarAsistencias(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Obtener sesión
	sesion, err := c.modelo.ObtenerSesionAsistencia(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Verificar que la sesión esté activa usando el patrón State
	ctx := &sesion_estado.Sesion{
		Fecha:      sesion.Fecha,
		HoraInicio: sesion.HoraInicio,
		HoraFin:    sesion.HoraFin,
	}
	if !ctx.CanRegistrarAsistencia() {
		http.Error(w, "Solo se pueden registrar asistencias en sesiones activas", http.StatusForbidden)
		return
	}

	// Obtener lista de estudiantes REALES de la base de datos
	estudiantesDB, err := c.estudianteModelo.MostrarEstudiantes()
	if err != nil {
		http.Error(w, "Error al obtener estudiantes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convertir a formato para la vista
	estudiantes := []struct {
		ID        string
		Nombre    string
		Apellidos string
		Registro  string
	}{}

	for _, est := range estudiantesDB {
		estudiantes = append(estudiantes, struct {
			ID        string
			Nombre    string
			Apellidos string
			Registro  string
		}{
			ID:        est.ID.String(),
			Nombre:    est.Nombre,
			Apellidos: est.Apellidos,
			Registro:  est.Registro,
		})
	}

	data := map[string]interface{}{
		"Sesion":      sesion,
		"Estudiantes": estudiantes,
	}

	c.vista.RenderizarRegistrarAsistencias(w, data)
}

func (c *SesionAsistenciaControlador) ProcesarSeleccionEstudiante(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error en el formulario", http.StatusBadRequest)
		return
	}

	// Obtener parámetros
	sesionIDStr := mux.Vars(r)["id"]
	estudianteIDStr := r.FormValue("estudiante_id")

	if estudianteIDStr == "" {
		http.Error(w, "Debe seleccionar un estudiante", http.StatusBadRequest)
		return
	}

	// Redirigir a la captura de foto para reconocimiento facial
	http.Redirect(w, r, "/capturar-foto?sesion="+sesionIDStr+"&estudiante="+estudianteIDStr, http.StatusSeeOther)
}

func (c *SesionAsistenciaControlador) MostrarFormularioFoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sesionIDStr := vars["id"]
	estudianteIDStr := vars["estudiante_id"]

	sesionID, err := uuid.Parse(sesionIDStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	estudianteID, err := uuid.Parse(estudianteIDStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Obtener sesión
	sesion, err := c.modelo.ObtenerSesionAsistencia(sesionID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Obtener estudiante
	estudiante, err := c.estudianteModelo.ObtenerEstudiantePorID(estudianteID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Verificar que la sesión esté activa usando el patrón State
	ctx := &sesion_estado.Sesion{
		Fecha:      sesion.Fecha,
		HoraInicio: sesion.HoraInicio,
		HoraFin:    sesion.HoraFin,
	}
	activa := ctx.CanVerRostro()

	// Verificar que el estudiante tenga foto de referencia
	tieneFotoReferencia := estudiante.FotoReferencia != ""

	data := map[string]interface{}{
		"Sesion":              sesion,
		"Estudiante":          estudiante,
		"Activa":              activa,
		"TieneFotoReferencia": tieneFotoReferencia,
	}

	c.vista.RenderizarFormularioFoto(w, data)
}
