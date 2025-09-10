package src

import (
	"github.com/MetaDandy/Assistense-System/config"
	"github.com/MetaDandy/Assistense-System/src/controlador"
	"github.com/MetaDandy/Assistense-System/src/modelo"
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	docenteModelo := modelo.NuevoDocenteModelo(config.DB)
	docenteControlador := controlador.NuevoDocenteControlador(docenteModelo)

	// HTML Templates routes (MVC Clásico) - directamente en la raíz
	// Página principal
	r.HandleFunc("/", docenteControlador.MostrarInicio).Methods("GET")

	// Formularios directamente en localhost:8000
	r.HandleFunc("/registro", docenteControlador.MostrarRegistro).Methods("GET")
	r.HandleFunc("/login", docenteControlador.MostrarLogin).Methods("GET")

	// Procesar formularios
	r.HandleFunc("/registro", docenteControlador.ProcesarRegistro).Methods("POST")
	r.HandleFunc("/login", docenteControlador.ProcesarLogin).Methods("POST")

	return r
}
