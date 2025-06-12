package domain

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

func (r Role) String() string {
	return string(r)
}

func IsValidRole(r string) bool {
	return r == RoleUser.String() || r == RoleAdmin.String()
}
