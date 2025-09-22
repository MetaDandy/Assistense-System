package controlador

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/MetaDandy/Assistense-System/src/modelo"
	"github.com/MetaDandy/Assistense-System/src/vista"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type EstudianteControlador struct {
	modelos   modelo.EstudianteModeloInterfaz
	vistaHTML *vista.EstudianteVistaHTML
}

type EstudianteControladorInterfaz interface {
	MostrarGestionarEstudiantes(w http.ResponseWriter, r *http.Request)
	ProcesarRegistrarEstudiante(w http.ResponseWriter, r *http.Request)
	MostrarEditarEstudiante(w http.ResponseWriter, r *http.Request)
	ProcesarEditarEstudiante(w http.ResponseWriter, r *http.Request)
}

func NuevoEstudianteControlador(modelos modelo.EstudianteModeloInterfaz, vistas *vista.EstudianteVistaHTML) EstudianteControladorInterfaz {
	return &EstudianteControlador{
		modelos:   modelos,
		vistaHTML: vistas,
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

	var estudiante modelo.RegistrarEstudianteDto

	// Verificar si es JSON o form data
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		// Manejo de JSON
		if err := json.NewDecoder(r.Body).Decode(&estudiante); err != nil {
			http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
			return
		}
	} else {
		// Manejo de form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error procesando formulario", http.StatusBadRequest)
			return
		}
		estudiante = modelo.RegistrarEstudianteDto{
			Nombre:         r.FormValue("nombre"),
			Apellidos:      r.FormValue("apellidos"),
			Registro:       r.FormValue("registro"),
			FotoReferencia: r.FormValue("foto_referencia"),
		}
	}

	if _, err := ec.modelos.RegistrarEstudiante(&estudiante); err != nil {
		if strings.Contains(contentType, "application/json") {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "Error al registrar estudiante", http.StatusInternalServerError)
		}
		return
	}

	if strings.Contains(contentType, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Estudiante registrado exitosamente"})
	} else {
		http.Redirect(w, r, "/gestionar-estudiantes", http.StatusSeeOther)
	}
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
	contentType := r.Header.Get("Content-Type")

	var actualizar modelo.ActualizarEstudiante

	if strings.Contains(contentType, "application/json") {
		// Manejo de JSON
		var data struct {
			Nombre         string `json:"nombre"`
			Apellidos      string `json:"apellidos"`
			Registro       string `json:"registro"`
			FotoReferencia string `json:"foto_referencia"`
		}

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Error decodificando JSON", http.StatusBadRequest)
			return
		}

		actualizar = modelo.ActualizarEstudiante{
			Nombre:         &data.Nombre,
			Apellidos:      &data.Apellidos,
			Registro:       &data.Registro,
			FotoReferencia: &data.FotoReferencia,
		}
	} else {
		// Manejo de form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error procesando formulario", http.StatusBadRequest)
			return
		}

		nombre := r.FormValue("nombre")
		apellidos := r.FormValue("apellidos")
		registro := r.FormValue("registro")
		fotoReferencia := r.FormValue("foto_referencia")

		actualizar = modelo.ActualizarEstudiante{
			Nombre:         &nombre,
			Apellidos:      &apellidos,
			Registro:       &registro,
			FotoReferencia: &fotoReferencia,
		}
	}

	if _, err := ec.modelos.ActualizarEstudiante(uuid.MustParse(id), &actualizar); err != nil {
		if strings.Contains(contentType, "application/json") {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Error(w, "Error al actualizar estudiante", http.StatusInternalServerError)
		}
		return
	}

	if strings.Contains(contentType, "application/json") {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Estudiante actualizado exitosamente"})
	} else {
		http.Redirect(w, r, "/gestionar-estudiantes", http.StatusSeeOther)
	}
}
