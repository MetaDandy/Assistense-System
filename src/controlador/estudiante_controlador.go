package controlador

import (
	"net/http"

	"github.com/MetaDandy/Assistense-System/src/modelo"
	"github.com/MetaDandy/Assistense-System/src/vista"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type EstudianteControlador struct {
	modelos   modelo.EstudianteInterfaz
	vistaHTML *vista.EstudianteVistaHTML
}

type EstudianteControladorInterfaz interface {
	MostrarGestionarEstudiantes(w http.ResponseWriter, r *http.Request)
	ProcesarRegistrarEstudiante(w http.ResponseWriter, r *http.Request)
	MostrarEditarEstudiante(w http.ResponseWriter, r *http.Request)
	ProcesarEditarEstudiante(w http.ResponseWriter, r *http.Request)
}

func NuevoEstudianteControlador(modelos modelo.EstudianteInterfaz) EstudianteControladorInterfaz {
	return &EstudianteControlador{
		modelos:   modelos,
		vistaHTML: vista.NewEstudianteVistaHTML(),
	}
}

func (ec *EstudianteControlador) MostrarGestionarEstudiantes(w http.ResponseWriter, r *http.Request) {
	estudiantes, err := ec.modelos.MostrarEstudiantes()
	if err != nil {
		http.Error(w, "Error al obtener la lista de estudiantes", http.StatusInternalServerError)
		return
	}
	ec.vistaHTML.RenderizarGestionarEstudiantes(w, map[string]interface{}{
		"Estudiantes": estudiantes,
	})
}

func (ec *EstudianteControlador) ProcesarRegistrarEstudiante(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/gestionar-estudiantes", http.StatusSeeOther)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error procesando formulario", http.StatusBadRequest)
		return
	}

	estudiante := modelo.RegistrarEstudianteDto{
		Nombre:    r.FormValue("nombre"),
		Apellidos: r.FormValue("apellidos"),
		Registro:  r.FormValue("registro"),
	}

	if _, err := ec.modelos.RegistrarEstudiante(&estudiante); err != nil {
		http.Error(w, "Error al registrar estudiante", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/gestionar-estudiantes", http.StatusSeeOther)
}

func (ec *EstudianteControlador) MostrarEditarEstudiante(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	estudiante, err := ec.modelos.ObtenerEstudiantePorID(uuid.MustParse(id))
	if err != nil {
		http.Error(w, "Estudiante no encontrado", http.StatusNotFound)
		return
	}
	ec.vistaHTML.RenderizarEditarEstudiante(w, estudiante)
}

func (ec *EstudianteControlador) ProcesarEditarEstudiante(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/gestionar-estudiantes", http.StatusSeeOther)
		return
	}

	id := mux.Vars(r)["id"]
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error procesando formulario", http.StatusBadRequest)
		return
	}

	nombre := r.FormValue("nombre")
	apellidos := r.FormValue("apellidos")
	registro := r.FormValue("registro")

	actualizar := modelo.ActualizarEstudiante{
		Nombre:    &nombre,
		Apellidos: &apellidos,
		Registro:  &registro,
	}

	if _, err := ec.modelos.ActualizarEstudiante(uuid.MustParse(id), &actualizar); err != nil {
		http.Error(w, "Error al actualizar estudiante", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/gestionar-estudiantes", http.StatusSeeOther)
}
