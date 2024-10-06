package kermesse

import (
	"context"
	"database/sql"
	goErrors "errors"

	"github.com/chall-goflutter-api/internal/types"
	"github.com/chall-goflutter-api/internal/user"
	"github.com/chall-goflutter-api/pkg/errors"
)

type KermesseService interface {
	GetAll(ctx context.Context) ([]types.Kermesse, error)
	Get(ctx context.Context, id int) (types.Kermesse, error)
	Create(ctx context.Context, input map[string]interface{}) error
	Update(ctx context.Context, id int, input map[string]interface{}) error
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
	kermesses, err := s.store.FindAll()
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return kermesses, nil
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
