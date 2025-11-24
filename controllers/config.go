package controllers

import (
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/pkg/config"
	"github.com/sanjayk-eng/UserMenagmentSystem_Backend/repositories"
)

// HandlerFunc holds dependencies
type HandlerFunc struct {
	Env   *config.ENV
	Query *repositories.Repository
}

// NewHandler initializes and returns a HandlerFunc
func NewHandler(env *config.ENV, query *repositories.Repository) *HandlerFunc {
	return &HandlerFunc{
		Env:   env,
		Query: query,
	}
}
