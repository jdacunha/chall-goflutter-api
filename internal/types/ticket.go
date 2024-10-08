package types

import "time"

type TicketUser struct {
	Id    int    `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Role  string `json:"role" db:"role"`
}

type TicketTombola struct {
	Id     int    `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Statut string `json:"statut" db:"statut"`
	Price  int    `json:"price" db:"price"`
	Lot    string `json:"lot" db:"lot"`
}

type TicketKermesse struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Statut      string `json:"statut" db:"statut"`
}

type Ticket struct {
	Id        int            `json:"id" db:"id"`
	Gagnant   bool           `json:"gagnant" db:"gagnant"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	User      TicketUser     `json:"user" db:"user"`
	Tombola   TicketTombola  `json:"tombola" db:"tombola"`
	Kermesse  TicketKermesse `json:"kermesse" db:"kermesse"`
}
