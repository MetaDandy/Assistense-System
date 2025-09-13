package modelo

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Asistencia struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;"`
	FechaHora        string    `gorm:"type:varchar(50);not null"`
	FotoVerificacion string    `gorm:"type:text"`         // Base64 de la foto de verificación
	Similitud        float64   `gorm:"type:decimal(5,4)"` // Porcentaje de similitud facial (0.0 - 1.0)

	// Relaciones
	EstudianteID       uuid.UUID `gorm:"type:uuid;not null"`
	SesionAsistenciaID uuid.UUID `gorm:"type:uuid;not null"`

	// Referencias
	Estudiante       Estudiante       `gorm:"foreignKey:EstudianteID"`
	SesionAsistencia SesionAsistencia `gorm:"foreignKey:SesionAsistenciaID"`
}

type RegistrarAsistenciaDto struct {
	FotoVerificacion   string    `json:"foto_verificacion" binding:"required"`
	Similitud          float64   `json:"similitud"`
	EstudianteID       uuid.UUID `json:"estudiante_id" binding:"required"`
	SesionAsistenciaID uuid.UUID `json:"sesion_asistencia_id" binding:"required"`
}

type AsistenciaInterfaz interface {
	RegistrarAsistencia(dto *RegistrarAsistenciaDto) (*Asistencia, error)
	ObtenerAsistenciasPorSesion(sesionID uuid.UUID) ([]Asistencia, error)
	VerificarAsistenciaExistente(estudianteID, sesionID uuid.UUID) (bool, error)
}

type AsistenciaModelo struct {
	db *gorm.DB
}

func NuevoAsistenciaModelo(db *gorm.DB) AsistenciaInterfaz {
	return &AsistenciaModelo{db: db}
}

func (am *AsistenciaModelo) RegistrarAsistencia(dto *RegistrarAsistenciaDto) (*Asistencia, error) {
	// Verificar si ya existe asistencia para este estudiante en esta sesión
	existeAsistencia, err := am.VerificarAsistenciaExistente(dto.EstudianteID, dto.SesionAsistenciaID)
	if err != nil {
		return nil, err
	}
	if existeAsistencia {
		return nil, gorm.ErrDuplicatedKey
	}

	asistencia := &Asistencia{
		ID:                 uuid.New(),
		FechaHora:          time.Now().Format("2006-01-02 15:04:05"),
		FotoVerificacion:   dto.FotoVerificacion,
		Similitud:          dto.Similitud,
		EstudianteID:       dto.EstudianteID,
		SesionAsistenciaID: dto.SesionAsistenciaID,
	}

	if err := am.db.Create(asistencia).Error; err != nil {
		return nil, err
	}

	return asistencia, nil
}

func (am *AsistenciaModelo) ObtenerAsistenciasPorSesion(sesionID uuid.UUID) ([]Asistencia, error) {
	var asistencias []Asistencia
	err := am.db.Preload("Estudiante").Where("sesion_asistencia_id = ?", sesionID).Find(&asistencias).Error
	return asistencias, err
}

func (am *AsistenciaModelo) VerificarAsistenciaExistente(estudianteID, sesionID uuid.UUID) (bool, error) {
	var count int64
	err := am.db.Model(&Asistencia{}).Where("estudiante_id = ? AND sesion_asistencia_id = ?", estudianteID, sesionID).Count(&count).Error
	return count > 0, err
}
