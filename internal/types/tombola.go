package types

const (
	TombolaStatutStarted = "STARTED"
	TombolaStatutEnded   = "ENDED"
)

type Tombola struct {
	Id         int    `json:"id" db:"id"`
	KermesseId int    `json:"kermesse_id" db:"kermesse_id"`
	Name       string `json:"name" db:"name"`
	Statut     string `json:"statut" db:"statut"`
	Price      int    `json:"price" db:"price"`
	Lot        string `json:"lot" db:"lot"`
}
