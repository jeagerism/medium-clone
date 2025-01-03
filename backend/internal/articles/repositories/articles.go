package repositories

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jeagerism/medium-clone/backend/internal/entities"
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

	// ปรับ query ให้รองรับกรณี tags ว่างหรือไม่
	query := `
		SELECT 
			a.id, 
			a.title, 
			LEFT(a.content, 100) AS content, -- จำกัดเนื้อหา content ไว้ 100 ตัวอักษร
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
			($1 = '' OR a.title ILIKE '%' || $1 || '%')
			AND ($2 = '' OR t.name ILIKE '%' || $2 || '%')
		GROUP BY 
			a.id
		ORDER BY 
			a.id
		LIMIT $3 OFFSET $4;
	`

	err := r.db.Select(&articles, query, params.Search, params.Tags, params.Limit, params.Offset)
	if err != nil {
		return nil, err
	}

	return articles, nil
}

func (r *articleRepository) CountRow(params entities.GetArticlesParams) int {
	var total int
	// Query for total count
	countQuery := `
		SELECT COUNT(DISTINCT a.id)
		FROM 
			articles a
		LEFT JOIN article_tags at2 ON a.id = at2.article_id
		LEFT JOIN tags t ON at2.tag_id = t.id
		WHERE 
			($1 = '' OR a.title ILIKE '%' || $1 || '%')
			-- ตรวจสอบเงื่อนไขสำหรับ tags
			 AND ($2 = '' OR t.name ILIKE '%' || $2 || '%') -- การใช้ ILIKE สำหรับการค้นหาด้วย pattern
	`
	err := r.db.Get(&total, countQuery, params.Search, params.Tags)
	if err != nil {
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
		a.content, -- ดึง content ทั้งหมด
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
		return entities.ArticleResponse{}, err
	}

	return article, nil
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
		return 0, err
	}
	return id, nil
}

func (r *articleRepository) CheckTag(tag string) (int, error) {
	var id int
	query := `
	SELECT id
	FROM tags
	WHERE name = $1
	`
	err := r.db.QueryRow(query, tag).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // Tag not found
		}
		return 0, err // Other errors
	}
	return id, nil
}

func (r *articleRepository) SaveTag(tag string) (int, error) {
	var tagID int
	query := `INSERT INTO tags (name) VALUES ($1) RETURNING id`
	err := r.db.QueryRow(query, tag).Scan(&tagID)
	if err != nil {
		return 0, err
	}
	return tagID, nil
}

func (r *articleRepository) SaveArticleTag(articleID, tagID int) error {
	query := `INSERT INTO article_tags (article_id, tag_id) VALUES ($1, $2)`
	_, err := r.db.Exec(query, articleID, tagID)
	if err != nil {
		return err
	}
	return nil
}

// func (r *articleRepository) UpdateArticle(req entities.UpdateArticleRequest) error {
// 	query := `
// 	UPDATE articles
// 	SET title = $1, content = $2, cover_image = $3
// 	WHERE id = $4;`
// 	result, err := r.db.Exec(query, req.Title, req.Content, req.Cover, req.Id)
// 	if err != nil {
// 		return err
// 	}

// 	// ตรวจสอบจำนวนแถวที่ได้รับผลกระทบ
// 	rowsAffected, err := result.RowsAffected()
// 	if err != nil {
// 		return err
// 	}

// 	if rowsAffected == 0 {
// 		return fmt.Errorf("no article found with id %d", req.Id)
// 	}

// 	return nil
// }

func (r *articleRepository) UpdateArticle(fields []string, args []interface{}, articleID int) error {
	// สร้าง query สำหรับอัปเดต
	query := fmt.Sprintf(`
		UPDATE articles
		SET %s
		WHERE id = $%d
	`, strings.Join(fields, ", "), len(fields)+1)
	args = append(args, articleID)

	// รัน query
	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	// ตรวจสอบจำนวนแถวที่ได้รับผลกระทบ
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	// หากไม่พบข้อมูลที่จะอัปเดต
	if rowsAffected == 0 {
		return fmt.Errorf("no article found with id %d", articleID)
	}

	return nil
}
