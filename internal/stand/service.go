package stand

import (
	"context"
	"database/sql"
	goErrors "errors"

	"github.com/chall-goflutter-api/internal/types"
	"github.com/chall-goflutter-api/pkg/errors"
)

type StandService interface {
	GetAll(ctx context.Context, params map[string]interface{}) ([]types.Stand, error)
	Get(ctx context.Context, id int) (types.Stand, error)
	Create(ctx context.Context, input map[string]interface{}) error
	Update(ctx context.Context, id int, input map[string]interface{}) error
	GetCurrent(ctx context.Context) (types.Stand, error)
	UpdateCurrent(ctx context.Context, input map[string]interface{}) error
}

type Service struct {
	store StandStore
}

func NewService(store StandStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetAll(ctx context.Context, params map[string]interface{}) ([]types.Stand, error) {
	filtres := map[string]interface{}{}
	if params["kermesse_id"] != nil {
		filtres["kermesse_id"] = params["kermesse_id"]
	}
	if params["is_libre"] != nil {
		filtres["is_libre"] = params["is_libre"]
	}

	stands, err := s.store.FindAll(params)
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return stands, nil
}

func (s *Service) Get(ctx context.Context, id int) (types.Stand, error) {
	stand, err := s.store.FindById(id)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return stand, errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return stand, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return stand, nil
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
	stand, err := s.store.FindById(id)
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

	userId, ok := ctx.Value(types.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("ID utilisateur non trouvé dans le contexte"),
		}
	}
	if stand.UserId != userId {
		return errors.CustomError{
			Key: errors.Forbidden,
			Err: goErrors.New("Interdit"),
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

func (s *Service) GetCurrent(ctx context.Context) (types.Stand, error) {
	userId, ok := ctx.Value(types.UserIDKey).(int)
	if !ok {
		return types.Stand{}, errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("ID utilisateur non trouvé dans le contexte"),
		}
	}

	stand, err := s.store.FindByUserId(userId)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return stand, errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return stand, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return stand, nil
}

func (s *Service) UpdateCurrent(ctx context.Context, input map[string]interface{}) error {
	userId, ok := ctx.Value(types.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("ID utilisateur non trouvé dans le contexte"),
		}
	}

	err := s.store.UpdateByUserId(userId, input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}
