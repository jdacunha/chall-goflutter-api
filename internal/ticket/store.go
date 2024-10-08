package ticket

import (
	"github.com/chall-goflutter-api/internal/types"
	"github.com/jmoiron/sqlx"
)

type TicketStore interface {
	FindAll() ([]types.Ticket, error)
	FindById(id int) (types.Ticket, error)
	Create(input map[string]interface{}) error
	CanCreate(input map[string]interface{}) (bool, error)
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
	queryFindAllTickets = "SELECT * FROM tickets"
	queryFindTicketById = "SELECT * FROM tickets WHERE id=$1"
	queryCreateTicket   = "INSERT INTO tickets (user_id, tombola_id) VALUES ($1, $2)"
)

func (s *Store) FindAll() ([]types.Ticket, error) {
	tickets := []types.Ticket{}
	err := s.db.Select(&tickets, queryFindAllTickets)

	return tickets, err
}

func (s *Store) FindById(id int) (types.Ticket, error) {
	ticket := types.Ticket{}
	err := s.db.Get(&ticket, queryFindTicketById, id)

	return ticket, err
}

func (s *Store) CanCreate(input map[string]interface{}) (bool, error) {
	var isAssociated bool
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM kermesses_users ku
			JOIN kermesses k ON k.id = ku.kermesse_id
			WHERE ku.kermesse_id = $1 AND ku.user_id = $2 AND k.statut = $3
		) AS is_associated
	`
	err := s.db.QueryRow(query, input["kermesse_id"], input["user_id"], types.KermesseStatutStarted).Scan(&isAssociated)

	return isAssociated, err
}

func (s *Store) Create(input map[string]interface{}) error {
	_, err := s.db.Exec(queryCreateTicket, input["user_id"], input["tombola_id"])

	return err
}
