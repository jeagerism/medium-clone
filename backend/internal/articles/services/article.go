package services

import (
	"fmt"
	"strings" // Add strings package for splitting

	"github.com/jeagerism/medium-clone/backend/internal/articles/repositories"
)

type articleService struct {
	articleRepository repositories.ArticleRepository
}

func NewArticleService(articleRepository repositories.ArticleRepository) ArticleService {
	return &articleService{articleRepository: articleRepository}
}

func (s *articleService) GetArticles() ([]article, error) {
	// Call the repository to fetch articles
	articleRepo, count, err := s.articleRepository.FindArticles("", []string{"Science", "Tech"}, 5, 0)
	if err != nil {
		return nil, err
	}
	fmt.Println(count)

	// Create an array of article that will be returned
	var articles []article

	// Process each article and convert Tags and Images fields to arrays (slices)
	for _, v := range articleRepo {
		articles = append(articles, article{
			ID:        v.ID,
			Title:     v.Title,
			Content:   v.Content,
			UserID:    v.UserID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Comments:  strings.Split(v.Comments, ", "), // Include the combined comments if applicable
			// Convert Images string into a slice
			Images: strings.Split(v.Images, ", "), // Split the Images string by comma
			// Convert Tags string into a slice
			Tags:      strings.Split(v.Tags, ", "), // Split the Tags string by comma
			LikeCount: v.LikeCount,
		})
	}

	return articles, nil
}
