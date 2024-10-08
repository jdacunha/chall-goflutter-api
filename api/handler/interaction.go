package handler

import (
	"net/http"
	"strconv"

	"github.com/chall-goflutter-api/api/middleware"
	"github.com/chall-goflutter-api/internal/interaction"
	"github.com/chall-goflutter-api/internal/types"
	"github.com/chall-goflutter-api/internal/user"
	"github.com/chall-goflutter-api/pkg/errors"
	"github.com/chall-goflutter-api/pkg/json"
	"github.com/chall-goflutter-api/pkg/utils"
	"github.com/gorilla/mux"
)

type InteractionHandler struct {
	service   interaction.InteractionService
	userStore user.UserStore
}

func NewInteractionHandler(service interaction.InteractionService, userStore user.UserStore) *InteractionHandler {
	return &InteractionHandler{
		service:   service,
		userStore: userStore,
	}
}

func (h *InteractionHandler) RegisterRoutes(mux *mux.Router) {
	mux.Handle("/interactions", errors.ErrorHandler(middleware.IsAuth(h.GetAll, h.userStore))).Methods(http.MethodGet)
	mux.Handle("/interactions/{id}", errors.ErrorHandler(middleware.IsAuth(h.Get, h.userStore))).Methods(http.MethodGet)
	mux.Handle("/interactions", errors.ErrorHandler(middleware.IsAuth(h.Create, h.userStore, types.UserRoleParent, types.UserRoleEnfant))).Methods(http.MethodPost)
	mux.Handle("/interactions/{id}", errors.ErrorHandler(middleware.IsAuth(h.Update, h.userStore, types.UserRoleTeneurStand))).Methods(http.MethodPatch)
}

func (h *InteractionHandler) GetAll(w http.ResponseWriter, r *http.Request) error {
	interactions, err := h.service.GetAll(r.Context(), utils.GetQueryParams(r))
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, interactions); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *InteractionHandler) Get(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	interaction, err := h.service.Get(r.Context(), id)
	if err != nil {
		return err
	}

	if err := json.Write(w, http.StatusOK, interaction); err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (h *InteractionHandler) Create(w http.ResponseWriter, r *http.Request) error {
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

func (h *InteractionHandler) Update(w http.ResponseWriter, r *http.Request) error {
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
