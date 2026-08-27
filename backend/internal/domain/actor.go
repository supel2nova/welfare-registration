package domain

import "github.com/google/uuid"

type ActorType string
type Role string

const (
	ActorUser   ActorType = "USER"
	ActorSystem ActorType = "SYSTEM"
)

const (
	RoleRegistrar Role = "REGISTRAR"
	RoleReviewer  Role = "REVIEWER"
	RoleAdmin     Role = "ADMIN"
)

type Actor struct {
	Type     ActorType
	OrgID    uuid.UUID
	UserID   *uuid.UUID
	ClientID *uuid.UUID
	Role     Role
}

func (a Actor) IsSystem() bool { return a.Type == ActorSystem }

func (a Actor) ActorID() uuid.UUID {
	if a.IsSystem() {
		if a.ClientID != nil {
			return *a.ClientID
		}
		return uuid.Nil
	}
	if a.UserID != nil {
		return *a.UserID
	}
	return uuid.Nil
}
