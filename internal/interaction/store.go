package interaction

import (
	"github.com/chall-goflutter-api/internal/types"
	"github.com/jmoiron/sqlx"
)

type InteractionStore interface {
	FindAll() ([]types.Interaction, error)
	FindById(id int) (types.Interaction, error)
	Create(input map[string]interface{}) error
	Update(id int, input map[string]interface{}) error
}

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

const (
	queryFindAllInteractions = "SELECT * FROM interactions"
	queryFindInteractionById = "SELECT * FROM interactions WHERE id=$1"
	queryCreateInteraction   = "INSERT INTO interactions (user_id, kermesse_id, stand_id, type, jetons) VALUES ($1, $2, $3, $4, $5)"
	queryUpdateInteraction   = "UPDATE interactions SET statut=$1, points=$2 WHERE id=$3"
)

func (s *Store) FindAll() ([]types.Interaction, error) {
	interactions := []types.Interaction{}
	err := s.db.Select(&interactions, queryFindAllInteractions)

	return interactions, err
}

func (s *Store) FindById(id int) (types.Interaction, error) {
	interaction := types.Interaction{}
	err := s.db.Get(&interaction, queryFindInteractionById, id)

	return interaction, err
}

func (s *Store) Create(input map[string]interface{}) error {
	_, err := s.db.Exec(queryCreateInteraction, input["user_id"], input["kermesse_id"], input["stand_id"], input["type"], input["jetons"])

	return err
}

func (s *Store) Update(id int, input map[string]interface{}) error {
	_, err := s.db.Exec(queryUpdateInteraction, input["statut"], input["points"], id)

	return err
}
