package types

type IMember interface {
	Roles() []IRole
}

type IRole interface{}
