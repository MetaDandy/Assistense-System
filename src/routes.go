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

	api := r.PathPrefix("/api").Subrouter()

	docentes := api.PathPrefix("/docente").Subrouter()

	docentes.HandleFunc("/registro", docenteControlador.RegistrarDocente).Methods("POST")
	docentes.HandleFunc("/login", docenteControlador.IniciarSesion).Methods("POST")
	docentes.HandleFunc("/{id}", docenteControlador.ObtenerDocentePorID).Methods("GET")
	docentes.HandleFunc("/{id}", docenteControlador.ActualizarDocente).Methods("PUT")

	return r
}
