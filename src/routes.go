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

	estudianteModelo := modelo.NuevoEstudianteModelo(config.DB)
	estudianteControlador := controlador.NuevoEstudianteControlador(estudianteModelo)

	sesionModelo := modelo.NuevaSesionAsistenciaModelo(config.DB)
	sesionVista := vista.NuevaSesionAsistenciaVistaHTML()
	asistenciaModelo := modelo.NuevoAsistenciaModelo(config.DB)
	sesionControlador := controlador.NuevoSesionAsistenciaControlador(sesionModelo, asistenciaModelo, estudianteModelo, sesionVista)

	asistenciaVista := vista.NuevaAsistenciaVistaHTML()
	asistenciaControlador := controlador.NuevoAsistenciaControlador(asistenciaModelo, estudianteModelo, asistenciaVista)

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
	r.HandleFunc("/sesion-asistencia/{id}/registrar", sesionControlador.MostrarRegistrarAsistencias).Methods("GET")
	r.HandleFunc("/sesion-asistencia/{id}/registrar", sesionControlador.ProcesarSeleccionEstudiante).Methods("POST")
	r.HandleFunc("/sesion-asistencia/{id}/estudiante/{estudiante_id}/foto", sesionControlador.MostrarFormularioFoto).Methods("GET")
	r.HandleFunc("/sesion-asistencia/{id}/listar", sesionControlador.MostrarListarAsistencias).Methods("GET")

	// Nueva ruta para gestionar sesiones (formulario + lista en una vista)
	r.HandleFunc("/gestionar-sesiones", sesionControlador.MostrarGestionarSesiones).Methods("GET")
	r.HandleFunc("/gestionar-sesiones", sesionControlador.ProcesarGestionarSesiones).Methods("POST")

	// Rutas para gestionar estudiantes
	r.HandleFunc("/gestionar-alumnos", estudianteControlador.MostrarGestionarEstudiantes).Methods("GET")
	r.HandleFunc("/gestionar-estudiantes", estudianteControlador.MostrarGestionarEstudiantes).Methods("GET")
	r.HandleFunc("/registrar-estudiante", estudianteControlador.ProcesarRegistrarEstudiante).Methods("POST")
	r.HandleFunc("/editar-estudiante/{id}", estudianteControlador.MostrarEditarEstudiante).Methods("GET")
	r.HandleFunc("/editar-estudiante/{id}", estudianteControlador.ProcesarEditarEstudiante).Methods("POST")

	// Rutas para asistencia (escaneo de QR)
	r.HandleFunc("/asistencia/confirmar", asistenciaControlador.MostrarConfirmarAsistencia).Methods("GET")
	r.HandleFunc("/api/registrar-asistencia", asistenciaControlador.ProcesarRegistrarAsistencia).Methods("POST")

	// Ruta para captura de foto (reconocimiento facial)
	r.HandleFunc("/capturar-foto", asistenciaControlador.MostrarCapturarFoto).Methods("GET")

	return r
}
