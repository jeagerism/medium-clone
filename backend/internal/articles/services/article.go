package services

import (
	"fmt"
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
		Cover:        article.Cover,
		CreatedAt:    articleRepo.CreatedAt,
		UpdatedAt:    articleRepo.UpdatedAt,
		Tags:         strings.Split(articleRepo.Tags, ", "),
		LikeCount:    articleRepo.LikeCount,
		CommentCount: articleRepo.CommentCount,
	}

	return article, nil
}

func (s *articleService) AddArticle(req entities.AddArticleRequest) error {
	// Save the article and get its ID
	arId, err := s.articleRepository.SaveArticle(req)
	if err != nil {
		return err
	}

	// Split the tags from the request
	tags := strings.Split(req.Tags, ",")
	for _, tag := range tags {
		// Clean up each tag (trim whitespace)
		tag = strings.TrimSpace(tag)
		if tag == "" {
			continue // Skip empty tags
		}

		var tagId int

		// Check if the tag exists
		tagId, err = s.articleRepository.CheckTag(tag)
		if err != nil {
			return err
		}

		// If the tag doesn't exist, save the tag and get its ID
		if tagId == 0 {
			tagId, err = s.articleRepository.SaveTag(tag)
			if err != nil {
				return err
			}
		}

		// Save the article-tag association
		err = s.articleRepository.SaveArticleTag(arId, tagId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *articleService) UpdateArticle(req entities.UpdateArticleRequest) error {
	// เก็บคอลัมน์ที่จะอัปเดต
	var fields []string
	var args []interface{}
	argID := 1

	// ตรวจสอบและเพิ่มข้อมูลที่ต้องการอัปเดต
	if req.Title != nil {
		fields = append(fields, fmt.Sprintf("title = $%d", argID))
		args = append(args, *req.Title)
		argID++
	}

	if req.Content != nil {
		fields = append(fields, fmt.Sprintf("content = $%d", argID))
		args = append(args, *req.Content)
		argID++
	}

	if req.Cover != nil {
		fields = append(fields, fmt.Sprintf("cover_image = $%d", argID))
		args = append(args, *req.Cover)
		argID++
	}

	// ถ้าไม่มีฟิลด์ใดถูกอัปเดต
	if len(fields) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// ส่งข้อมูลไปยัง repository
	return s.articleRepository.UpdateArticle(fields, args, req.Id)
}
