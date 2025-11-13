package cadena_responsabilidad

import (
	"fmt"
)

// ValidadorEstudiante valida que el estudiante exista en la base de datos
type ValidadorEstudiante struct {
	siguiente           Validador
	verificarEstudiante CallbackVerificarEstudiante
}

// NewValidadorEstudiante crea una nueva instancia de ValidadorEstudiante
// Recibe un callback para verificar si el estudiante existe
func NewValidadorEstudiante(callback CallbackVerificarEstudiante) *ValidadorEstudiante {
	return &ValidadorEstudiante{
		verificarEstudiante: callback,
	}
}

// SetSiguiente establece el siguiente validador en la cadena
func (v *ValidadorEstudiante) SetSiguiente(validador Validador) Validador {
	v.siguiente = validador
	return validador
}

// Validar implementa la validación de existencia del estudiante
// Verifica que el estudiante exista, luego delega al siguiente validador
func (v *ValidadorEstudiante) Validar(solicitud *SolicitudAsistencia) error {
	_, err := v.verificarEstudiante(solicitud.EstudianteID)
	if err != nil {
		return fmt.Errorf("estudiante no encontrado: %v", err)
	}

	// Validación exitosa, pasar al siguiente validador
	if v.siguiente != nil {
		return v.siguiente.Validar(solicitud)
	}

	// Fin de la cadena
	return nil
}
