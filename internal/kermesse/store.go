package kermesse

import (
	"github.com/jmoiron/sqlx"
	"github.com/kermanager/internal/types"
)

type KermesseStore interface {
	FindAll() ([]types.Kermesse, error)
	FindById(id int) (types.Kermesse, error)
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
	queryFindAllKermesses = "SELECT * FROM kermesses"
	queryFindKermesseById = "SELECT * FROM kermesses WHERE id=$1"
	queryCreateKermesse   = "INSERT INTO kermesses (user_id, name, description) VALUES ($1, $2, $3)"
	queryUpdateKermesse   = "UPDATE kermesses SET name=$1, description=$2 WHERE id=$3"
)

func (s *Store) FindAll() ([]types.Kermesse, error) {
	kermesses := []types.Kermesse{}
	err := s.db.Select(&kermesses, queryFindAllKermesses)

	return kermesses, err
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
