package config

import (
	"log"

	"github.com/MetaDandy/Assistense-System/src/modelo"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	log.Println("Starting migration...")

	if err := db.AutoMigrate(
		&modelo.Docente{},
		&modelo.Estudiante{},
		&modelo.SesionAsistencia{},
	); err != nil {
		log.Fatal("Failed to migrate database: " + err.Error())
	}
	log.Println("AutoMigrate completed")
}
