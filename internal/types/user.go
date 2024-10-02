package types

type contextKey string

const (
	UserIDKey   contextKey = "userId"
	UserRoleKey contextKey = "userRole"
)

const (
	UserRoleOrganisateur string = "ORGANISATEUR"
	UserRoleTeneurStand  string = "TENEUR_STAND"
	UserRoleParent       string = "PARENT"
	UserRoleEnfant       string = "ENFANT"
)

type User struct {
	Id           int    `json:"id" db:"id"`
	ParentId     *int   `json:"parentId" db:"parent_id"`
	Name         string `json:"name" db:"name"`
	Email        string `json:"email" db:"email"`
	PasswordHash string `json:"password" db:"password_hash"`
	Role         string `json:"role" db:"role"`
	Jetons       int    `json:"jetons" db:"jetons"`
}

type UserBasic struct {
	Id     int    `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Email  string `json:"email" db:"email"`
	Role   string `json:"role" db:"role"`
	Jetons int    `json:"jetons" db:"jetons"`
}

type UserBasicWithToken struct {
	Id     int    `json:"id" db:"id"`
	Name   string `json:"name" db:"name"`
	Email  string `json:"email" db:"email"`
	Role   string `json:"role" db:"role"`
	Jetons int    `json:"jetons" db:"jetons"`
	Token  string `json:"token"`
}
