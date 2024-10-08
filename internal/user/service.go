package user

import (
	"context"
	"database/sql"
	goErrors "errors"
	"os"
	"strconv"

	"github.com/chall-goflutter-api/internal/types"
	"github.com/chall-goflutter-api/pkg/errors"
	"github.com/chall-goflutter-api/pkg/generator"
	"github.com/chall-goflutter-api/pkg/hasher"
	"github.com/chall-goflutter-api/pkg/jwt"
	"github.com/chall-goflutter-api/pkg/utils"
	goJwt "github.com/golang-jwt/jwt/v5"
)

type UserService interface {
	GetAll(ctx context.Context, params map[string]interface{}) ([]types.UserBasic, error)
	GetChildren(ctx context.Context, params map[string]interface{}) ([]types.UserBasic, error)
	Get(ctx context.Context, id int) (types.UserBasic, error)
	UpdatePassword(ctx context.Context, id int, input map[string]interface{}) error
	UpdateJetons(userId, credit int) error
	Invite(ctx context.Context, input map[string]interface{}) error
	Distribute(ctx context.Context, input map[string]interface{}) error
	Register(ctx context.Context, input map[string]interface{}) error
	Login(ctx context.Context, input map[string]interface{}) (types.UserBasicWithToken, error)
	GetMe(ctx context.Context) (types.UserBasic, error)
}

type Service struct {
	store UserStore
}

func NewService(store UserStore) *Service {
	return &Service{
		store: store,
	}
}

func (s *Service) GetAll(ctx context.Context, params map[string]interface{}) ([]types.UserBasic, error) {
	filtres := map[string]interface{}{}
	if params["kermesse_id"] != nil {
		filtres["kermesse_id"] = params["kermesse_id"]
	}

	users, err := s.store.FindAll(filtres)
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return users, nil
}

func (s *Service) GetChildren(ctx context.Context, params map[string]interface{}) ([]types.UserBasic, error) {
	filtres := map[string]interface{}{}
	if params["kermesse_id"] != nil {
		filtres["kermesse_id"] = params["kermesse_id"]
	}

	userId, ok := ctx.Value(types.UserIDKey).(int)
	if !ok {
		return nil, errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("ID utilisateur non trouvé dans le contexte"),
		}
	}

	users, err := s.store.FindChildren(userId, filtres)
	if err != nil {
		return nil, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return users, nil
}

func (s *Service) Get(ctx context.Context, id int) (types.UserBasic, error) {
	user, err := s.store.FindById(id)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return types.UserBasic{}, errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return types.UserBasic{}, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return types.UserBasic{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
		Jetons: user.Jetons,
	}, nil
}

func (s *Service) UpdatePassword(ctx context.Context, id int, input map[string]interface{}) error {
	user, err := s.store.FindById(id)
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
	if user.Id != userId {
		return errors.CustomError{
			Key: errors.Forbidden,
			Err: goErrors.New("Interdit"),
		}
	}

	if !hasher.Compare(user.PasswordHash, input["password"].(string)) {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("Mot de passe incorrect"),
		}
	}

	hashedPassword, err := hasher.Hash(input["new_password"].(string))
	if err != nil {
		return err
	}
	input["new_password"] = hashedPassword

	err = s.store.UpdatePassword(id, input)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (s *Service) UpdateJetons(userId, credit int) error {
	user, err := s.store.FindById(userId)
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
	if user.Role != types.UserRoleParent {
		return errors.CustomError{
			Key: errors.Forbidden,
			Err: goErrors.New("Interdit"),
		}
	}

	err = s.store.UpdateJetons(userId, credit)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (s *Service) Invite(ctx context.Context, input map[string]interface{}) error {
	_, err := s.store.FindByEmail(input["email"].(string))
	if err == nil {
		return errors.CustomError{
			Key: errors.EmailAlreadyExists,
			Err: goErrors.New("Email déjà utilisé"),
		}
	}

	randomPassword, err := generator.RandomPassword(8)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	hashedPassword, err := hasher.Hash(randomPassword)
	if err != nil {
		return err
	}

	userId, ok := ctx.Value(types.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("ID utilisateur non trouvé dans le contexte"),
		}
	}

	err = s.store.Create(map[string]interface{}{
		"name":      input["name"],
		"email":     input["email"],
		"password":  hashedPassword,
		"role":      types.UserRoleEnfant,
		"parent_id": userId,
	})
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}
	// Faire ensuite l'envoi d'email
	return nil
}

