package stand

import (
	"github.com/chall-goflutter-api/internal/types"
	"github.com/jmoiron/sqlx"
)

type StandStore interface {
	FindAll() ([]types.Stand, error)
	FindById(id int) (types.Stand, error)
	Create(input map[string]interface{}) error
	Update(id int, input map[string]interface{}) error
	UpdateStock(id int, n int) error
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
	queryFindAllStands = "SELECT * FROM stands"
	queryFindStandById = "SELECT * FROM stands WHERE id=$1"
	queryCreateStand   = "INSERT INTO stands (user_id, name, description, type, price, stock) VALUES ($1, $2, $3, $4, $5, $6)"
	queryUpdateStand   = "UPDATE stands SET name=$1, description=$2, price=$3, stock=$4 WHERE id=$5"
	queryUpdateStock   = "UPDATE stands SET stock=stock+$1 WHERE id=$2"
)

func (s *Store) FindAll() ([]types.Stand, error) {
	stands := []types.Stand{}
	err := s.db.Select(&stands, queryFindAllStands)

	return stands, err
}

func (s *Store) FindById(id int) (types.Stand, error) {
	stand := types.Stand{}
	err := s.db.Get(&stand, queryFindStandById, id)

	return stand, err
}

func (s *Store) Create(input map[string]interface{}) error {
	_, err := s.db.Exec(queryCreateStand, input["user_id"], input["name"], input["description"], input["type"], input["price"], input["stock"])

	return err
}

func (s *Store) Update(id int, input map[string]interface{}) error {
	_, err := s.db.Exec(queryUpdateStand, input["name"], input["description"], input["price"], input["stock"], id)

	return err
}

func (s *Store) UpdateStock(id int, quantity int) error {
	_, err := s.db.Exec(queryUpdateStock, quantity, id)

	return err
}
