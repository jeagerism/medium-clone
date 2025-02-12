package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jeagerism/medium-clone/backend/internal/articles/entities"
	"github.com/jeagerism/medium-clone/backend/pkg/logger"
	"github.com/jmoiron/sqlx"
)

type articleRepository struct {
	db *sqlx.DB
}

func NewArticleRepository(db *sqlx.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) FindArticles(params entities.GetArticlesParams) ([]entities.ArticleResponse, error) {
	var articles []entities.ArticleResponse
	query := `
	SELECT 
		a.id, 
		a.title, 
		LEFT(a.content, 100) AS content, 
		a.user_id, 
		a.created_at, 
		a.updated_at,
		a.cover_image,
		COALESCE(STRING_AGG(DISTINCT t.name , ', '), '') AS tags, 
		COUNT(DISTINCT l.id) AS like_count,
		COUNT(DISTINCT c.id) AS comment_count
	FROM 
		articles a
	LEFT JOIN article_tags at2 ON a.id = at2.article_id
	LEFT JOIN tags t ON at2.tag_id = t.id
	LEFT JOIN likes l ON a.id = l.article_id
	LEFT JOIN comments c ON a.id = c.article_id
	WHERE 
		(COALESCE($1, '') = '' OR a.title ILIKE '%' || COALESCE($1, '') || '%')
		AND (COALESCE($2, '') = '' OR t.name ILIKE '%' || COALESCE($2, '') || '%')
	GROUP BY 
		a.id
	ORDER BY 
		a.id
	LIMIT $3 OFFSET $4;

	`

	err := r.db.Select(&articles, query, params.Search, params.Tags, params.Limit, params.Offset)
	if err != nil {
		logger.LogError(fmt.Errorf("database error in FindArticles: %w", err))
		return nil, fmt.Errorf("failed to find articles: %w", err)
	}

	return articles, nil
}

func (r *articleRepository) CountRow(params entities.GetArticlesParams) int {
	var total int
	countQuery := `
		SELECT COUNT(DISTINCT a.id)
		FROM 
			articles a
		LEFT JOIN article_tags at2 ON a.id = at2.article_id
		LEFT JOIN tags t ON at2.tag_id = t.id
		WHERE 
			($1 = '' OR a.title ILIKE '%' || $1 || '%')
			AND ($2 = '' OR t.name ILIKE '%' || $2 || '%')
	`
	err := r.db.Get(&total, countQuery, params.Search, params.Tags)
	if err != nil {
		logger.LogError(fmt.Errorf("database error in CountRow: %w", err))
		return 0
	}
	return total
}

func (r *articleRepository) FindByID(id int) (entities.ArticleResponse, error) {
	var article entities.ArticleResponse
	query := `
		SELECT 
			a.id, 
			a.title, 
			a.content, 
			a.user_id, 
			a.created_at, 
			a.updated_at,
			a.cover_image,
			STRING_AGG(DISTINCT t.name , ', ') AS tags, 
			COUNT(DISTINCT l.id) AS like_count,
			COUNT(DISTINCT c.id) AS comment_count
		FROM 
			articles a
		LEFT JOIN article_tags at2 ON a.id = at2.article_id
		LEFT JOIN tags t ON at2.tag_id = t.id
		LEFT JOIN likes l ON a.id = l.article_id
		LEFT JOIN comments c ON a.id = c.article_id
		WHERE 
			a.id = $1
		GROUP BY 
			a.id;
	`
	err := r.db.Get(&article, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			logger.LogError(fmt.Errorf("article with ID %d not found: %w", id, err))
			return entities.ArticleResponse{}, fmt.Errorf("article with ID %d not found: %w", id, err)
		}
		logger.LogError(fmt.Errorf("failed to find article by ID %d: %w", id, err))
		return entities.ArticleResponse{}, fmt.Errorf("failed to find article by ID %d: %w", id, err)
	}

	return article, nil
}

func (r *articleRepository) FindArticlesByUserID(req entities.GetArticlesByUserIDParams) ([]entities.ArticleResponse, error) {
	query := `
		SELECT 
			a.id, 
			a.title, 
			LEFT(a.content, 100) AS content, 
			a.user_id, 
			a.created_at, 
			a.updated_at,
			a.cover_image,
			COUNT(DISTINCT l.id) AS like_count,
			COUNT(DISTINCT c.id) AS comment_count
		FROM 
			articles a
		LEFT JOIN likes l ON a.id = l.article_id
		LEFT JOIN comments c ON a.id = c.article_id
		WHERE 
			a.user_id = $1
		GROUP BY 
			a.id
		ORDER BY 
			a.id
		LIMIT $2 OFFSET $3;
	`

	var articles []entities.ArticleResponse
	err := r.db.Select(&articles, query, req.ID, req.Limit, req.Offset)
	if err != nil {
		logger.LogError(fmt.Errorf("database error in FindArticlesByUserID: %w", err))
		return nil, fmt.Errorf("failed to find articles by user ID: %w", err)
	}

	return articles, nil
}

func (r *articleRepository) SaveArticle(req entities.AddArticleRequest) (int, error) {
	var id int
	query := `
		INSERT INTO articles (title, content, cover_image, user_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id;
	`
	err := r.db.QueryRow(query, req.Title, req.Content, req.Cover, req.UserID).Scan(&id)
	if err != nil {
		logger.LogError(fmt.Errorf("failed to save article: %w", err))
		return 0, fmt.Errorf("failed to save article: %w", err)
	}
	return id, nil
}

