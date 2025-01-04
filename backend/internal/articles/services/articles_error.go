package services

import "errors"

// Article-related errors
var (
	ErrNoFieldsUpdate     = errors.New("no fields to update")
	ErrArticleNotFound    = errors.New("article not found")
	ErrArticlesNotFound   = errors.New("articles not found")
	ErrFailedToUpdate     = errors.New("failed to update article")
	ErrFailedToAddArticle = errors.New("failed to add new article")
	ErrFailedToFind       = errors.New("failed to find article")
	ErrFailedToDelete     = errors.New("failed to delete article")
)

// Tag-related errors
var (
	ErrCheckTag             = errors.New("failed to check tag existence")
	ErrFailedToAddTag       = errors.New("failed to add new tag")
	ErrFailedToAssociateTag = errors.New("failed to associate tag with article")
)
