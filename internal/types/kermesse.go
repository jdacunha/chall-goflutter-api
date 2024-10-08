package types

const (
	KermesseStatutStarted string = "STARTED"
	KermesseStatutEnded   string = "ENDED"
)

type Kermesse struct {
	Id          int    `json:"id" db:"id"`
	UserId      int    `json:"user_id" db:"user_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Statut      string `json:"statut" db:"statut"`
}
type KermesseStats struct {
	UserCount         int `json:"user_count"`
	StandCount        int `json:"stand_count"`
	InteractionCount  int `json:"interaction_count"`
	InteractionIncome int `json:"interaction_income"`
	TicketCount       int `json:"ticket_count"`
	TombolaIncome     int `json:"tombola_income"`
	PointsLadder      int `json:"points"`
}

type KermesseWithStats struct {
	Id                int    `json:"id" db:"id"`
	UserId            int    `json:"user_id" db:"user_id"`
	Name              string `json:"name" db:"name"`
	Description       string `json:"description" db:"description"`
	Statut            string `json:"statut" db:"statut"`
	UserCount         int    `json:"user_count"`
	StandCount        int    `json:"stand_count"`
	InteractionCount  int    `json:"interaction_count"`
	InteractionIncome int    `json:"interaction_income"`
	TicketCount       int    `json:"ticket_count"`
	TombolaIncome     int    `json:"tombola_income"`
	PointsLadder      int    `json:"points"`
}
