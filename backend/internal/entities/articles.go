package entities

import "time"

type ArticleResponse struct {
	ID           int       `db:"id" json:"id"`
	Title        string    `db:"title" json:"title"`
	Content      string    `db:"content" json:"content"`
	UserID       int       `db:"user_id" json:"user_id"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updeated_at"`
	Tags         string    `db:"tags" json:"tags"`
	LikeCount    int       `db:"like_count" json:"like_count"`
	CommentCount int       `db:"comment_count" json:"comment_count"` // Add the comment count field
}

type GetArticlesParams struct {
	Search string `form:"search"` // ดึงค่าจาก query string
	Tags   string `form:"tags"`
	Limit  int    `form:"limit"`
	Page   int    `form:"page"`
	Sort   string `form:"sort"`
	Offset int
}
