package kermesse

import (
	"fmt"

	"github.com/chall-goflutter-api/internal/types"
	"github.com/jmoiron/sqlx"
)

type KermesseStore interface {
	FindAll(filtres map[string]interface{}) ([]types.Kermesse, error)
	FindUsersInvite(id int) ([]types.UserBasic, error)
	FindById(id int) (types.Kermesse, error)
	Create(input map[string]interface{}) error
	Update(id int, input map[string]interface{}) error
	AddParticipant(input map[string]interface{}) error
	CanAddStand(standId int) (bool, error)
	AddStand(input map[string]interface{}) error
	CanEnd(id int) (bool, error)
	End(id int) error
	// TODO // Stats(id int, filtres map[string]interface{}) (types.KermesseStats, error)
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
	queryFindAllKermesses = "SELECT * FROM kermesses"
	queryFindKermesseById = "SELECT * FROM kermesses WHERE id=$1"
	queryCreateKermesse   = "INSERT INTO kermesses (user_id, name, description) VALUES ($1, $2, $3)"
	queryUpdateKermesse   = "UPDATE kermesses SET name=$1, description=$2 WHERE id=$3"
	queryAddParticipant   = "INSERT INTO kermesses_users (kermesse_id, user_id) VALUES ($1, $2)"
	queryAddStand         = "INSERT INTO kermesses_stands (kermesse_id, stand_id) VALUES ($1, $2)"
	queryCanEnd           = "SELECT EXISTS ( SELECT 1 FROM tombolas WHERE kermesse_id = $1 AND statut = $2 ) AS is_true"
	queryEnd              = "UPDATE kermesses SET statut=$1 WHERE id=$2"
)

func (s *Store) FindAll(filtres map[string]interface{}) ([]types.Kermesse, error) {
	kermesses := []types.Kermesse{}
	query := `
		SELECT DISTINCT
			k.id AS id,
			k.user_id AS user_id,
			k.name AS name,
			k.description AS description,
			k.statut AS statut
		FROM kermesses k
		FULL OUTER JOIN kermesses_users ku ON k.id = ku.kermesse_id
		FULL OUTER JOIN kermesses_stands ks ON k.id = ks.kermesse_id
		FULL OUTER JOIN stands s ON ks.stand_id = s.id
		WHERE 1=1
	`
	if filtres["organisateur_id"] != nil {
		query += fmt.Sprintf(" AND k.user_id = %v", filtres["organisateur_id"])
	}
	if filtres["parent_id"] != nil {
		query += fmt.Sprintf(" AND ku.user_id = %v", filtres["parent_id"])
	}
	if filtres["child_id"] != nil {
		query += fmt.Sprintf(" AND ku.user_id = %v", filtres["child_id"])
	}
	if filtres["teneur_stand_id"] != nil {
		query += fmt.Sprintf(" AND ks.stand_id IS NOT NULL AND s.user_id = %v", filtres["teneur_stand_id"])
	}
	err := s.db.Select(&kermesses, query)

	return kermesses, err
}

func (s *Store) FindUsersInvite(id int) ([]types.UserBasic, error) {
	users := []types.UserBasic{}
	query := `
		SELECT DISTINCT
			u.id AS id,
			u.name AS name,
			u.email AS email,
			u.role AS role,
			u.jetons AS jetons
		FROM users u
		LEFT JOIN kermesses_users ku ON u.id = ku.user_id AND ku.kermesse_id = $1
		WHERE u.role = 'ENFANT'
		AND ku.user_id IS NULL;
	`
	err := s.db.Select(&users, query, id)

	return users, err
}

func (s *Store) FindById(id int) (types.Kermesse, error) {
	kermesse := types.Kermesse{}
	err := s.db.Get(&kermesse, queryFindKermesseById, id)

	return kermesse, err
}

func (s *Store) Create(input map[string]interface{}) error {
	_, err := s.db.Exec(queryCreateKermesse, input["user_id"], input["name"], input["description"])

	return err
}

func (s *Store) Update(id int, input map[string]interface{}) error {
	_, err := s.db.Exec(queryUpdateKermesse, input["name"], input["description"], id)

	return err
}

func (s *Store) AddParticipant(input map[string]interface{}) error {
	_, err := s.db.Exec(queryAddParticipant, input["kermesse_id"], input["user_id"])

	return err
}

func (s *Store) CanAddStand(standId int) (bool, error) {
	var isTrue bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM kermesses_stands ks
  		JOIN kermesses k ON ks.kermesse_id = k.id
  		WHERE ks.stand_id = $1 AND k.statut = $2
		) AS is_associated
 	`
	err := s.db.QueryRow(query, standId, types.KermesseStatutStarted).Scan(&isTrue)

	return !isTrue, err
}

func (s *Store) AddStand(input map[string]interface{}) error {
	_, err := s.db.Exec(queryAddStand, input["kermesse_id"], input["stand_id"])

	return err
}

func (s *Store) CanEnd(id int) (bool, error) {
	var isTrue bool
	err := s.db.QueryRow(queryCanEnd, id, types.TombolaStatutStarted).Scan(&isTrue)

	return !isTrue, err
}

func (s *Store) End(id int) error {
	_, err := s.db.Exec(queryEnd, types.KermesseStatutEnded, id)

	return err
}
