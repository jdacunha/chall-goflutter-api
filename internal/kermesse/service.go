package kermesse

import (
	"context"
	"database/sql"
	goErrors "errors"

	"github.com/chall-goflutter-api/internal/types"
	"github.com/chall-goflutter-api/internal/user"
	"github.com/chall-goflutter-api/pkg/errors"
	"github.com/chall-goflutter-api/pkg/utils"
)

type KermesseService interface {
	GetAll(ctx context.Context) ([]types.Kermesse, error)
	GetUsersInvite(ctx context.Context, id int) ([]types.UserBasic, error)
	Get(ctx context.Context, id int) (types.Kermesse, error)
	Create(ctx context.Context, input map[string]interface{}) error
	Update(ctx context.Context, id int, input map[string]interface{}) error
	AddParticipant(ctx context.Context, input map[string]interface{}) error
	AddStand(ctx context.Context, input map[string]interface{}) error
	End(ctx context.Context, id int) error
}

type Service struct {
	store     KermesseStore
	userStore user.UserStore
}

func NewService(store KermesseStore, userStore user.UserStore) *Service {
	return &Service{
		store:     store,
		userStore: userStore,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]types.Kermesse, error) {
	userId, ok := ctx.Value(types.UserIDKey).(int)
	if !ok {
		return nil, errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("Id utilisateur non trouvé dans le contexte"),
		}
	}
	userRole, ok := ctx.Value(types.UserRoleKey).(string)
	if !ok {
		return nil, errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("Role utilisateur non trouvé dans le contexte"),
		}
	}

	filtres := map[string]interface{}{}
	if userRole == types.UserRoleOrganisateur {
		filtres["organisateur_id"] = userId
	} else if userRole == types.UserRoleParent {
		filtres["parent_id"] = userId
	} else if userRole == types.UserRoleEnfant {
		filtres["child_id"] = userId
	} else if userRole == types.UserRoleTeneurStand {
		filtres["teneur_stand_id"] = userId
	}

	kermesses, err := s.store.FindAll(filtres)
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return kermesses, nil
}

func (s *Service) GetUsersInvite(ctx context.Context, id int) ([]types.UserBasic, error) {
	users, err := s.store.FindUsersInvite(id)
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return users, nil
}

func (s *Service) Get(ctx context.Context, id int) (types.Kermesse, error) {
	kermesse, err := s.store.FindById(id)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return kermesse, errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return kermesse, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return kermesse, nil
}

func (s *Service) Create(ctx context.Context, input map[string]interface{}) error {
	userId, ok := ctx.Value(types.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("ID utilisateur non trouvé dans le contexte"),
		}
	}
	input["user_id"] = userId

	err := s.store.Create(input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (s *Service) Update(ctx context.Context, id int, input map[string]interface{}) error {
	kermesse, err := s.store.FindById(id)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if kermesse.Statut == types.KermesseStatutEnded {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("La kermesse est déjà terminée"),
		}
	}

	userId, ok := ctx.Value(types.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("ID utilisateur non trouvé dans le contexte"),
		}
	}
	if kermesse.UserId != userId {
		return errors.CustomError{
			Key: errors.Forbidden,
			Err: goErrors.New("Interdit de modifier la kermesse d'un autre utilisateur"),
		}
	}

	err = s.store.Update(id, input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (s *Service) AddParticipant(ctx context.Context, input map[string]interface{}) error {
	kermesse, err := s.store.FindById(input["kermesse_id"].(int))
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if kermesse.Statut == types.KermesseStatutEnded {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("La kermesse est déjà terminée"),
		}
	}

	managerId, ok := ctx.Value(types.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("Id utilisateur non trouvé dans le contexte"),
		}
	}
	if kermesse.UserId != managerId {
		return errors.CustomError{
			Key: errors.Forbidden,
			Err: goErrors.New("Interdit"),
		}
	}

	childId, error := utils.GetIntFromMap(input, "user_id")
	if error != nil {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: error,
		}
	}
	child, err := s.userStore.FindById(childId)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}
	if child.Role != types.UserRoleEnfant {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("L'utilisateur n'est pas un enfant"),
		}
	}

	// Inviter l'enfant
	err = s.store.AddParticipant(input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	// Inviter le parent de l'enfant
	if child.ParentId != nil {
		input["user_id"] = child.ParentId
		err = s.store.AddParticipant(input)
		if err != nil {
			return errors.CustomError{
				Key: errors.InternalServerError,
				Err: err,
			}
		}
	}

	return nil
}

func (s *Service) AddStand(ctx context.Context, input map[string]interface{}) error {
	kermesse, err := s.store.FindById(input["kermesse_id"].(int))
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if kermesse.Statut == types.KermesseStatutEnded {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("La kermesse est déjà terminée"),
		}
	}

	standId, err := utils.GetIntFromMap(input, "stand_id")
	if err != nil {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: err,
		}
	}
	canAddStand, err := s.store.CanAddStand(standId)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}
	if !canAddStand {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("Le stand est déjà pris"),
		}
	}

	userId, ok := ctx.Value(types.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("ID utilisateur non trouvé dans le contexte"),
		}
	}
	if kermesse.UserId != userId {
		return errors.CustomError{
			Key: errors.Forbidden,
			Err: goErrors.New("Interdit"),
		}
	}

	err = s.store.AddStand(input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (s *Service) End(ctx context.Context, id int) error {
	kermesse, err := s.store.FindById(id)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if kermesse.Statut == types.KermesseStatutEnded {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("La kermeese est déjà terminée"),
		}
	}

	canEnd, err := s.store.CanEnd(id)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}
	if !canEnd {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("La kermesse ne peut pas être terminée, car il y a une tombola en cours"),
		}
	}

	userId, ok := ctx.Value(types.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("ID utilisateur non trouvé dans le contexte"),
		}
	}
	if kermesse.UserId != userId {
		return errors.CustomError{
			Key: errors.Forbidden,
			Err: goErrors.New("Interdit"),
		}
	}

	err = s.store.End(id)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}
