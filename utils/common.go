package utils

import "github.com/google/uuid"

type Common struct {
	Component  string
	Action     string
	FromUserID uuid.UUID
}

func NewCommon(component, action string, fromUserID uuid.UUID) *Common {
	return &Common{
		Component:  component,
		Action:     action,
		FromUserID: fromUserID,
	}
}