func (s *Service) Distribute(ctx context.Context, input map[string]interface{}) error {
	childId, err := utils.GetIntFromMap(input, "child_id")
	if err != nil {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: err,
		}
	}
	child, err := s.store.FindById(childId)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("Enfant non trouvé"),
		}
	}

	parentId, ok := ctx.Value(types.UserIDKey).(int)
	if !ok {
		return errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("ID utilisateur non trouvé dans le contexte"),
		}
	}
	parent, err := s.store.FindById(parentId)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("Parent non trouvé"),
		}
	}

	if child.ParentId == nil || *child.ParentId != parent.Id {
		return errors.CustomError{
			Key: errors.Forbidden,
			Err: goErrors.New("Interdit"),
		}
	}

	amount, error := utils.GetIntFromMap(input, "montant")
	if error != nil {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: error,
		}
	}
	if parent.Jetons < amount {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("Jetons insuffisants"),
		}
	}

	err = s.store.UpdateJetons(childId, amount)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	err = s.store.UpdateJetons(parentId, -amount)
	if err != nil {
		return errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return nil
}

func (s *Service) Register(ctx context.Context, input map[string]interface{}) error {
	_, err := s.store.FindByEmail(input["email"].(string))
	if err == nil {
		return errors.CustomError{
			Key: errors.EmailAlreadyExists,
			Err: goErrors.New("Email déjà utilisé"),
		}
	}

	hashedPassword, err := hasher.Hash(input["password"].(string))
	if err != nil {
		return err
	}
	input["password"] = hashedPassword
	input["parent_id"] = nil

	if input["role"] == types.UserRoleEnfant {
		return errors.CustomError{
			Key: errors.BadRequest,
			Err: goErrors.New("Role non autorisé"),
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

func (s *Service) Login(ctx context.Context, input map[string]interface{}) (types.UserBasicWithToken, error) {
	user, err := s.store.FindByEmail(input["email"].(string))
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return types.UserBasicWithToken{}, errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return types.UserBasicWithToken{}, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	if !hasher.Compare(user.PasswordHash, input["password"].(string)) {
		return types.UserBasicWithToken{}, errors.CustomError{
			Key: errors.InvalidCredentials,
			Err: goErrors.New("Invalid credentials"),
		}
	}

	expiresIn, err := strconv.Atoi(os.Getenv("JWT_EXPIRES_IN"))
	if err != nil {
		return types.UserBasicWithToken{}, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	token, err := jwt.Create(os.Getenv("JWT_SECRET"), expiresIn, user.Id)
	if err != nil {
		if goErrors.Is(err, goJwt.ErrTokenExpired) || goErrors.Is(err, goJwt.ErrSignatureInvalid) {
			return types.UserBasicWithToken{}, errors.CustomError{
				Key: errors.Unauthorized,
				Err: err,
			}
		}
		return types.UserBasicWithToken{}, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return types.UserBasicWithToken{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
		Jetons: user.Jetons,
		Token:  token,
	}, nil
}

func (s *Service) GetMe(ctx context.Context) (types.UserBasic, error) {
	userId, ok := ctx.Value(types.UserIDKey).(int)
	if !ok {
		return types.UserBasic{}, errors.CustomError{
			Key: errors.Unauthorized,
			Err: goErrors.New("Id utilisateur non trouvé dans le contexte"),
		}
	}

	user, err := s.store.FindById(userId)
	if err != nil {
		if goErrors.Is(err, sql.ErrNoRows) {
			return types.UserBasic{}, errors.CustomError{
				Key: errors.NotFound,
				Err: err,
			}
		}
		return types.UserBasic{}, errors.CustomError{
			Key: errors.InternalServerError,
			Err: err,
		}
	}

	return types.UserBasic{
		Id:     user.Id,
		Name:   user.Name,
		Email:  user.Email,
		Role:   user.Role,
		Jetons: user.Jetons,
	}, nil
}
