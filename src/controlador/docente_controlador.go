package controlador

import (
	"encoding/json"
	"net/http"

	"github.com/MetaDandy/Assistense-System/helper"
	"github.com/MetaDandy/Assistense-System/src/modelo"
	"github.com/MetaDandy/Assistense-System/src/vista"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type DocenteControlador struct {
	modelos modelo.DocenteInterfaz
}

type DocenteControladorInterfaz interface {
	RegistrarDocente(w http.ResponseWriter, r *http.Request)
	IniciarSesion(w http.ResponseWriter, r *http.Request)
	ObtenerDocentePorID(w http.ResponseWriter, r *http.Request)
	ActualizarDocente(w http.ResponseWriter, r *http.Request)
}

func NuevoDocenteControlador(modelos modelo.DocenteInterfaz) DocenteControladorInterfaz {
	return &DocenteControlador{modelos: modelos}
}

func (dc *DocenteControlador) RegistrarDocente(w http.ResponseWriter, r *http.Request) {
	var registro modelo.RegistrarDocenteDto

	if err := json.NewDecoder(r.Body).Decode(&registro); err != nil {
		helper.EnviarJson(w, http.StatusBadRequest, map[string]string{
			"error": "JSON inválido",
		})
		return
	}

	docente, token, err := dc.modelos.RegistrarDocente(&registro)
	if err != nil {
		helper.EnviarJson(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	respuesta := vista.DocenteARespuesta(docente)
	respuesta.Token = token

	helper.EnviarJson(w, http.StatusCreated, respuesta)
}

func (dc *DocenteControlador) IniciarSesion(w http.ResponseWriter, r *http.Request) {
	var login modelo.IniciarSesionDto

	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		helper.EnviarJson(w, http.StatusBadRequest, map[string]string{
			"error": "JSON inválido",
		})
		return
	}

	docente, token, err := dc.modelos.IniciarSesion(login)
	if err != nil {
		helper.EnviarJson(w, http.StatusUnauthorized, map[string]string{
			"error": err.Error(),
		})
		return
	}

	respuesta := vista.DocenteARespuesta(docente)
	respuesta.Token = token

	helper.EnviarJson(w, http.StatusOK, respuesta)
}

func (dc *DocenteControlador) ObtenerDocentePorID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	docenteID, err := uuid.Parse(id)
	if err != nil {
		helper.EnviarJson(w, http.StatusBadRequest, map[string]string{
			"error": "ID inválido",
		})
		return
	}

	docente, err := dc.modelos.ObtenerDocentePorID(docenteID)
	if err != nil {
		helper.EnviarJson(w, http.StatusNotFound, map[string]string{
			"error": err.Error(),
		})
		return
	}

	respuesta := vista.DocenteARespuesta(docente)
	helper.EnviarJson(w, http.StatusOK, respuesta)
}

func (dc *DocenteControlador) ActualizarDocente(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var actualizar modelo.ActualizarDocente
	if err := json.NewDecoder(r.Body).Decode(&actualizar); err != nil {
		helper.EnviarJson(w, http.StatusBadRequest, map[string]string{
			"error": "JSON inválido",
		})
		return
	}

	docenteID, err := uuid.Parse(id)
	if err != nil {
		helper.EnviarJson(w, http.StatusBadRequest, map[string]string{
			"error": "ID inválido",
		})
		return
	}

	docente, err := dc.modelos.ActualizarDocente(docenteID, &actualizar)
	if err != nil {
		helper.EnviarJson(w, http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	respuesta := vista.DocenteARespuesta(docente)
	helper.EnviarJson(w, http.StatusOK, respuesta)
}
