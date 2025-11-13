package cadena_responsabilidad

import (
	"fmt"
)

// ValidadorDuplicado valida que no exista una asistencia duplicada para la misma sesión y estudiante
// Este es típicamente el último validador en la cadena
type ValidadorDuplicado struct {
	siguiente          Validador
	permitirDuplicados bool
	verificarDuplicado CallbackVerificarDuplicado
}

// NewValidadorDuplicado crea una nueva instancia de ValidadorDuplicado
// Recibe un callback para verificar si ya existe una asistencia registrada
func NewValidadorDuplicado(callback CallbackVerificarDuplicado) *ValidadorDuplicado {
	return &ValidadorDuplicado{
		permitirDuplicados: false,
		verificarDuplicado: callback,
	}
}

// SetSiguiente establece el siguiente validador en la cadena
func (v *ValidadorDuplicado) SetSiguiente(validador Validador) Validador {
	v.siguiente = validador
	return validador
}

// Validar implementa la validación de asistencia duplicada
// Es típicamente el último validador en la cadena
func (v *ValidadorDuplicado) Validar(solicitud *SolicitudAsistencia) error {
	if !v.permitirDuplicados {
		existe, err := v.verificarDuplicado(solicitud.EstudianteID, solicitud.SesionID)
		if err != nil {
			return fmt.Errorf("error al verificar asistencia duplicada: %v", err)
		}
		if existe {
			return fmt.Errorf("ya existe asistencia registrada para esta sesión")
		}
	}

	// Validación exitosa, pasar al siguiente validador (si existe)
	if v.siguiente != nil {
		return v.siguiente.Validar(solicitud)
	}

	// Fin de la cadena - todas las validaciones pasaron
	return nil
}
