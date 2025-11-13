package cadena_responsabilidad

import (
	"fmt"

	"github.com/MetaDandy/Assistense-System/helper"
)

// ValidadorSimilitud valida la similitud entre rostros usando comparación de histogramas
type ValidadorSimilitud struct {
	siguiente    Validador
	similitudMin float64
	obtenerFoto  CallbackObtenerFotoReferencia
}

// NewValidadorSimilitud crea una nueva instancia de ValidadorSimilitud
// Recibe un callback para obtener la foto de referencia
// La comparación de rostros se realiza directamente con el helper
func NewValidadorSimilitud(callbackFoto CallbackObtenerFotoReferencia) *ValidadorSimilitud {
	return &ValidadorSimilitud{
		similitudMin: 0.6,
		obtenerFoto:  callbackFoto,
	}
}

// SetSiguiente establece el siguiente validador en la cadena
func (v *ValidadorSimilitud) SetSiguiente(validador Validador) Validador {
	v.siguiente = validador
	return validador
}

// Validar implementa la validación de similitud de rostro
// Compara la foto de referencia con la foto de verificación usando el helper, luego delega al siguiente
func (v *ValidadorSimilitud) Validar(solicitud *SolicitudAsistencia) error {
	fotoReferencia, err := v.obtenerFoto(solicitud.EstudianteID)
	if err != nil {
		return fmt.Errorf("error al obtener foto de referencia para validación: %v", err)
	}

	// Comparar rostros directamente con el helper
	_, similitud, err := helper.CompararRostros(fotoReferencia, solicitud.FotoVerificacion)
	if err != nil {
		return fmt.Errorf("error al comparar rostros: %v", err)
	}

	if similitud < v.similitudMin {
		return fmt.Errorf("rostro no coincide (similitud: %.2f%% < %.2f%% requerido)", similitud*100, v.similitudMin*100)
	}

	// Almacenar la similitud en la solicitud para uso posterior
	solicitud.Similitud = similitud

	// Validación exitosa, pasar al siguiente validador
	if v.siguiente != nil {
		return v.siguiente.Validar(solicitud)
	}

	// Fin de la cadena
	return nil
}
