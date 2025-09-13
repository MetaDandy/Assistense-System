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
	modelo           modelo.SesionAsistenciaInterfaz
	asistenciaModelo modelo.AsistenciaInterfaz
	estudianteModelo modelo.EstudianteInterfaz
	vista            *vista.SesionAsistenciaVistaHTML
}

func NuevoSesionAsistenciaControlador(m modelo.SesionAsistenciaInterfaz, am modelo.AsistenciaInterfaz, em modelo.EstudianteInterfaz, v *vista.SesionAsistenciaVistaHTML) *SesionAsistenciaControlador {
	return &SesionAsistenciaControlador{
		modelo:           m,
		asistenciaModelo: am,
		estudianteModelo: em,
		vista:            v,
	}
}

// esSesionActiva verifica si una sesión está activa considerando fecha y hora
func esSesionActiva(fechaSesion, horaInicio, horaFin string) bool {
	now := time.Now()
	fechaActual := now.Format("2006-01-02")
	horaActual := now.Format("15:04")

	// Solo está activa si es hoy Y está dentro del rango de horas
	return fechaSesion == fechaActual && horaActual >= horaInicio && horaActual <= horaFin
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

	for _, s := range sesiones {
		// Verificar si la sesión está activa considerando fecha y hora
		activa := esSesionActiva(s.Fecha, s.HoraInicio, s.HoraFin)

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

	activa := esSesionActiva(sesion.Fecha, sesion.HoraInicio, sesion.HoraFin)

	data := map[string]interface{}{
		"Sesion": sesion,
		"Activa": activa,
	}

	c.vista.RenderizarDetalle(w, data)
}

// GET /gestionar-sesiones - Mostrar formulario y lista en una sola vista
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
		// Verificar si la sesión está activa considerando fecha y hora
		activa := esSesionActiva(s.Fecha, s.HoraInicio, s.HoraFin)

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

// POST /gestionar-sesiones - Procesar registro de nueva sesión
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

// Función auxiliar para mostrar la vista con error
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
		activa := esSesionActiva(s.Fecha, s.HoraInicio, s.HoraFin)
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

// Función auxiliar para mostrar la vista con éxito
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
		activa := esSesionActiva(s.Fecha, s.HoraInicio, s.HoraFin)
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

// GET /sesion-asistencia/{id}/registrar - Mostrar SELECTOR de estudiantes
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

	// Verificar que la sesión esté activa
	activa := esSesionActiva(sesion.Fecha, sesion.HoraInicio, sesion.HoraFin)
	if !activa {
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

// POST /sesion-asistencia/{id}/registrar - Procesar selección de estudiante
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

// GET /sesion-asistencia/{id}/estudiante/{estudiante_id}/foto - Mostrar formulario de captura de foto
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

	// Verificar que la sesión esté activa
	activa := esSesionActiva(sesion.Fecha, sesion.HoraInicio, sesion.HoraFin)

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

// GET /sesion-asistencia/{id}/listar - Mostrar lista de asistencias
func (c *SesionAsistenciaControlador) MostrarListarAsistencias(w http.ResponseWriter, r *http.Request) {
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

	// Verificar que el docente tenga acceso a esta sesión
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

	// Verificar que la sesión pertenece al docente
	if sesion.DocenteID != docenteID {
		http.Error(w, "No tiene acceso a esta sesión", http.StatusForbidden)
		return
	}

	// Obtener asistencias reales de la base de datos
	asistenciasReales, err := c.asistenciaModelo.ObtenerAsistenciasPorSesion(id)
	if err != nil {
		http.Error(w, "Error al obtener asistencias: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Convertir a formato para la vista
	asistencias := []struct {
		ID               string
		EstudianteNombre string
		FechaHora        string
		Similitud        float64
		FotoVerificacion string
	}{}

	for _, a := range asistenciasReales {
		estudianteNombre := "Estudiante Desconocido"
		if a.Estudiante.Nombre != "" {
			estudianteNombre = a.Estudiante.Nombre + " " + a.Estudiante.Apellidos
		}

		asistencias = append(asistencias, struct {
			ID               string
			EstudianteNombre string
			FechaHora        string
			Similitud        float64
			FotoVerificacion string
		}{
			ID:               a.ID.String(),
			EstudianteNombre: estudianteNombre,
			FechaHora:        a.FechaHora,
			Similitud:        a.Similitud * 100, // Convertir a porcentaje
			FotoVerificacion: a.FotoVerificacion,
		})
	}

	data := map[string]interface{}{
		"Sesion":           sesion,
		"Asistencias":      asistencias,
		"TotalAsistencias": len(asistencias),
	}

	c.vista.RenderizarListarAsistencias(w, data)
}
