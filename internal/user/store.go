package user

import (
	"fmt"

	"github.com/chall-goflutter-api/internal/types"
	"github.com/jmoiron/sqlx"
)

type UserStore interface {
	FindAll(filtres map[string]interface{}) ([]types.UserBasic, error)
	FindChildren(id int, filtres map[string]interface{}) ([]types.UserBasic, error)
	FindById(id int) (types.User, error)
	FindByEmail(email string) (types.User, error)
	Create(input map[string]interface{}) error
	UpdatePassword(id int, input map[string]interface{}) error
	UpdateJetons(id int, amount int) error
	HasStand(id int) (bool, error)
}

type Store struct {
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) FindAll(filtres map[string]interface{}) ([]types.UserBasic, error) {
	users := []types.UserBasic{}
	query := `
		SELECT DISTINCT
			u.id AS id,
			u.name AS name,
			u.email AS email,
			u.role AS role,
			u.jetons AS jetons
		FROM users u
		FULL OUTER JOIN kermesses_users ku ON u.id = ku.user_id
		WHERE 1=1
	`
	if filtres["kermesse_id"] != nil {
		query += fmt.Sprintf(" AND ku.kermesse_id = %v", filtres["kermesse_id"])
	}
	err := s.db.Select(&users, query)

	return users, err
}

func (s *Store) FindChildren(id int, filtres map[string]interface{}) ([]types.UserBasic, error) {
	users := []types.UserBasic{}
	query := `
		SELECT DISTINCT
			u.id AS id,
			u.name AS name,
			u.email AS email,
			u.role AS role,
			u.jetons AS jetons
		FROM users u
		FULL OUTER JOIN kermesses_users ku ON u.id = ku.user_id
		WHERE u.role=$1 AND u.parent_id=$2
	`
	if filtres["kermesse_id"] != nil {
		query += fmt.Sprintf(" AND ku.kermesse_id = %v", filtres["kermesse_id"])
	}
	err := s.db.Select(&users, query, types.UserRoleEnfant, id)

	return users, err
}

func (s *Store) FindById(id int) (types.User, error) {
	user := types.User{}
	query := "SELECT * FROM users WHERE id=$1"
	err := s.db.Get(&user, query, id)

	return user, err
}

func (s *Store) FindByEmail(email string) (types.User, error) {
	user := types.User{}
	query := "SELECT * FROM users WHERE email=$1"
	err := s.db.Get(&user, query, email)

	return user, err
}

func (s *Store) Create(input map[string]interface{}) error {
	query := "INSERT INTO users (parent_id, name, email, password_hash, role) VALUES ($1, $2, $3, $4, $5)"
	_, err := s.db.Exec(query, input["parent_id"], input["name"], input["email"], input["password"], input["role"])

	return err
}

func (s *Store) UpdatePassword(id int, input map[string]interface{}) error {
	query := "UPDATE users SET password_hash=$1 WHERE id=$2"
	_, err := s.db.Exec(query, input["new_password"], id)

	return err
}

func (s *Store) UpdateJetons(id int, amount int) error {
	query := "UPDATE users SET jetons=jetons+$1 WHERE id=$2"
	_, err := s.db.Exec(query, amount, id)

	return err
}

func (s *Store) HasStand(id int) (bool, error) {
	var count int
	query := "SELECT COUNT(*) FROM stands WHERE user_id=$1"
	err := s.db.Get(&count, query, id)

	return count >= 1, err
}
