package store

import "github.com/aserawayneasera/research-tracker/internal/models"

type ListParams struct {
	Status string
	Year   int
	Q      string
	Limit  int
	Offset int
}

type PaperStore interface {
	Create(p *models.Paper) (*models.Paper, error)
	Get(id int64) (*models.Paper, error)
	Update(id int64, p *models.Paper) (*models.Paper, error)
	Delete(id int64) error
	List(params ListParams) ([]models.Paper, int, error)
}
