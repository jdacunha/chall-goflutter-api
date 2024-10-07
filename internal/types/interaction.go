package types

const (
	InteractionTypeTransaction string = "TRANSACTION"
	InteractionTypeActivite    string = "ACTIVITE"
	InteractionStatutStarted   string = "STARTED"
	InteractionStatutEnded     string = "ENDED"
)

type Interaction struct {
	Id         int    `json:"id" db:"id"`
	UserId     int    `json:"user_id" db:"user_id"`
	KermesseId int    `json:"kermesse_id" db:"kermesse_id"`
	StandId    int    `json:"stand_id" db:"stand_id"`
	Type       string `json:"type" db:"type"`
	Status     string `json:"status" db:"status"`
	Jetons     int    `json:"jetons" db:"jetons"`
	Points     int    `json:"points" db:"points"`
}
