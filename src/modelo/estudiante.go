package modelo

import (
	"fmt"

	"github.com/MetaDandy/Assistense-System/helper"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Estudiante struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Nombre         string    `gorm:"type:varchar(100);not null"`
	Apellidos      string    `gorm:"type:varchar(100);not null"`
	Registro       string    `gorm:"type:varchar(10);uniqueIndex;not null"`
	FotoReferencia string    `gorm:"type:text"`
}

type RegistrarEstudianteDto struct {
	Nombre         string `json:"nombre" binding:"required"`
	Apellidos      string `json:"apellidos" binding:"required"`
	Registro       string `json:"registro" binding:"required,max=10"`
	FotoReferencia string `json:"foto_referencia,omitempty"` // Base64 de la foto
}

type ActualizarEstudiante struct {
	Nombre         *string `json:"nombre" binding:"required"`
	Apellidos      *string `json:"apellidos" binding:"required"`
	Registro       *string `json:"registro" binding:"required,max=10"`
	FotoReferencia *string `json:"foto_referencia,omitempty"`
}

type EstudianteModeloInterfaz interface {
	RegistrarEstudiante(estudiante *RegistrarEstudianteDto) (*Estudiante, error)
	ActualizarEstudiante(id uuid.UUID, estudiante *ActualizarEstudiante) (*Estudiante, error)
	MostrarEstudiantes() ([]Estudiante, error)
	ObtenerEstudiantePorID(id uuid.UUID) (*Estudiante, error)
}

type EstudianteModelo struct {
	db *gorm.DB
}

func NuevoEstudianteModelo(db *gorm.DB) EstudianteModeloInterfaz {
	return &EstudianteModelo{db: db}
}

func (em *EstudianteModelo) RegistrarEstudiante(estudiante *RegistrarEstudianteDto) (*Estudiante, error) {
	var existe Estudiante

	if err := em.db.Where("registro = ?", estudiante.Registro).First(&existe).Error; err == nil {
		return nil, gorm.ErrRegistered
	}

	// Validar foto de referencia si se proporciona
	if estudiante.FotoReferencia != "" {
		if err := helper.ValidarImagenBase64(estudiante.FotoReferencia); err != nil {
			return nil, fmt.Errorf("foto de referencia inválida: %v", err)
		}
	}

	nuevoEstudiante := Estudiante{}
	copier.Copy(&nuevoEstudiante, estudiante)
	nuevoEstudiante.ID = uuid.New()

	if err := em.db.Create(&nuevoEstudiante).Error; err != nil {
		return nil, err
	}

	return &nuevoEstudiante, nil
}

func (em *EstudianteModelo) ActualizarEstudiante(id uuid.UUID, estudiante *ActualizarEstudiante) (*Estudiante, error) {
	var existente Estudiante

	if err := em.db.First(&existente, "id = ?", id).Error; err != nil {
		return nil, err
	}

	// Validar foto de referencia si se proporciona
	if estudiante.FotoReferencia != nil && *estudiante.FotoReferencia != "" {
		if err := helper.ValidarImagenBase64(*estudiante.FotoReferencia); err != nil {
			return nil, fmt.Errorf("foto de referencia inválida: %v", err)
		}
	}

	opt := copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}

	if err := copier.CopyWithOption(&existente, estudiante, opt); err != nil {
		return nil, fmt.Errorf("failed to update fields: %w", err)
	}
	if err := em.db.Save(&existente).Error; err != nil {
		return nil, err
	}

	return &existente, nil
}

func (em *EstudianteModelo) MostrarEstudiantes() ([]Estudiante, error) {
	var estudiantes []Estudiante

	if err := em.db.Find(&estudiantes).Error; err != nil {
		return nil, err
	}

	return estudiantes, nil
}

func (em *EstudianteModelo) ObtenerEstudiantePorID(id uuid.UUID) (*Estudiante, error) {
	var estudiante Estudiante

	if err := em.db.First(&estudiante, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &estudiante, nil
}
