package modelo

import (
	"time"

	"github.com/MetaDandy/Assistense-System/src/modelo/cadena_responsabilidad"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Asistencia struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;"`
	FechaHora        string    `gorm:"type:varchar(50);not null"`
	FotoVerificacion string    `gorm:"type:text"`
	Similitud        float64   `gorm:"type:decimal(5,4)"`

	EstudianteID       uuid.UUID `gorm:"type:uuid;not null"`
	SesionAsistenciaID uuid.UUID `gorm:"type:uuid;not null"`

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
	db               *gorm.DB
	estudianteModelo EstudianteModeloInterfaz
	sesionModelo     SesionAsistenciaInterfaz
}

func NuevoAsistenciaModelo(db *gorm.DB, estudianteModelo EstudianteModeloInterfaz, sesionModelo SesionAsistenciaInterfaz) AsistenciaInterfaz {
	return &AsistenciaModelo{
		db:               db,
		estudianteModelo: estudianteModelo,
		sesionModelo:     sesionModelo,
	}
}

func (am *AsistenciaModelo) RegistrarAsistencia(dto *RegistrarAsistenciaDto) (*Asistencia, error) {
	// Crear la solicitud que viajará por la cadena de validadores
	solicitud := &cadena_responsabilidad.SolicitudAsistencia{
		FotoVerificacion: dto.FotoVerificacion,
		SesionID:         dto.SesionAsistenciaID,
		EstudianteID:     dto.EstudianteID,
	}

	// Construir la cadena de validadores
	primerValidador := am.construirCadenaValidadores()

	// Validar usando la cadena de responsabilidad
	// Iniciar la cadena desde el primer validador (ValidadorImagen)
	if err := primerValidador.Validar(solicitud); err != nil {
		return nil, err
	}

	// Si todas las validaciones pasaron, registrar la asistencia
	asistencia := &Asistencia{
		ID:                 uuid.New(),
		FechaHora:          time.Now().Format("2006-01-02 15:04:05"),
		FotoVerificacion:   dto.FotoVerificacion,
		Similitud:          solicitud.Similitud,
		EstudianteID:       dto.EstudianteID,
		SesionAsistenciaID: dto.SesionAsistenciaID,
	}

	if err := am.db.Create(asistencia).Error; err != nil {
		return nil, err
	}

	return asistencia, nil
}

// construirCadenaValidadores construye la cadena de responsabilidad con los validadores
// Sigue el patrón: Validador → ValidadorImagen → ValidadorUUID → ValidadorEstudiante →
// ValidadorFotoReferencia → ValidadorSimilitud → ValidadorDuplicado
func (am *AsistenciaModelo) construirCadenaValidadores() cadena_responsabilidad.Validador {
	// Definir callbacks para evitar ciclos de importación

	// Callback para verificar existencia del estudiante
	callbackEstudiante := func(estudianteID uuid.UUID) (string, error) {
		_, err := am.estudianteModelo.ObtenerEstudiantePorID(estudianteID)
		if err != nil {
			return "", err
		}
		return "", nil
	}

	// Callback para obtener foto de referencia
	callbackFotoRef := func(estudianteID uuid.UUID) (string, error) {
		estudiante, err := am.estudianteModelo.ObtenerEstudiantePorID(estudianteID)
		if err != nil {
			return "", err
		}
		return estudiante.FotoReferencia, nil
	}

	// Callback para verificar asistencia duplicada
	callbackDuplicado := func(estudianteID uuid.UUID, sesionID uuid.UUID) (bool, error) {
		return am.VerificarAsistenciaExistente(estudianteID, sesionID)
	}

	// Crear instancias de cada validador usando el paquete cadena_responsabilidad
	v1 := cadena_responsabilidad.NewValidadorImagen()
	v2 := cadena_responsabilidad.NewValidadorUUID()
	v3 := cadena_responsabilidad.NewValidadorEstudiante(callbackEstudiante)
	v4 := cadena_responsabilidad.NewValidadorFotoReferencia(callbackFotoRef)
	v5 := cadena_responsabilidad.NewValidadorSimilitud(callbackFotoRef)
	v6 := cadena_responsabilidad.NewValidadorDuplicado(callbackDuplicado)

	// Encadenar los validadores: v1 → v2 → v3 → v4 → v5 → v6
	// SetSiguiente retorna el siguiente, permitiendo encadenamiento fluido
	v1.SetSiguiente(v2)
	v2.SetSiguiente(v3)
	v3.SetSiguiente(v4)
	v4.SetSiguiente(v5)
	v5.SetSiguiente(v6)

	// Retornar el primer manejador de la cadena
	return v1
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
