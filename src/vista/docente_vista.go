package vista

import (
	"time"

	"github.com/MetaDandy/Assistense-System/src/modelo"
	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type DocenteRespuesta struct {
	ID        uuid.UUID `json:"id"`
	Correo    string    `json:"correo"`
	Nombre    string    `json:"nombre"`
	Apellidos string    `json:"apellidos"`
	Token     string    `json:"token,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func DocenteARespuesta(u *modelo.Docente) *DocenteRespuesta {
	var docente DocenteRespuesta

	copier.Copy(&docente, u)

	return &docente
}
