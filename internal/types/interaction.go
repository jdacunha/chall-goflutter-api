package types

import "time"

const (
	InteractionTypeTransaction string = "TRANSACTION"
	InteractionTypeActivite    string = "ACTIVITE"
	InteractionStatutStarted   string = "STARTED"
	InteractionStatutEnded     string = "ENDED"
)

type InteractionUser struct {
	Id    int    `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Email string `json:"email" db:"email"`
	Role  string `json:"role" db:"role"`
}

type InteractionStand struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Type        string `json:"type" db:"type"`
	Price       int    `json:"price" db:"price"`
}

type InteractionKermesse struct {
	Id          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Statut      string `json:"statut" db:"statut"`
}

type Interaction struct {
	Id        int                 `json:"id" db:"id"`
	Type      string              `json:"type" db:"type"`
	Statut    string              `json:"statut" db:"statut"`
	Jetons    int                 `json:"jetons" db:"jetons"`
	Points    int                 `json:"points" db:"points"`
	CreatedAt time.Time           `json:"created_at" db:"created_at"`
	User      InteractionUser     `json:"user" db:"user"`
	Stand     InteractionStand    `json:"stand" db:"stand"`
	Kermesse  InteractionKermesse `json:"kermesse" db:"kermesse"`
}

type InteractionBasic struct {
	Id        int              `json:"id" db:"id"`
	Type      string           `json:"type" db:"type"`
	Statut    string           `json:"statut" db:"statut"`
	Jetons    int              `json:"jetons" db:"jetons"`
	Points    int              `json:"points" db:"points"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
	User      InteractionUser  `json:"user" db:"user"`
	Stand     InteractionStand `json:"stand" db:"stand"`
}
