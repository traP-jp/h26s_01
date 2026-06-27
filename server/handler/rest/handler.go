package rest

import "github.com/traP-jp/h26s_01/server/repository"

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}
