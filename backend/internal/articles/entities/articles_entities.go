package entities

import "time"

type ArticleResponse struct {
	ID           int       `db:"id" json:"id"`
	Title        string    `db:"title" json:"title"`
	Content      string    `db:"content" json:"content"`
	UserID       int       `db:"user_id" json:"user_id"`
	Cover        string    `db:"cover_image" json:"cover_image"`
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

type GetArticlesByUserIDParams struct {
	ID     int    `form:"id"`
	Limit  int    `form:"limit"`
	Page   int    `form:"page"`
	Sort   string `form:"sort"`
	Offset int
}

type AddArticleRequest struct {
	Title   string `json:"title" binding:"required,min=3,max=100"` // Title is required, with a length between 3 and 100
	Content string `json:"content" binding:"required,min=10"`      // Content is required, minimum 10 characters
	UserID  int    `json:"user_id" binding:"required,gt=0"`        // UserID is required and must be greater than 0
	Tags    string `json:"tags" binding:"required"`                // Tags are required (can also use custom validation for format)
	Cover   string `json:"cover_image" binding:"required,url"`     // Cover image URL is required and must be a valid URL
}

type UpdateArticleRequest struct {
	Id      int     `json:"id" binding:"required"`
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Cover   *string `json:"cover_image"`
}

type DeleteArticleRequest struct {
	ID     int `json:"id" binding:"required"`
	UserID int `json:"user_id" binding:"required"`
}

type AddCommentRequest struct {
	ID      int    `json:"article_id" binding:"required"`
	UserID  int    `json:"user_id" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type DeleteCommentRequest struct {
	ID     int `json:"article_id" binding:"required"`
	UserID int `json:"user_id" binding:"required"`
}

type GetArticleCommentsResponse struct {
	ID         int       `db:"id" json:"comment_id"`
	AID        int       `db:"article_id" json:"article_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
	Content    string    `db:"content" json:"content"`
	UserID     int       `db:"user_id" json:"user_id"`
	UserName   string    `db:"name" json:"name"`
	ProfileImg string    `db:"profile_image" json:"profile_image"`
}
