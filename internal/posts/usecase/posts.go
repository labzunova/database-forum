package usecase

import (
	"DBproject/internal/posts"
	"DBproject/models"
)

type postsUsecase struct {
	postsRepository posts.PostsRepo
}

func NewPostsUsecase(repo posts.PostsRepo) posts.PostsUsecase {
	return &postsUsecase{
		postsRepository: repo,
	}
}

func (p postsUsecase) GetPost() ([]models.Post, models.Error) {
	return p.postsRepository.GetPost()
}

func (p postsUsecase) UpdatePost(id int64, post models.Post) (models.Post, models.Error) {
	return p.postsRepository.UpdatePost(id, post)
}

func (p postsUsecase) CreatePosts(slug string, posts []models.Post) ([]models.Post, models.Error) {
	return p.postsRepository.CreatePosts(slug, posts)
}
