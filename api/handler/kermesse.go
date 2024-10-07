package handler

import (
	"net/http"
	"strconv"

	"github.com/chall-goflutter-api/api/middleware"
	"github.com/chall-goflutter-api/internal/kermesse"
	"github.com/chall-goflutter-api/internal/types"
	"github.com/chall-goflutter-api/internal/user"
	"github.com/chall-goflutter-api/pkg/errors"
	"github.com/chall-goflutter-api/pkg/json"
	"github.com/gorilla/mux"
)

type KermesseHandler struct {
	service   kermesse.KermesseService
	userStore user.UserStore
}

func NewKermesseHandler(service kermesse.KermesseService, userStore user.UserStore) *KermesseHandler {
	return &KermesseHandler{
		service:   service,
		userStore: userStore,
	}
}

func (h *KermesseHandler) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/kermesses", errors.ErrorHandler(middleware.IsAuth(h.GetAll, h.userStore))).Methods(http.MethodGet)
	mux.Handle("/kermesses/{id}", errors.ErrorHandler(middleware.IsAuth(h.Get, h.userStore))).Methods(http.MethodGet)
	mux.Handle("/kermesses", errors.ErrorHandler(middleware.IsAuth(h.Create, h.userStore, types.UserRoleOrganisateur))).Methods(http.MethodPost)
	mux.Handle("/kermesses/{id}", errors.ErrorHandler(middleware.IsAuth(h.Update, h.userStore, types.UserRoleOrganisateur))).Methods(http.MethodPatch)
	mux.Handle("/kermesses/{id}/participant", errors.ErrorHandler(middleware.IsAuth(h.AddParticipant, h.userStore))).Methods(http.MethodPost)
	mux.Handle("/kermesses/{id}/stand", errors.ErrorHandler(middleware.IsAuth(h.AddStand, h.userStore))).Methods(http.MethodPost)
}

func (h *KermesseHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	kermesses, err := h.service.GetAll(r.Context())
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, kermesses); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *KermesseHandler) Get(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	kermesse, err := h.service.Get(r.Context(), id)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, kermesse); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *KermesseHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.Create(r.Context(), input); err != nil {
		return err
	}

	if err := json.Write(w, http.StatusCreated, nil); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *KermesseHandler) Update(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.Update(r.Context(), id, input); err != nil {
		return err
	}

	if err := json.Write(w, http.StatusAccepted, nil); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *KermesseHandler) AddParticipant(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}
	input["kermesse_id"] = id

	if err := h.service.AddParticipant(r.Context(), input); err != nil {
		return err
	}

	if err := json.Write(w, http.StatusCreated, nil); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *KermesseHandler) AddStand(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}
	input["kermesse_id"] = id

	if err := h.service.AddStand(r.Context(), input); err != nil {
		return err
	}

	if err := json.Write(w, http.StatusCreated, nil); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}
