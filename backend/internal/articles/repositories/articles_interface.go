package repositories

import (
	"github.com/jeagerism/medium-clone/backend/internal/articles/entities"
)

type ArticleRepository interface {
	FindArticles(params entities.GetArticlesParams) ([]entities.ArticleResponse, error)
	CountRow(params entities.GetArticlesParams) int
	FindByID(id int) (entities.ArticleResponse, error)
	FindArticlesByUserID(req entities.GetArticlesByUserIDParams) ([]entities.ArticleResponse, error)
	SaveArticle(req entities.AddArticleRequest) (int, error)
	CheckTag(tag string) (int, error)
	SaveTag(tag string) (int, error)
	SaveArticleTag(articleID, tagID int) error
	UpdateArticle(fields []string, args []interface{}, articleID int) error
	RemoveArticle(id int) error
	SaveComment(req entities.AddCommentRequest) error
	RemoveComment(id int) error
	FindArticleComments(id int) ([]entities.GetArticleCommentsResponse, error)
	FindCommentByID(commentID int) (*entities.Comment, error)
}