func (r *articleRepository) CheckTag(tag string) (int, error) {
	var id int
	query := `
		SELECT id
		FROM tags
		WHERE name = $1;
	`
	err := r.db.QueryRow(query, tag).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // Tag not found
		}
		logger.LogError(fmt.Errorf("failed to check tag '%s': %w", tag, err))
		return 0, fmt.Errorf("failed to check tag '%s': %w", tag, err)
	}
	return id, nil
}

func (r *articleRepository) SaveTag(tag string) (int, error) {
	var tagID int
	query := `INSERT INTO tags (name) VALUES ($1) RETURNING id;`
	err := r.db.QueryRow(query, tag).Scan(&tagID)
	if err != nil {
		logger.LogError(fmt.Errorf("failed to save tag '%s': %w", tag, err))
		return 0, fmt.Errorf("failed to save tag '%s': %w", tag, err)
	}
	return tagID, nil
}

func (r *articleRepository) SaveArticleTag(articleID, tagID int) error {
	query := `INSERT INTO article_tags (article_id, tag_id) VALUES ($1, $2);`
	_, err := r.db.Exec(query, articleID, tagID)
	if err != nil {
		logger.LogError(fmt.Errorf("failed to save article-tag association '%d' & '%d': %w", articleID, tagID, err))
		return fmt.Errorf("failed to associate tag %d with article %d: %w", tagID, articleID, err)
	}
	return nil
}

func (r *articleRepository) UpdateArticle(fields []string, args []interface{}, articleID int) error {
	setClauses := make([]string, len(fields))
	for i, field := range fields {
		setClauses[i] = fmt.Sprintf("%s = $%d", field, i+1)
	}

	query := fmt.Sprintf(`
		UPDATE articles
		SET %s
		WHERE id = $%d
	`, strings.Join(setClauses, ", "), len(fields)+1)

	args = append(args, articleID)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		logger.LogError(fmt.Errorf("failed to update article with ID %d: %w", articleID, err))
		return fmt.Errorf("failed to update article with ID %d: %w", articleID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.LogError(fmt.Errorf("failed to check rows affected for article ID %d: %w", articleID, err))
		return fmt.Errorf("failed to check rows affected for article ID %d: %w", articleID, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no article found with ID %d", articleID)
	}

	return nil
}

func (r *articleRepository) RemoveArticle(id int) error {
	query := `DELETE FROM articles WHERE id = $1;`
	result, err := r.db.Exec(query, id)
	if err != nil {
		logger.LogError(fmt.Errorf("failed to delete article with ID %d: %w", id, err))
		return fmt.Errorf("failed to delete article with ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.LogError(fmt.Errorf("failed to check rows affected for article ID %d: %w", id, err))
		return fmt.Errorf("failed to check rows affected for article ID %d: %w", id, err)
	}

	if rowsAffected == 0 {
		err := fmt.Errorf("no article found with ID %d", id)
		logger.LogError(err)
		return err
	}

	return nil
}

func (r *articleRepository) SaveComment(req entities.AddCommentRequest) error {
	query := `INSERT INTO comments (article_id,user_id,content) VALUES ($1,$2,$3)`
	_, err := r.db.Exec(query, req.ID, req.UserID, req.Content)
	if err != nil {
		logger.LogError(fmt.Errorf("failed to save comment %d , %d:%w", req.ID, req.UserID, err))
		return fmt.Errorf("failed to save comment %d , %d:%w", req.ID, req.UserID, err)
	}
	return nil
}

func (r *articleRepository) RemoveComment(id int) error {
	query := `DELETE FROM comments WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		logger.LogError(fmt.Errorf("failed to delete comment with ID %d: %w", id, err))
		return fmt.Errorf("failed to delete comment with ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		logger.LogError(fmt.Errorf("failed to check rows affected for comment ID %d: %w", id, err))
		return fmt.Errorf("failed to check rows affected for comment ID %d: %w", id, err)
	}

	if rowsAffected == 0 {
		err := fmt.Errorf("no comment found with ID %d", id)
		logger.LogError(err)
		return err
	}

	return nil
}

func (r *articleRepository) FindArticleComments(id int) ([]entities.GetArticleCommentsResponse, error) {
	var comments []entities.GetArticleCommentsResponse
	query := `
	SELECT  
		c.id,
		c.content,
		c.created_at,
		u.id AS user_id ,
		u.name,
		u.profile_image
	FROM 
		comments c
	LEFT JOIN users u ON c.user_id = u.id 
	WHERE article_id = $1`
	err := r.db.Select(&comments, query, id)
	if err != nil {
		logger.LogError(fmt.Errorf("database error in Find Article Comments: %w", err))
		return nil, fmt.Errorf("failed to find comment: %w", err)
	}

	return comments, nil
}

func (r *articleRepository) FindCommentByID(commentID int) (*entities.Comment, error) {
	var comment entities.Comment
	query := `SELECT id, user_id FROM comments WHERE id = $1`
	err := r.db.Get(&comment, query, commentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("comment not found")
		}
		return nil, fmt.Errorf("failed to retrieve comment: %w", err)
	}
	return &comment, nil
}
