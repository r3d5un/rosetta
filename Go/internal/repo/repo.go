package repo

import "github.com/r3d5un/rosetta/Go/internal/data"

type Repository struct {
	models      *data.Models
	ForumReader ForumReader
	ForumWriter ForumWriter
}

func NewRepository(models *data.Models) Repository {
	forumRepo := NewForumRepository(models)

	return Repository{
		models:      models,
		ForumReader: &forumRepo,
		ForumWriter: &forumRepo,
	}
}
