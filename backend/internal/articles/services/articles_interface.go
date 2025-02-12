package services

import (
	"time"

	"github.com/jeagerism/medium-clone/backend/internal/articles/entities"
)

type getAllResponse struct {
	Count int                        `json:"count"`
	Data  []entities.ArticleResponse `json:"data"`
}

type getArticleByIDResponse struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	UserID       int       `json:"user_id"`
	Cover        string    `json:"cover_image"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updeated_at"`
	Tags         []string  `json:"tags"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"` // Add the comment count fieldF
}

type ArticleService interface {
	// Articles Services
	GetArticles(params entities.GetArticlesParams) (getAllResponse, error)
	GetArticleByID(id int) (getArticleByIDResponse, error)
	AddArticle(req entities.AddArticleRequest) error
	UpdateArticle(req entities.UpdateArticleRequest) error
	DeleteArticle(id int) error
	GetArticleByUserID(req entities.GetArticlesByUserIDParams) ([]entities.ArticleResponse, error)

	// Comment Services
	AddComment(req entities.AddCommentRequest) error
	DeleteComment(id int) error
	GetArticleComments(id int) ([]entities.GetArticleCommentsResponse, error)
	GetCommentByID(commentID int) (*entities.Comment, error)
}
