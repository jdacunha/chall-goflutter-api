package tombola

import (
	"fmt"

	"github.com/chall-goflutter-api/internal/types"
	"github.com/jmoiron/sqlx"
)

type TombolaStore interface {
	FindAll(filters map[string]interface{}) ([]types.Tombola, error)
	FindById(id int) (types.Tombola, error)
	Create(input map[string]interface{}) error
	Update(id int, input map[string]interface{}) error
	SelectGagnant(id int) error
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
	queryFindTombolaById = "SELECT * FROM tombolas WHERE id=$1"
	queryCreateTombola   = "INSERT INTO tombolas (kermesse_id, name, price, lot) VALUES ($1, $2, $3, $4)"
	queryUpdateTombola   = "UPDATE tombolas SET name=$1, price=$2, lot=$3 WHERE id=$4"
	queryUpdateStatut    = "UPDATE tombolas SET statut=$1 WHERE id=$2"
)

func (s *Store) FindAll(filters map[string]interface{}) ([]types.Tombola, error) {
	tombolas := []types.Tombola{}
	query := `
		SELECT DISTINCT
			t.id AS id,
			t.kermesse_id AS kermesse_id,
			t.name AS name,
			t.statut AS statut,
			t.price AS price,
			t.lot AS lot
		FROM tombolas t
		WHERE 1=1
	`
	if filters["kermesse_id"] != nil {
		query += fmt.Sprintf(" AND t.kermesse_id = %v", filters["kermesse_id"])
	}
	err := s.db.Select(&tombolas, query)
	return tombolas, err
}

func (s *Store) FindById(id int) (types.Tombola, error) {
	tombola := types.Tombola{}
	err := s.db.Get(&tombola, queryFindTombolaById, id)

	return tombola, err
}

func (s *Store) Create(input map[string]interface{}) error {
	_, err := s.db.Exec(queryCreateTombola, input["kermesse_id"], input["name"], input["price"], input["lot"])

	return err
}

func (s *Store) Update(id int, input map[string]interface{}) error {
	_, err := s.db.Exec(queryUpdateTombola, input["name"], input["price"], input["lot"], id)

	return err
}

// Sélectionne un gagnant aléatoire parmi les tickets d'une tombola et termine la tombola.
func (s *Store) SelectGagnant(id int) error {
	tx, err := s.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	if _, err = tx.Exec(queryUpdateStatut, types.TombolaStatutEnded, id); err != nil {
		return err
	}

	const querySelectGagnant = `
		UPDATE tickets
		SET gagnant = true
		WHERE id = (
			SELECT id
			FROM tickets
			WHERE tombola_id = $1
			ORDER BY RANDOM()
			LIMIT 1
		)
		AND tombola_id = $1`
	_, err = tx.Exec(querySelectGagnant, id)
	return err
}
