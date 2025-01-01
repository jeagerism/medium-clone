package services

import "time"

type article struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Comments  []string  `json:"comments"`   // Aggregated comment content
	Images    []string  `json:"images"`     // Aggregated image URLs
	Tags      []string  `json:"tags"`       // Aggregated tags
	LikeCount int       `json:"like_count"` // Total count of likes
}

type ArticleService interface {
	GetArticles() ([]article, error)
}
