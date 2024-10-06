package handler

import (
	"net/http"
	"strconv"

	"github.com/chall-goflutter-api/api/middleware"
	"github.com/chall-goflutter-api/internal/stand"
	"github.com/chall-goflutter-api/internal/types"
	"github.com/chall-goflutter-api/internal/user"
	"github.com/chall-goflutter-api/pkg/errors"
	"github.com/chall-goflutter-api/pkg/json"
	"github.com/gorilla/mux"
)

type StandHandler struct {
	service   stand.StandService
	userStore user.UserStore
}

func NewStandHandler(service stand.StandService, userStore user.UserStore) *StandHandler {
	return &StandHandler{
		service:   service,
		userStore: userStore,
	}
}

func (h *StandHandler) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/stands", errors.ErrorHandler(middleware.IsAuth(h.GetAll, h.userStore))).Methods(http.MethodGet)
	mux.Handle("/stands/{id}", errors.ErrorHandler(middleware.IsAuth(h.Get, h.userStore))).Methods(http.MethodGet)
	mux.Handle("/stands", errors.ErrorHandler(middleware.IsAuth(h.Create, h.userStore, types.UserRoleTeneurStand))).Methods(http.MethodPost)
	mux.Handle("/stands/{id}", errors.ErrorHandler(middleware.IsAuth(h.Update, h.userStore, types.UserRoleTeneurStand))).Methods(http.MethodPatch)
}

func (h *StandHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	stands, err := h.service.GetAll(r.Context())
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, stands); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *StandHandler) Get(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	stand, err := h.service.Get(r.Context(), id)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, stand); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *StandHandler) Create(w http.ResponseWriter, r *http.Request) error {
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

func (h *StandHandler) Update(w http.ResponseWriter, r *http.Request) error {
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
