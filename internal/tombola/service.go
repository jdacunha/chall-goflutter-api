package tombola

import (
	"context"
	"database/sql"
	goErrors "errors"

	"github.com/chall-goflutter-api/internal/kermesse"
	"github.com/chall-goflutter-api/internal/types"
	"github.com/chall-goflutter-api/pkg/errors"
	"github.com/chall-goflutter-api/pkg/utils"
)

type TombolaService interface {
	GetAll(ctx context.Context, params map[string]interface{}) ([]types.Tombola, error)
	Get(ctx context.Context, id int) (types.Tombola, error)
	Create(ctx context.Context, input map[string]interface{}) error
	Update(ctx context.Context, id int, input map[string]interface{}) error
	End(ctx context.Context, id int) error
}

type Service struct {
	store         TombolaStore
	kermesseStore kermesse.KermesseStore
}

func NewService(store TombolaStore, kermesseStore kermesse.KermesseStore) *Service {
	return &Service{
		store:         store,
		kermesseStore: kermesseStore,
	}
}

func (s *Service) GetAll(ctx context.Context, params map[string]interface{}) ([]types.Tombola, error) {
	filters := map[string]interface{}{}
	if params["kermesse_id"] != nil {
		filters["kermesse_id"] = params["kermesse_id"]
	}

	tombolas, err := s.store.FindAll(filters)
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return tombolas, nil
}
func (s *Service) Get(ctx context.Context, id int) (types.Tombola, error) {
	tombola, err := s.store.FindById(id)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return tombola, errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return tombola, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return tombola, nil
}

func (s *Service) Create(ctx context.Context, input map[string]interface{}) error {
	kermesseId, error := utils.GetIntFromMap(input, "kermesse_id")
	if error != nil {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: error,
		}
	}
	kermesse, err := s.kermesseStore.FindById(kermesseId)
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
			Err: goErrors.New("La kermesse est terminée"),
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

	err = s.store.Create(input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (s *Service) Update(ctx context.Context, id int, input map[string]interface{}) error {
	tombola, err := s.store.FindById(id)
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

	kermesse, err := s.kermesseStore.FindById(tombola.KermesseId)
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
			Err: goErrors.New("La kermesse est terminée"),
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
			Err: goErrors.New("Inrerdit"),
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

func (s *Service) End(ctx context.Context, id int) error {
	tombola, err := s.store.FindById(id)
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

	kermesse, err := s.kermesseStore.FindById(tombola.KermesseId)
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
			Err: goErrors.New("La kermeese est terminée"),
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

	if tombola.Statut == types.TombolaStatutEnded {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("La tombola est déjà terminée"),
		}
	}

	err = s.store.SelectGagnant(id)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}
