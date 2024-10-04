package handler

import (
	"net/http"
	"strconv"

	"github.com/chall-goflutter-api/api/middleware"
	"github.com/chall-goflutter-api/internal/types"
	"github.com/chall-goflutter-api/internal/user"
	"github.com/chall-goflutter-api/pkg/errors"
	"github.com/chall-goflutter-api/pkg/json"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	service user.UserService
	store   user.UserStore
}

func NewUserHandler(service user.UserService, store user.UserStore) *UserHandler {
	return &UserHandler{
		service: service,
		store:   store,
	}
}

func (h *UserHandler) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/users/{id}", errors.ErrorHandler(middleware.IsAuth(h.Get, h.store))).Methods(http.MethodGet)
	mux.Handle("/users/invite", errors.ErrorHandler(middleware.IsAuth(h.Invite, h.store, types.UserRoleParent))).Methods(http.MethodPost)
	mux.Handle("/users/distribute", errors.ErrorHandler(middleware.IsAuth(h.Distribute, h.store, types.UserRoleParent))).Methods(http.MethodPatch)

	mux.Handle("/register", errors.ErrorHandler(h.Register)).Methods(http.MethodPost)
	mux.Handle("/login", errors.ErrorHandler(h.Login)).Methods(http.MethodPost)
	mux.Handle("/me", errors.ErrorHandler(middleware.IsAuth(h.GetMe, h.store))).Methods(http.MethodGet)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	user, err := h.service.Get(r.Context(), id)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, user); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *UserHandler) Invite(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.Invite(r.Context(), input); err != nil {
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

func (h *UserHandler) Distribute(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.Distribute(r.Context(), input); err != nil {
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

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if err := h.service.Register(r.Context(), input); err != nil {
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

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) error {
	var input map[string]interface{}
	if err := json.Parse(r, &input); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	response, err := h.service.Login(r.Context(), input)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, response); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) error {
	response, err := h.service.GetMe(r.Context())
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, response); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}
