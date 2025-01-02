package services

import (
	"time"

	"github.com/jeagerism/medium-clone/backend/internal/entities"
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
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updeated_at"`
	Tags         []string  `json:"tags"`
	LikeCount    int       `json:"like_count"`
	CommentCount int       `json:"comment_count"` // Add the comment count fieldF
}

type ArticleService interface {
	GetArticles(params entities.GetArticlesParams) (getAllResponse, error)
	GetArticleByID(id int) (getArticleByIDResponse, error)
}
