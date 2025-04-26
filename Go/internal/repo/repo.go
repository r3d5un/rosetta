package repo

import "github.com/r3d5un/rosetta/Go/internal/data"

type Repository struct {
	models *data.Models
}

func NewRepository(models *data.Models) Repository {
	return Repository{models: models}
}
