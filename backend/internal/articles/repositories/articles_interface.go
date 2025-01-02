package repositories

import (
	"github.com/jeagerism/medium-clone/backend/internal/entities"
)

type ArticleRepository interface {
	FindArticles(params entities.GetArticlesParams) ([]entities.ArticleResponse, error)
	CountRow(params entities.GetArticlesParams) int
	FindByID(id int) (entities.ArticleResponse, error)
}
