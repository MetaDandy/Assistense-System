package cadena_responsabilidad

import (
	"github.com/google/uuid"
)

// SolicitudAsistencia es el objeto que viaja a través de la cadena de validadores
type SolicitudAsistencia struct {
	FotoVerificacion string
	SesionID         uuid.UUID
	EstudianteID     uuid.UUID
	Similitud        float64
}

// Validador es la interfaz que define el contrato para todos los validadores
// Cada validador implementa su propia lógica de validación y mantiene referencia al siguiente
type Validador interface {
	// SetSiguiente establece el siguiente validador en la cadena
	SetSiguiente(validador Validador) Validador

	// Validar implementa la lógica específica de validación
	// Cada validador concreto decide si procesa o delega al siguiente
	Validar(solicitud *SolicitudAsistencia) error
}

// ========== CALLBACKS PARA INYECCIÓN DE DEPENDENCIAS ==========

// CallbackVerificarEstudiante verifica si un estudiante existe por su ID
type CallbackVerificarEstudiante func(estudianteID uuid.UUID) (fotoReferencia string, err error)

// CallbackObtenerFotoReferencia obtiene la foto de referencia de un estudiante
type CallbackObtenerFotoReferencia func(estudianteID uuid.UUID) (foto string, err error)

// CallbackVerificarDuplicado verifica si ya existe asistencia registrada
type CallbackVerificarDuplicado func(estudianteID uuid.UUID, sesionID uuid.UUID) (existe bool, err error)
