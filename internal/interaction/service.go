package interaction

import (
	"context"
	"database/sql"
	goErrors "errors"

	"github.com/chall-goflutter-api/internal/kermesse"
	"github.com/chall-goflutter-api/internal/stand"
	"github.com/chall-goflutter-api/internal/types"
	"github.com/chall-goflutter-api/internal/user"
	"github.com/chall-goflutter-api/pkg/errors"
	"github.com/chall-goflutter-api/pkg/utils"
)

type InteractionService interface {
	GetAll(ctx context.Context) ([]types.Interaction, error)
	Get(ctx context.Context, id int) (types.Interaction, error)
	Create(ctx context.Context, input map[string]interface{}) error
	Update(ctx context.Context, id int, input map[string]interface{}) error
}

type Service struct {
	store         InteractionStore
	standStore    stand.StandStore
	userStore     user.UserStore
	kermesseStore kermesse.KermesseStore
}

func NewService(store InteractionStore, standStore stand.StandStore, userStore user.UserStore, kermesseStore kermesse.KermesseStore) *Service {
	return &Service{
		store:         store,
		standStore:    standStore,
		userStore:     userStore,
		kermesseStore: kermesseStore,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]types.Interaction, error) {
	interactions, err := s.store.FindAll()
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return interactions, nil
}

func (s *Service) Get(ctx context.Context, id int) (types.Interaction, error) {
	interaction, err := s.store.FindById(id)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return interaction, errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return interaction, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return interaction, nil
}

func (s *Service) Create(ctx context.Context, input map[string]interface{}) error {
	standId, err := utils.GetIntFromMap(input, "stand_id")
	if err != nil {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: err,
		}
	}
	stand, err := s.standStore.FindById(standId)
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

	userId, err := utils.GetIntFromMap(input, "user_id")
	if err != nil {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: err,
		}
	}
	user, err := s.userStore.FindById(userId)
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

	if stand.Type == types.InteractionTypeTransaction {
		quantity, err := utils.GetIntFromMap(input, "quantity")
		if err != nil {
			return errors.CustomError{
				Key: errors.BadRequest,
				Err: err,
			}
		}
		// mettre à jour le stock du stand
		if stand.Stock < quantity {
			return errors.CustomError{
				Key: errors.BadRequest,
				Err: goErrors.New("Pas assez de stock"),
			}
		}
		err = s.standStore.UpdateStock(standId, -quantity)
		if err != nil {
			return errors.CustomError{
				Key: errors.InternalServerError,
				Err: err,
			}
		}
		// mettre à jour les jetons de l'utilisateur
		totalPrice := stand.Price * quantity
		if user.Jetons < totalPrice {
			return errors.CustomError{
				Key: errors.BadRequest,
				Err: goErrors.New("Pas assez de jetons"),
			}
		}
		err = s.userStore.UpdateJetons(userId, -totalPrice)
		if err != nil {
			return errors.CustomError{
				Key: errors.InternalServerError,
				Err: err,
			}
		}
	} else {
		// mettre à jour les jetons de l'utilisateur
		totalPrice := stand.Price
		if user.Jetons < totalPrice {
			return errors.CustomError{
				Key: errors.BadRequest,
				Err: goErrors.New("Pas assez de jetons"),
			}
		}
		err = s.userStore.UpdateJetons(userId, -totalPrice)
		if err != nil {
			return errors.CustomError{
				Key: errors.InternalServerError,
				Err: err,
			}
		}
	}

	input["type"] = stand.Type
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
	interaction, err := s.store.FindById(id)
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

	if interaction.Type != types.InteractionTypeActivite {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("L'interaction n'est pas une activité"),
		}
	}

	err = s.store.Update(id, map[string]interface{}{
		"statut": types.InteractionStatutEnded,
		"points": input["points"],
	})
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}
