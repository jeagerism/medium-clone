package repositories

import (
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

// ถ้า content ข้อมูลเยอะมาก ต้องการจำกัดข้อมูลในการ FindArticles แต่ใน FindById เอาเต็มที่มี
