package modelo

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SesionAsistencia struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Fecha      string    `gorm:"type:varchar(10);not null"`
	HoraInicio string    `gorm:"type:varchar(5);not null"`
	HoraFin    string    `gorm:"type:varchar(5);not null"`

	DocenteID uuid.UUID `gorm:"type:uuid;not null"`
	Docente   Docente   `gorm:"foreignKey:DocenteID"`
}

type RegistrarSesionAsistenciaDto struct {
	Fecha      string    `json:"fecha" binding:"required"`
	HoraInicio string    `json:"hora_inicio" binding:"required"`
	HoraFin    string    `json:"hora_fin" binding:"required"`
	DocenteID  uuid.UUID `json:"docente_id" binding:"required"`
}

type SesionAsistenciaInterfaz interface {
	RegistrarSesionAsistencia(dto *RegistrarSesionAsistenciaDto) (*SesionAsistencia, error)
	ObtenerSesionAsistencia(id uuid.UUID) (*SesionAsistencia, error)
	ObtenerSesionesAsistencia(DocenteID uuid.UUID) ([]SesionAsistencia, error)
}

type SesionAsistenciaModelo struct {
	db *gorm.DB
}

func NuevaSesionAsistenciaModelo(db *gorm.DB) SesionAsistenciaInterfaz {
	return &SesionAsistenciaModelo{db: db}
}

func (sam *SesionAsistenciaModelo) RegistrarSesionAsistencia(dto *RegistrarSesionAsistenciaDto) (*SesionAsistencia, error) {
	var sesion SesionAsistencia

	// Guardar directamente como strings sin conversión
	sesion.ID = uuid.New()
	sesion.Fecha = dto.Fecha
	sesion.HoraInicio = dto.HoraInicio
	sesion.HoraFin = dto.HoraFin
	sesion.DocenteID = dto.DocenteID

	if err := sam.db.Create(&sesion).Error; err != nil {
		return nil, err
	}

	return &sesion, nil
}

func (sam *SesionAsistenciaModelo) ObtenerSesionAsistencia(id uuid.UUID) (*SesionAsistencia, error) {
	var sesion SesionAsistencia

	if err := sam.db.Where("id = ?", id).First(&sesion).Error; err != nil {
		return nil, err
	}

	return &sesion, nil
}

func (sam *SesionAsistenciaModelo) ObtenerSesionesAsistencia(DocenteID uuid.UUID) ([]SesionAsistencia, error) {
	var sesiones []SesionAsistencia

	if err := sam.db.Where("docente_id = ?", DocenteID).Find(&sesiones).Error; err != nil {
		return nil, err
	}

	return sesiones, nil
}
