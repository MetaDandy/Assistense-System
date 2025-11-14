package modelo

import (
	"github.com/MetaDandy/Assistense-System/src/modelo/template_method"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Estudiante es el alias del modelo en template_method
type Estudiante = template_method.Estudiante

// RegistrarEstudianteDto es el alias del DTO en template_method
type RegistrarEstudianteDto = template_method.RegistrarEstudianteDto

// ActualizarEstudiante es el alias del DTO en template_method
type ActualizarEstudiante = template_method.ActualizarEstudianteDto

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

func (em *EstudianteModelo) RegistrarEstudiante(datos *RegistrarEstudianteDto) (*Estudiante, error) {
	// Crear el procesador específico de registrar
	proc := template_method.NewProcesadorRegistrar(em.db, datos)

	// Crear la plantilla base
	base := template_method.NewProcesadorBase()

	// Ejecutar el template method
	if err := base.Procesar(proc); err != nil {
		return nil, err
	}

	// Retornar resultado
	return proc.ObtenerResultado().(*Estudiante), nil
}

func (em *EstudianteModelo) ActualizarEstudiante(id uuid.UUID, datos *ActualizarEstudiante) (*Estudiante, error) {
	// Convertir alias al tipo real
	dtoTemplate := (*template_method.ActualizarEstudianteDto)(datos)

	// Crear el procesador específico de actualizar
	proc := template_method.NewProcesadorActualizar(em.db, id.String(), dtoTemplate)

	// Crear la plantilla base
	base := template_method.NewProcesadorBase()

	// Ejecutar el template method
	if err := base.Procesar(proc); err != nil {
		return nil, err
	}

	// Retornar resultado
	return proc.ObtenerResultado().(*Estudiante), nil
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
