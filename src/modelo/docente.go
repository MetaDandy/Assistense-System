package modelo

import (
	"fmt"
	"time"

	"github.com/MetaDandy/Assistense-System/helper"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Docente struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Correo     string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Nombre     string    `gorm:"type:varchar(100);not null"`
	Apellidos  string    `gorm:"type:varchar(100);not null"`
	Contraseña string    `gorm:"not null"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegistrarDocenteDto struct {
	Correo     string `json:"correo" binding:"required,email"`
	Nombre     string `json:"nombre" binding:"required"`
	Apellidos  string `json:"apellidos" binding:"required"`
	Contraseña string `json:"contraseña" binding:"required,min=6"`
}

type IniciarSesionDto struct {
	Correo     string `json:"correo" binding:"required,email"`
	Contraseña string `json:"contraseña" binding:"required"`
}

type ActualizarDocente struct {
	Nombre    *string `json:"nombre" binding:"required"`
	Apellidos *string `json:"apellidos" binding:"required"`
}

type DocenteModelo struct {
	db *gorm.DB
}

type DocenteInterfaz interface {
	RegistrarDocente(docente *RegistrarDocenteDto) (*Docente, string, error)
	IniciarSesion(inicio IniciarSesionDto) (*Docente, string, error)
	ObtenerDocentePorID(id uuid.UUID) (*Docente, error)
	ActualizarDocente(id uuid.UUID, docente *ActualizarDocente) (*Docente, error)
}

func NuevoDocenteModelo(db *gorm.DB) DocenteInterfaz {
	return &DocenteModelo{db: db}
}

func (dm *DocenteModelo) RegistrarDocente(docente *RegistrarDocenteDto) (*Docente, string, error) {
	var existe Docente

	if err := dm.db.Where("correo = ?", docente.Correo).First(&existe).Error; err == nil {
		return nil, "", fmt.Errorf("el correo ya está registrado")
	}

	hash, err := helper.HashPassword(docente.Contraseña)
	if err != nil {
		return nil, "", fmt.Errorf("error al generar el hash de la contraseña")
	}

	nuevoDocente := Docente{
		ID:         uuid.New(),
		Correo:     docente.Correo,
		Nombre:     docente.Nombre,
		Apellidos:  docente.Apellidos,
		Contraseña: hash,
	}

	if err := dm.db.Create(&nuevoDocente).Error; err != nil {
		return nil, "", fmt.Errorf("error al registrar el docente")
	}

	token, err := helper.GenerateJwt(nuevoDocente.ID.String(), nuevoDocente.Correo)
	if err != nil {
		return nil, "", fmt.Errorf("error al generar el token JWT")
	}

	return &nuevoDocente, token, nil
}

func (dm *DocenteModelo) IniciarSesion(inicio IniciarSesionDto) (*Docente, string, error) {
	var docente Docente

	if err := dm.db.Where("correo = ?", inicio.Correo).First(&docente).Error; err != nil {
		return nil, "", fmt.Errorf("credenciales inválidas")
	}

	if !helper.CheckPasswordHash(inicio.Contraseña, docente.Contraseña) {
		return nil, "", fmt.Errorf("credenciales inválidas")
	}

	token, err := helper.GenerateJwt(docente.ID.String(), docente.Correo)
	if err != nil {
		return nil, "", fmt.Errorf("error al generar el token JWT")
	}

	return &docente, token, nil
}

func (dm *DocenteModelo) ObtenerDocentePorID(id uuid.UUID) (*Docente, error) {
	var docente Docente

	if err := dm.db.First(&docente, "id = ?", id).Error; err != nil {
		return nil, fmt.Errorf("docente no encontrado")
	}

	return &docente, nil
}

func (dm *DocenteModelo) ActualizarDocente(id uuid.UUID, docente *ActualizarDocente) (*Docente, error) {
	docenteExistente, err := dm.ObtenerDocentePorID(id)
	if err != nil {
		return nil, err
	}

	opt := copier.Option{
		IgnoreEmpty: true,
		DeepCopy:    true,
	}

	if err := copier.CopyWithOption(docenteExistente, docente, opt); err != nil {
		return nil, fmt.Errorf("failed to update fields: %w", err)
	}

	if err := dm.db.Save(docenteExistente).Error; err != nil {
		return nil, fmt.Errorf("error al actualizar el docente")
	}

	return docenteExistente, nil
}
