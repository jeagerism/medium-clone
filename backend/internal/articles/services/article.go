package services

import (
	"fmt"
	"strings"

	"github.com/jeagerism/medium-clone/backend/internal/articles/entities"
	"github.com/jeagerism/medium-clone/backend/internal/articles/repositories"
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
		return getAllResponse{}, ErrArticlesNotFound
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
		if strings.Contains(err.Error(), "not found") {
			return getArticleByIDResponse{}, ErrArticleNotFound
		}
		return getArticleByIDResponse{}, ErrFailedToFind
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

func (s *articleService) GetArticleByUserID(req entities.GetArticlesByUserIDParams) ([]entities.ArticleResponse, error) {
	req.Offset = utils.CalculateOffset(req.Page, req.Limit)

	articles, err := s.articleRepository.FindArticlesByUserID(req)
	if err != nil {
		return []entities.ArticleResponse{}, ErrArticlesNotFound
	}
	return articles, nil
}

func (s *articleService) AddArticle(req entities.AddArticleRequest) error {
	// Save the article and get its ID
	arId, err := s.articleRepository.SaveArticle(req)
	if err != nil {
		return ErrFailedToAddArticle
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
			return ErrFailedToAddArticle
		}

		// If the tag doesn't exist, save the tag and get its ID
		if tagId == 0 {
			tagId, err = s.articleRepository.SaveTag(tag)
			if err != nil {
				return ErrFailedToAddTag
			}
		}

		// Save the article-tag association
		err = s.articleRepository.SaveArticleTag(arId, tagId)
		if err != nil {
			return ErrFailedToAssociateTag
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
		return ErrNoFieldsUpdate
	}

	// ส่งข้อมูลไปยัง repository
	err := s.articleRepository.UpdateArticle(fields, args, req.Id)
	if err != nil {
		if strings.Contains(err.Error(), "no article found") {
			return ErrArticleNotFound
		}
		return ErrFailedToUpdate
	}

	return nil
}

func (s *articleService) DeleteArticle(id int) error {
	err := s.articleRepository.RemoveArticle(id)
	if err != nil {
		if strings.Contains(err.Error(), "no article found") {
			// แปลงเป็น error ที่เหมาะสมกับ business logic
			return ErrArticleNotFound
		}

		// error อื่นๆ ส่งไปที่ Handler/Controller
		return ErrFailedToDelete
	}

	return nil
}

// comments Services
func (s *articleService) AddComment(req entities.AddCommentRequest) error {
	err := s.articleRepository.SaveComment(req)
	if err != nil {
		return ErrFailedToAddComment
	}
	return nil
}

func (s *articleService) DeleteComment(id int) error {
	err := s.articleRepository.RemoveComment(id)
	if err != nil {
		if strings.Contains(err.Error(), "no comment found") {
			// แปลงเป็น error ที่เหมาะสมกับ business logic
			return ErrCommentNotFound
		}

		// error อื่นๆ ส่งไปที่ Handler/Controller
		return ErrFailedToDeleteComment
	}

	return nil
}

func (s *articleService) GetArticleComments(id int) ([]entities.GetArticleCommentsResponse, error) {
	comments, err := s.articleRepository.FindArticleComments(id)
	if err != nil {
		return comments, ErrCommentNotFound
	}
	return comments, nil
}

func (s *articleService) GetCommentByID(commentID int) (*entities.Comment, error) {
	comment, err := s.articleRepository.FindCommentByID(commentID)
	if err != nil {
		return nil, fmt.Errorf("comment not found: %w", err)
	}
	return comment, nil
}
