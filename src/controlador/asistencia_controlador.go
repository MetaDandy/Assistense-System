package controlador

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/MetaDandy/Assistense-System/helper"
	"github.com/MetaDandy/Assistense-System/src/modelo"
	"github.com/MetaDandy/Assistense-System/src/vista"
	"github.com/google/uuid"
)

type AsistenciaControlador struct {
	modelo           modelo.AsistenciaInterfaz
	estudianteModelo modelo.EstudianteInterfaz
	vista            *vista.AsistenciaVistaHTML
}

func NuevoAsistenciaControlador(m modelo.AsistenciaInterfaz, em modelo.EstudianteInterfaz, v *vista.AsistenciaVistaHTML) *AsistenciaControlador {
	return &AsistenciaControlador{modelo: m, estudianteModelo: em, vista: v}
}

// GET /asistencia/confirmar?sesion=uuid&estudiante=uuid
func (c *AsistenciaControlador) MostrarConfirmarAsistencia(w http.ResponseWriter, r *http.Request) {
	sesionID := r.URL.Query().Get("sesion")
	estudianteID := r.URL.Query().Get("estudiante")

	// Validar parámetros
	if sesionID == "" || estudianteID == "" {
		http.Error(w, "Parámetros inválidos", http.StatusBadRequest)
		return
	}

	// Validar UUIDs
	_, err := uuid.Parse(sesionID)
	if err != nil {
		http.Error(w, "ID de sesión inválido", http.StatusBadRequest)
		return
	}

	_, err = uuid.Parse(estudianteID)
	if err != nil {
		http.Error(w, "ID de estudiante inválido", http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"SesionID":     sesionID,
		"EstudianteID": estudianteID,
	}

	c.vista.RenderizarConfirmarAsistencia(w, data)
}

// GET /capturar-foto?sesion=uuid&estudiante=uuid
func (c *AsistenciaControlador) MostrarCapturarFoto(w http.ResponseWriter, r *http.Request) {
	sesionID := r.URL.Query().Get("sesion")
	estudianteID := r.URL.Query().Get("estudiante")

	// Validar parámetros
	if sesionID == "" || estudianteID == "" {
		http.Error(w, "Parámetros inválidos", http.StatusBadRequest)
		return
	}

	// Validar UUIDs
	_, err := uuid.Parse(sesionID)
	if err != nil {
		http.Error(w, "ID de sesión inválido", http.StatusBadRequest)
		return
	}

	estudianteUUID, err := uuid.Parse(estudianteID)
	if err != nil {
		http.Error(w, "ID de estudiante inválido", http.StatusBadRequest)
		return
	}

	// Obtener estudiante
	estudiante, err := c.estudianteModelo.ObtenerEstudiantePorID(estudianteUUID)
	if err != nil {
		http.Error(w, "Estudiante no encontrado", http.StatusNotFound)
		return
	}

	// Verificar que el estudiante tenga foto de referencia
	if estudiante.FotoReferencia == "" {
		http.Error(w, "El estudiante no tiene foto de referencia registrada", http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"SesionID":     sesionID,
		"EstudianteID": estudianteID,
		"Estudiante":   estudiante,
	}

	c.vista.RenderizarCapturarFoto(w, data)
}

// POST /api/registrar-asistencia
func (c *AsistenciaControlador) ProcesarRegistrarAsistencia(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método no permitido"})
		return
	}

	// Parsear JSON del request
	var request struct {
		FotoVerificacion string `json:"foto_verificacion"`
		SesionID         string `json:"sesion_id"`
		EstudianteID     string `json:"estudiante_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al procesar datos"})
		return
	}

	// Validar datos requeridos
	if request.FotoVerificacion == "" || request.SesionID == "" || request.EstudianteID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Datos requeridos faltantes"})
		return
	}

	// Validar formato de la imagen
	if err := helper.ValidarImagenBase64(request.FotoVerificacion); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Formato de imagen inválido: " + err.Error()})
		return
	}

	// Validar UUIDs
	sesionUUID, err := uuid.Parse(request.SesionID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID de sesión inválido"})
		return
	}

	estudianteUUID, err := uuid.Parse(request.EstudianteID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID de estudiante inválido"})
		return
	}

	// Obtener la foto de referencia del estudiante
	estudiante, err := c.estudianteModelo.ObtenerEstudiantePorID(estudianteUUID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Estudiante no encontrado"})
		return
	}

	if estudiante.FotoReferencia == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "El estudiante no tiene foto de referencia registrada"})
		return
	}

	// Comparar rostros
	esIgual, similitud, err := helper.CompararRostros(estudiante.FotoReferencia, request.FotoVerificacion)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al procesar reconocimiento facial: " + err.Error()})
		return
	}

	// Verificar si es la misma persona (usando umbral más permisivo)
	if !esIgual || similitud < 0.6 {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": fmt.Sprintf("El rostro no coincide con el registrado (similitud: %.1f%%)", similitud*100),
		})
		return
	}

	// Crear DTO para registrar asistencia
	dto := &modelo.RegistrarAsistenciaDto{
		FotoVerificacion:   request.FotoVerificacion,
		Similitud:          similitud,
		EstudianteID:       estudianteUUID,
		SesionAsistenciaID: sesionUUID,
	}

	// Registrar asistencia
	asistencia, err := c.modelo.RegistrarAsistencia(dto)
	if err != nil {
		// Si es error de duplicado
		if err.Error() == "UNIQUE constraint failed" || err.Error() == "duplicate key value" {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "Ya se registró asistencia para esta sesión"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al registrar asistencia: " + err.Error()})
		return
	}

	// Respuesta exitosa
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":   true,
		"message":   "Asistencia registrada exitosamente",
		"id":        asistencia.ID.String(),
		"similitud": similitud,
	})
}
