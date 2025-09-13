package config

import (
	"log"

	"github.com/MetaDandy/Assistense-System/src/modelo"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	log.Println("Starting migration...")

	// Primero aplicar AutoMigrate para crear/actualizar tablas
	if err := db.AutoMigrate(
		&modelo.Docente{},
		&modelo.Estudiante{},
		&modelo.SesionAsistencia{},
		&modelo.Asistencia{},
	); err != nil {
		log.Fatal("Failed to migrate database: " + err.Error())
	}

	log.Println("Migration completed")
}
