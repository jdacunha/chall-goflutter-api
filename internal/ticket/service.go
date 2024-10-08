package ticket

import (
	"context"
	"database/sql"
	goErrors "errors"

	"github.com/chall-goflutter-api/internal/tombola"
	"github.com/chall-goflutter-api/internal/types"
	"github.com/chall-goflutter-api/internal/user"
	"github.com/chall-goflutter-api/pkg/errors"
	"github.com/chall-goflutter-api/pkg/utils"
)

type TicketService interface {
	GetAll(ctx context.Context) ([]types.Ticket, error)
	Get(ctx context.Context, id int) (types.Ticket, error)
	Create(ctx context.Context, input map[string]interface{}) error
}

type Service struct {
	store        TicketStore
	tombolaStore tombola.TombolaStore
	userStore    user.UserStore
}

func NewService(store TicketStore, tombolaStore tombola.TombolaStore, userStore user.UserStore) *Service {
	return &Service{
		store:        store,
		tombolaStore: tombolaStore,
		userStore:    userStore,
	}
}

func (s *Service) GetAll(ctx context.Context) ([]types.Ticket, error) {
	tickets, err := s.store.FindAll()
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return tickets, nil
}

func (s *Service) Get(ctx context.Context, id int) (types.Ticket, error) {
	ticket, err := s.store.FindById(id)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return ticket, errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return ticket, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return ticket, nil
}

func (s *Service) Create(ctx context.Context, input map[string]interface{}) error {
	tombolaId, err := utils.GetIntFromMap(input, "tombola_id")
	if err != nil {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: err,
		}
	}
	tombola, err := s.tombolaStore.FindById(tombolaId)
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

	if tombola.Statut == types.TombolaStatutEnded {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("La tombola est déjà terminée"),
		}
	}

	userId, ok := ctx.Value(types.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("ID utilisateur non trouvé dans le contexte"),
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

	// Check si l'utilisateur a assez de jetons
	if user.Jetons < tombola.Price {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("Pas assez de jetons"),
		}
	}

	// Check si l'utilisateur peut participer à la tombola
	canCreate, err := s.store.CanCreate(map[string]interface{}{
		"kermesse_id": tombola.KermesseId,
		"user_id":     userId,
	})
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}
	if !canCreate {
		return errors.CustomError{
			Key: errors.Forbidden,
			Err: goErrors.New("Interdit"),
		}
	}

	// Mettre à jour les jetons de l'utilisateur
	err = s.userStore.UpdateJetons(userId, -tombola.Price)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	input["user_id"] = userId

	err = s.store.Create(input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}
