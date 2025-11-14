package template_method

import (
	"errors"
	"fmt"

	"github.com/MetaDandy/Assistense-System/helper"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

// ProcesadorRegistrar implementa ProcesadorEstudiante para REGISTRAR estudiantes
type ProcesadorRegistrar struct {
	db               *gorm.DB
	datosEntrada     *RegistrarEstudianteDto
	estudianteResult *Estudiante
}

// NewProcesadorRegistrar crea una instancia
func NewProcesadorRegistrar(db *gorm.DB, datos *RegistrarEstudianteDto) *ProcesadorRegistrar {
	return &ProcesadorRegistrar{
		db:               db,
		datosEntrada:     datos,
		estudianteResult: nil,
	}
}

// ValidarEntrada: Verificar que TODOS los campos requeridos existan
func (r *ProcesadorRegistrar) ValidarEntrada() error {
	if r.datosEntrada == nil {
		return errors.New("datos de entrada nulos")
	}
	if r.datosEntrada.Nombre == "" {
		return errors.New("nombre requerido")
	}
	if r.datosEntrada.Apellidos == "" {
		return errors.New("apellidos requeridos")
	}
	if r.datosEntrada.Registro == "" {
		return errors.New("registro requerido")
	}
	if r.datosEntrada.FotoReferencia == "" {
		return errors.New("foto de referencia requerida")
	}
	return nil
}

// VerificarPrecondicion: Para REGISTRAR, verificar que NO existe
func (r *ProcesadorRegistrar) VerificarPrecondicion() error {
	var existe Estudiante
	result := r.db.Where("registro = ?", r.datosEntrada.Registro).First(&existe)

	if result.Error == nil {
		return fmt.Errorf("estudiante con registro %s ya existe", r.datosEntrada.Registro)
	}

	if result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	return nil
}

// PrepararEstudiante: Generar UUID y copiar datos
func (r *ProcesadorRegistrar) PrepararEstudiante() error {
	r.estudianteResult = &Estudiante{
		ID: uuid.New(),
	}

	if err := copier.Copy(r.estudianteResult, r.datosEntrada); err != nil {
		return fmt.Errorf("error al copiar datos: %w", err)
	}

	return nil
}

// ValidarFotoReferencia: Validar que la foto sea base64 válido
func (r *ProcesadorRegistrar) ValidarFotoReferencia() error {
	if r.estudianteResult == nil {
		return errors.New("estudiante no preparado")
	}
	if r.estudianteResult.FotoReferencia == "" {
		return nil
	}
	return helper.ValidarImagenBase64(r.estudianteResult.FotoReferencia)
}

// GuardarEnBD: Crear estudiante en BD (CREATE para REGISTRAR)
func (r *ProcesadorRegistrar) GuardarEnBD() error {
	if r.estudianteResult == nil {
		return errors.New("estudiante no preparado")
	}
	return r.db.Create(r.estudianteResult).Error
}

// ObtenerResultado: Retornar el estudiante creado
func (r *ProcesadorRegistrar) ObtenerResultado() interface{} {
	return r.estudianteResult
}
