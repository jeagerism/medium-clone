package repositories

import "time"

type article struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	UserID    int       `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Comment   string    `db:"comment"`    // Aggregated comment content
	Images    string    `db:"images"`     // Aggregated image URLs
	Tags      string    `db:"tags"`       // Aggregated tags
	LikeCount int       `db:"like_count"` // Total count of likes
}

type ArticleRepository interface {
	FindArticles() ([]article, error)
}
