package repo

import "github.com/r3d5un/rosetta/Go/internal/data"

type Repository struct {
	models       *data.Models
	ForumReader  ForumReader
	ForumWriter  ForumWriter
	ThreadReader ThreadReader
	ThreadWriter ThreadWriter
	PostReader   PostReader
	PostWriter   PostWriter
	UserReader   UserReader
	UserWriter   UserWriter
}

func NewRepository(models *data.Models) Repository {
	userRepo := NewUserRepository(models)
	forumRepo := NewForumRepository(models, &userRepo)
	threadRepo := NewThreadRepository(models, &forumRepo, &userRepo)
	postRepo := NewPostRepository(models, &threadRepo, &userRepo)

	return Repository{
		models:       models,
		ForumReader:  &forumRepo,
		ForumWriter:  &forumRepo,
		ThreadReader: &threadRepo,
		ThreadWriter: &threadRepo,
		PostReader:   &postRepo,
		PostWriter:   &postRepo,
		UserReader:   &userRepo,
		UserWriter:   &userRepo,
	}
}
