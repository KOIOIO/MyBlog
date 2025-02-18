package appTypes

type RoleID int

const (
	Guest RoleID = iota
	User
	Admin
)
