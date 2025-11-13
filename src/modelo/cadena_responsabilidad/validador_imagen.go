package cadena_responsabilidad

import (
	"fmt"

	"github.com/MetaDandy/Assistense-System/helper"
)

// ValidadorImagen valida que la foto de verificación sea un Base64 válido
type ValidadorImagen struct {
	siguiente Validador
	tipo      string // "BASE64_IMAGE"
}

// NewValidadorImagen crea una nueva instancia de ValidadorImagen
func NewValidadorImagen() *ValidadorImagen {
	return &ValidadorImagen{
		tipo: "BASE64_IMAGE",
	}
}

// SetSiguiente establece el siguiente validador en la cadena
func (v *ValidadorImagen) SetSiguiente(validador Validador) Validador {
	v.siguiente = validador
	return validador
}

// Validar implementa la validación de imagen Base64
// Si la imagen es válida, delega al siguiente validador
func (v *ValidadorImagen) Validar(solicitud *SolicitudAsistencia) error {
	if err := helper.ValidarImagenBase64(solicitud.FotoVerificacion); err != nil {
		return fmt.Errorf("imagen base64 inválida: %v", err)
	}

	// Validación exitosa, pasar al siguiente validador
	if v.siguiente != nil {
		return v.siguiente.Validar(solicitud)
	}

	// Fin de la cadena
	return nil
}
