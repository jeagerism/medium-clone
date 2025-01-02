package services

import (
	"strings"

	"github.com/jeagerism/medium-clone/backend/internal/articles/repositories"
	"github.com/jeagerism/medium-clone/backend/internal/entities"
	"github.com/jeagerism/medium-clone/backend/pkg/utils"
)

type articleService struct {
	articleRepository repositories.ArticleRepository
}

func NewArticleService(articleRepository repositories.ArticleRepository) ArticleService {
	return &articleService{articleRepository: articleRepository}
}

func (s *articleService) GetArticles(params entities.GetArticlesParams) (getAllResponse, error) {
	params.Offset = utils.CalculateOffset(params.Page, params.Limit)
	// Call the repository to fetch articles
	articleRepo, err := s.articleRepository.FindArticles(params)
	if err != nil {
		return getAllResponse{}, err
	}

	var count int
	if len(articleRepo) == 0 {
		count = 0
	} else {
		count = s.articleRepository.CountRow(params)
	}
	// Return the response structure
	return getAllResponse{
		Count: count,
		Data:  articleRepo,
	}, nil
}

func (s *articleService) GetArticleByID(id int) (getArticleByIDResponse, error) {
	var article getArticleByIDResponse
	articleRepo, err := s.articleRepository.FindByID(id)
	if err != nil {
		return article, err
	}

	article = getArticleByIDResponse{
		ID:           articleRepo.ID,
		Title:        articleRepo.Title,
		Content:      articleRepo.Content,
		UserID:       articleRepo.UserID,
		CreatedAt:    articleRepo.CreatedAt,
		UpdatedAt:    articleRepo.UpdatedAt,
		Tags:         strings.Split(articleRepo.Tags, ", "),
		LikeCount:    articleRepo.LikeCount,
		CommentCount: articleRepo.CommentCount,
	}

	return article, nil
}
