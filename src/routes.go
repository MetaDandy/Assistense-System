package src

import (
	"github.com/MetaDandy/Assistense-System/config"
	"github.com/MetaDandy/Assistense-System/src/controlador"
	"github.com/MetaDandy/Assistense-System/src/modelo"
	"github.com/MetaDandy/Assistense-System/src/vista"
	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	docenteModelo := modelo.NuevoDocenteModelo(config.DB)
	docenteControlador := controlador.NuevoDocenteControlador(docenteModelo)

	sesionModelo := modelo.NuevaSesionAsistenciaModelo(config.DB)
	sesionVista := vista.NuevaSesionAsistenciaVistaHTML()
	sesionControlador := controlador.NuevoSesionAsistenciaControlador(sesionModelo, sesionVista)

	// Página principal
	r.HandleFunc("/", docenteControlador.MostrarInicio).Methods("GET")

	// Formularios directamente en localhost:8000
	r.HandleFunc("/registro", docenteControlador.MostrarRegistro).Methods("GET")
	r.HandleFunc("/login", docenteControlador.MostrarLogin).Methods("GET")

	// Procesar formularios
	r.HandleFunc("/registro", docenteControlador.ProcesarRegistro).Methods("POST")
	r.HandleFunc("/login", docenteControlador.ProcesarLogin).Methods("POST")
	r.HandleFunc("/panel-docente", docenteControlador.MostrarPanelDocente).Methods("GET")

	// Rutas para sesiones de asistencia
	r.HandleFunc("/sesion-asistencia/registrar", sesionControlador.MostrarRegistrar).Methods("GET")
	r.HandleFunc("/sesion-asistencia/registrar", sesionControlador.ProcesarRegistrar).Methods("POST")
	r.HandleFunc("/sesion-asistencia/listar", sesionControlador.ListarSesiones).Methods("GET")
	r.HandleFunc("/sesion-asistencia/{id}", sesionControlador.MostrarDetalle).Methods("GET")

	return r
}
