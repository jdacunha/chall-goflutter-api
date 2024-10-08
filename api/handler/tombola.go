package handler

import (
	"net/http"
	"strconv"

	"github.com/chall-goflutter-api/api/middleware"
	"github.com/chall-goflutter-api/internal/tombola"
	"github.com/chall-goflutter-api/internal/types"
	"github.com/chall-goflutter-api/internal/user"
	"github.com/chall-goflutter-api/pkg/errors"
	"github.com/chall-goflutter-api/pkg/json"
	"github.com/chall-goflutter-api/pkg/utils"
	"github.com/gorilla/mux"
)

type TombolaHandler struct {
	service   tombola.TombolaService
	userStore user.UserStore
}

func NewTombolaHandler(service tombola.TombolaService, userStore user.UserStore) *TombolaHandler {
	return &TombolaHandler{
		service:   service,
		userStore: userStore,
	}
}

func (h *TombolaHandler) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/tombolas", errors.ErrorHandler(middleware.IsAuth(h.GetAll, h.userStore))).Methods(http.MethodGet)
	mux.Handle("/tombolas/{id}", errors.ErrorHandler(middleware.IsAuth(h.Get, h.userStore))).Methods(http.MethodGet)
	mux.Handle("/tombolas", errors.ErrorHandler(middleware.IsAuth(h.Create, h.userStore, types.UserRoleOrganisateur))).Methods(http.MethodPost)
	mux.Handle("/tombolas/{id}", errors.ErrorHandler(middleware.IsAuth(h.Update, h.userStore, types.UserRoleOrganisateur))).Methods(http.MethodPatch)
	mux.Handle("/tombolas/{id}/end", errors.ErrorHandler(middleware.IsAuth(h.End, h.userStore, types.UserRoleOrganisateur))).Methods(http.MethodPatch)
}

func (h *TombolaHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	tombolas, err := h.service.GetAll(r.Context(), utils.GetQueryParams(r))
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, tombolas); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *TombolaHandler) Get(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	tombola, err := h.service.Get(r.Context(), id)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, tombola); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *TombolaHandler) Create(w http.ResponseWriter, r *http.Request) error {
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

func (h *TombolaHandler) Update(w http.ResponseWriter, r *http.Request) error {
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

func (h *TombolaHandler) End(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.End(r.Context(), id); err != nil {
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
