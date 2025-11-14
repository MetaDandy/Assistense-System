package template_method

import (
	"errors"
	"fmt"

	"github.com/MetaDandy/Assistense-System/helper"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// ProcesadorActualizar implementa ProcesadorEstudiante para ACTUALIZAR estudiantes
type ProcesadorActualizar struct {
	db               *gorm.DB
	estudianteID     string
	datosEntrada     *ActualizarEstudianteDto
	estudianteResult *Estudiante
}

// NewProcesadorActualizar crea una instancia
func NewProcesadorActualizar(db *gorm.DB, estudianteID string, datos *ActualizarEstudianteDto) *ProcesadorActualizar {
	return &ProcesadorActualizar{
		db:               db,
		estudianteID:     estudianteID,
		datosEntrada:     datos,
		estudianteResult: nil,
	}
}

// ValidarEntrada: Para ACTUALIZAR, verificar que al menos un campo sea != nil
func (a *ProcesadorActualizar) ValidarEntrada() error {
	if a.datosEntrada == nil {
		return errors.New("datos de entrada nulos")
	}

	if a.datosEntrada.Nombre == nil &&
		a.datosEntrada.Apellidos == nil &&
		a.datosEntrada.Registro == nil &&
		a.datosEntrada.FotoReferencia == nil {
		return errors.New("debe proporcionar al menos un campo para actualizar")
	}

	return nil
}

// VerificarPrecondicion: Para ACTUALIZAR, verificar que SÍ existe
func (a *ProcesadorActualizar) VerificarPrecondicion() error {
	a.estudianteResult = &Estudiante{}
	result := a.db.First(a.estudianteResult, a.estudianteID)

	if result.Error == gorm.ErrRecordNotFound {
		return errors.New("estudiante no encontrado")
	}

	if result.Error != nil {
		return result.Error
	}

	return nil
}

// PrepararEstudiante: Copiar datos del DTO al estudiante existente
func (a *ProcesadorActualizar) PrepararEstudiante() error {
	if a.estudianteResult == nil {
		return errors.New("estudiante no encontrado para actualizar")
	}

	if err := copier.CopyWithOption(a.estudianteResult, a.datosEntrada, copier.Option{
		IgnoreEmpty: true,
	}); err != nil {
		return fmt.Errorf("error al copiar datos: %w", err)
	}

	return nil
}

// ValidarFotoReferencia: Validar que la foto sea base64 válido (si se actualiza)
func (a *ProcesadorActualizar) ValidarFotoReferencia() error {
	if a.estudianteResult == nil {
		return errors.New("estudiante no preparado")
	}

	// Solo validar si la foto se actualiza
	if a.datosEntrada.FotoReferencia != nil && *a.datosEntrada.FotoReferencia != "" {
		return helper.ValidarImagenBase64(*a.datosEntrada.FotoReferencia)
	}

	return nil
}

// GuardarEnBD: Actualizar estudiante en BD (SAVE para ACTUALIZAR)
func (a *ProcesadorActualizar) GuardarEnBD() error {
	if a.estudianteResult == nil {
		return errors.New("estudiante no preparado")
	}
	return a.db.Save(a.estudianteResult).Error
}

// ObtenerResultado: Retornar el estudiante actualizado
func (a *ProcesadorActualizar) ObtenerResultado() interface{} {
	return a.estudianteResult
}
