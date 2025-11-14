package template_method

import "github.com/google/uuid"

// Estudiante es el modelo de la base de datos
type Estudiante struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Nombre         string    `gorm:"type:varchar(100);not null"`
	Apellidos      string    `gorm:"type:varchar(100);not null"`
	Registro       string    `gorm:"type:varchar(10);uniqueIndex;not null"`
	FotoReferencia string    `gorm:"type:text"`
}

// RegistrarEstudianteDto DTO para registrar un estudiante
type RegistrarEstudianteDto struct {
	Nombre         string `json:"nombre" binding:"required"`
	Apellidos      string `json:"apellidos" binding:"required"`
	Registro       string `json:"registro" binding:"required,max=10"`
	FotoReferencia string `json:"foto_referencia,omitempty"` // Base64
}

// ActualizarEstudianteDto DTO para actualizar un estudiante
type ActualizarEstudianteDto struct {
	Nombre         *string `json:"nombre" binding:"required"`
	Apellidos      *string `json:"apellidos" binding:"required"`
	Registro       *string `json:"registro" binding:"required,max=10"`
	FotoReferencia *string `json:"foto_referencia,omitempty"`
}
