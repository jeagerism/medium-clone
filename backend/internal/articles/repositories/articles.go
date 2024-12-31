package repositories

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type articleRepository struct {
	db *sqlx.DB
}

func NewArticleRepository(db *sqlx.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) FindArticles() ([]article, error) {
	query := `
	SELECT 
		a.id, 
		a.title, 
		a.content, 
		a.user_id, 
		a.created_at, 
		a.updated_at,
		STRING_AGG(DISTINCT c."content", ', ') AS comment,
		STRING_AGG(DISTINCT i.image_url, ', ') AS images,
		STRING_AGG(DISTINCT t."name", ', ') AS tags,
		COUNT(DISTINCT l.id) AS like_count
	FROM 
		articles a
	LEFT JOIN comments c ON a.id = c.article_id 
	LEFT JOIN images i ON a.id = i.article_id 
	LEFT JOIN article_tags at2 ON a.id = at2.article_id 
	LEFT JOIN tags t ON at2.tag_id = t.id
	LEFT JOIN likes l ON a.id = l.article_id 
	GROUP BY 
		a.id, a.title, a.content, a.user_id, a.created_at, a.updated_at
	ORDER BY 
		a.id;`

	var articles []article
	err := r.db.Select(&articles, query)
	if err != nil {
		log.Fatalln("Database Query Error")
		return nil, err
	}
	return articles, nil
}
