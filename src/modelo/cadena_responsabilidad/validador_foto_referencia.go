package cadena_responsabilidad

import (
	"fmt"
)

// ValidadorFotoReferencia valida que el estudiante tenga una foto de referencia registrada
type ValidadorFotoReferencia struct {
	siguiente   Validador
	obtenerFoto CallbackObtenerFotoReferencia
}

// NewValidadorFotoReferencia crea una nueva instancia de ValidadorFotoReferencia
// Recibe un callback para obtener la foto de referencia del estudiante
func NewValidadorFotoReferencia(callback CallbackObtenerFotoReferencia) *ValidadorFotoReferencia {
	return &ValidadorFotoReferencia{
		obtenerFoto: callback,
	}
}

// SetSiguiente establece el siguiente validador en la cadena
func (v *ValidadorFotoReferencia) SetSiguiente(validador Validador) Validador {
	v.siguiente = validador
	return validador
}

// Validar implementa la validación de foto de referencia
// Verifica que exista una foto de referencia, luego delega al siguiente validador
func (v *ValidadorFotoReferencia) Validar(solicitud *SolicitudAsistencia) error {
	foto, err := v.obtenerFoto(solicitud.EstudianteID)
	if err != nil {
		return fmt.Errorf("error al obtener foto de referencia: %v", err)
	}
	if foto == "" {
		return fmt.Errorf("estudiante no tiene foto de referencia registrada")
	}

	// Validación exitosa, pasar al siguiente validador
	if v.siguiente != nil {
		return v.siguiente.Validar(solicitud)
	}

	// Fin de la cadena
	return nil
}
