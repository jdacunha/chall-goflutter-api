package interaction

import (
	"fmt"

	"github.com/chall-goflutter-api/internal/types"
	"github.com/jmoiron/sqlx"
)

type InteractionStore interface {
	FindAll(filters map[string]interface{}) ([]types.InteractionBasic, error)
	FindById(id int) (types.Interaction, error)
	CanCreate(input map[string]interface{}) (bool, error)
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
	queryCreateInteraction = "INSERT INTO interactions (user_id, kermesse_id, stand_id, type, jetons) VALUES ($1, $2, $3, $4, $5)"
	queryUpdateInteraction = "UPDATE interactions SET statut=$1, points=$2 WHERE id=$3"
)

func (s *Store) FindAll(filters map[string]interface{}) ([]types.InteractionBasic, error) {
	interactions := []types.InteractionBasic{}
	query := `
		SELECT DISTINCT
			i.id AS id,
			i.type AS type,
			i.statut AS statut,
			i.jetons AS jetons,
			i.points AS points,
			i.created_at AS created_at,
			u.id AS "user.id",
			u.name AS "user.name",
			u.email AS "user.email",
			u.role AS "user.role",
			s.id AS "stand.id",
			s.name AS "stand.name",
			s.description AS "stand.description",
			s.type AS "stand.type",
			s.price AS "stand.price"
		FROM interactions i
		JOIN users u ON i.user_id = u.id
		JOIN stands s ON i.stand_id = s.id
		WHERE 1=1
	`
	if filters["kermesse_id"] != nil {
		query += fmt.Sprintf(" AND i.kermesse_id = %v", filters["kermesse_id"])
	}
	if filters["parent_id"] != nil {
		query += fmt.Sprintf(" AND (u.id = %v OR u.parent_id = %v)", filters["parent_id"], filters["parent_id"])
	}
	if filters["enfant_id"] != nil {
		query += fmt.Sprintf(" AND u.id = %v", filters["enfant_id"])
	}
	if filters["teneur_stand_id"] != nil {
		query += fmt.Sprintf(" AND s.user_id = %v", filters["teneur_stand_id"])
	}
	query += " ORDER BY i.created_at DESC"
	err := s.db.Select(&interactions, query)

	return interactions, err
}

func (s *Store) FindById(id int) (types.Interaction, error) {
	interaction := types.Interaction{}
	query := `
		SELECT
			i.id AS id,
			i.type AS type,
			i.statut AS statut,
			i.jetons AS jetons,
			i.points AS points,
			i.created_at AS created_at,
			u.id AS "user.id",
			u.name AS "user.name",
			u.email AS "user.email",
			u.role AS "user.role",
			s.id AS "stand.id",
			s.name AS "stand.name",
			s.description AS "stand.description",
			s.type AS "stand.type",
			s.price AS "stand.price",
			k.id AS "kermesse.id",
			k.name AS "kermesse.name",
			k.description AS "kermesse.description",
			k.statut AS "kermesse.statut"
		FROM interactions i
		JOIN users u ON i.user_id = u.id
		JOIN stands s ON i.stand_id = s.id
		JOIN kermesses k ON i.kermesse_id = k.id
		WHERE i.id=$1
	`
	err := s.db.Get(&interaction, query, id)

	return interaction, err
}

func (s *Store) CanCreate(input map[string]interface{}) (bool, error) {
	var isAssociated bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM kermesses_users ku
  		JOIN kermesses_stands ks ON ku.kermesse_id = ks.kermesse_id
			JOIN kermesses k ON ku.kermesse_id = k.id
  		WHERE ku.user_id = $1 AND ks.stand_id = $2 AND k.statut = $3
		) AS is_associated
 	`
	err := s.db.QueryRow(query, input["user_id"], input["stand_id"], types.KermesseStatutStarted).Scan(&isAssociated)

	return isAssociated, err
}

func (s *Store) Create(input map[string]interface{}) error {
	_, err := s.db.Exec(queryCreateInteraction, input["user_id"], input["kermesse_id"], input["stand_id"], input["type"], input["jetons"])

	return err
}

func (s *Store) Update(id int, input map[string]interface{}) error {
	_, err := s.db.Exec(queryUpdateInteraction, input["statut"], input["points"], id)

	return err
}
