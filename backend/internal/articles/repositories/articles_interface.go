package repositories

import "time"

type article struct {
	ID        int       `db:"id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	UserID    int       `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Comments  string    `db:"comments"` // Store comments as a slice of strings
	Images    string    `db:"images"`   // Store images as a slice of strings
	Tags      string    `db:"tags"`     // Store tags as a slice of strings
	LikeCount int       `db:"like_count"`
}

type ArticleRepository interface {
	FindArticles(search string, tags []string, limit, offset int) ([]article, int, error)
}
