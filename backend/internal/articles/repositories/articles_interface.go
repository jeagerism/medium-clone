package repositories

import (
	"github.com/jeagerism/medium-clone/backend/internal/entities"
)

type ArticleRepository interface {
	FindArticles(params entities.GetArticlesParams) ([]entities.ArticleResponse, error)
	CountRow(params entities.GetArticlesParams) int
	FindByID(id int) (entities.ArticleResponse, error)
	SaveArticle(req entities.AddArticleRequest) (int, error)
	CheckTag(tag string) (int, error)
	SaveTag(tag string) (int, error)
	SaveArticleTag(articleID, tagID int) error
	UpdateArticle(fields []string, args []interface{}, articleID int) error
}
