package repositories

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type articleRepository struct {
	db *sqlx.DB
}

func NewArticleRepository(db *sqlx.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) FindArticles(search string, tags []string, limit, offset int) ([]article, int, error) {
	var articles []article
	var total int

	// Dynamically construct the WHERE condition for tags
	var tagConditions string
	if len(tags) > 0 {
		tagConditions = " AND ("
		for i, tag := range tags {
			if i > 0 {
				tagConditions += " OR"
			}
			tagConditions += fmt.Sprintf(" t.name ILIKE '%%%s%%'", tag)
		}
		tagConditions += ")"
	}

	// Prepare the query
	query := fmt.Sprintf(`
		SELECT 
			a.id, 
			a.title, 
			a.content, 
			a.user_id, 
			a.created_at, 
			a.updated_at,
			STRING_AGG(DISTINCT c."content", ', ') AS comments,  -- Combine unique comments
			STRING_AGG(DISTINCT i.image_url, ', ') AS images,   -- Combine unique image URLs
			STRING_AGG(DISTINCT t."name", ', ') AS tags, 
			COUNT(DISTINCT l.id) AS like_count
		FROM 
			articles a
		LEFT JOIN comments c ON a.id = c.article_id 
		LEFT JOIN images i ON a.id = i.article_id 
		LEFT JOIN article_tags at2 ON a.id = at2.article_id 
		LEFT JOIN tags t ON at2.tag_id = t.id
		LEFT JOIN likes l ON a.id = l.article_id 
		WHERE 
			($1 = '' OR a.title ILIKE $1 OR a.content ILIKE $1)
		%s
		GROUP BY 
			a.id
		ORDER BY 
			a.id
		LIMIT %d OFFSET %d;
	`, tagConditions, limit, offset)

	// Execute the query and scan directly into the article struct
	err := r.db.Select(&articles, query, search)
	if err != nil {
		return nil, 0, err
	}

	// No need to modify the result after scanning as pq.Array() is not required
	// The query will return the appropriate arrays for comments, images, and tags

	// Count query
	countQuery := fmt.Sprintf(`
		SELECT COUNT(DISTINCT a.id)
		FROM 
			articles a
		LEFT JOIN article_tags at2 ON a.id = at2.article_id 
		LEFT JOIN tags t ON at2.tag_id = t.id
		WHERE 
			($1 = '' OR a.title ILIKE $1 OR a.content ILIKE $1)
		%s
	`, tagConditions)

	err = r.db.Get(&total, countQuery, search)
	if err != nil {
		return nil, 0, err
	}

	return articles, total, nil
}
