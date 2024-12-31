package services

import (
	"errors"

	"github.com/jeagerism/medium-clone/backend/internal/articles/repositories"
)

type articleService struct {
	articleRepository repositories.ArticleRepository
}

func NewArticleService(articleRepository repositories.ArticleRepository) ArticleService {
	return &articleService{articleRepository: articleRepository}
}

func (s *articleService) GetArticles() ([]article, error) {
	// เรียกข้อมูลจาก repository
	articleRepo, err := s.articleRepository.FindArticles()
	if err != nil {
		return nil, errors.New("failed to get articles")
	}

	// สร้าง array ของ article ที่จะส่งกลับ
	var articles []article

	// แปลงข้อมูลจาก repository ไปเป็น struct ที่ใช้งานใน service
	for _, v := range articleRepo {
		articles = append(articles, article{
			ID:        v.ID,
			Title:     v.Title,
			Content:   v.Content,
			UserID:    v.UserID,
			CreatedAt: v.CreatedAt,
			UpdatedAt: v.UpdatedAt,
			Comment:   v.Comment, // ถ้ามีการรวม comment ไว้ใน repository
			Images:    v.Images,  // ถ้ามีการรวม images ไว้ใน repository
			Tags:      v.Tags,    // ถ้ามีการรวม tags ไว้ใน repository
			LikeCount: v.LikeCount,
		})
	}

	return articles, nil
}
