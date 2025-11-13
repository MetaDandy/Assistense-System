package cadena_responsabilidad

import (
	"fmt"

	"github.com/google/uuid"
)

// ValidadorUUID valida que los UUIDs sean válidos
type ValidadorUUID struct {
	siguiente Validador
	formato   string // "UUID_V4"
}

// NewValidadorUUID crea una nueva instancia de ValidadorUUID
func NewValidadorUUID() *ValidadorUUID {
	return &ValidadorUUID{
		formato: "UUID_V4",
	}
}

// SetSiguiente establece el siguiente validador en la cadena
func (v *ValidadorUUID) SetSiguiente(validador Validador) Validador {
	v.siguiente = validador
	return validador
}

// Validar implementa la validación de UUIDs
// Verifica que SesionID y EstudianteID sean válidos, luego delega al siguiente
func (v *ValidadorUUID) Validar(solicitud *SolicitudAsistencia) error {
	if solicitud.SesionID == uuid.Nil {
		return fmt.Errorf("ID de sesión inválido")
	}
	if solicitud.EstudianteID == uuid.Nil {
		return fmt.Errorf("ID de estudiante inválido")
	}

	// Validación exitosa, pasar al siguiente validador
	if v.siguiente != nil {
		return v.siguiente.Validar(solicitud)
	}

	// Fin de la cadena
	return nil
}
