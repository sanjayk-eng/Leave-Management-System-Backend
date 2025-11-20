package controllers

import (
	"github.com/jmoiron/sqlx"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/pkg/config"
)

// HandlerFunc holds dependencies
type HandlerFunc struct {
	Env *config.ENV
	DB  *sqlx.DB
}

// NewHandler initializes and returns a HandlerFunc
func NewHandler(env *config.ENV, db *sqlx.DB) *HandlerFunc {
	return &HandlerFunc{
		Env: env,
		DB:  db,
	}
}
