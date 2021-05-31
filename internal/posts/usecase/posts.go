package usecase

import (
	"DBproject/internal/posts"
	"DBproject/models"
	"strconv"
)

type postsUsecase struct {
	postsRepository posts.PostsRepo
}

func NewPostsUsecase(repo posts.PostsRepo) posts.PostsUsecase {
	return &postsUsecase{
		postsRepository: repo,
	}
}

func (p postsUsecase) GetPost() (models.Post, models.Error) {
	return p.postsRepository.GetPost()
}

func (p postsUsecase) UpdatePost(id int,  message string) (models.Post, models.Error) {
	return p.postsRepository.UpdatePost(id, message)
}

func (p postsUsecase) CreatePosts(slug string, posts []models.Post) ([]models.Post, models.Error) {
	var thread models.Thread
	var err models.Error
	id, errID := strconv.Atoi(slug)
	if errID != nil {
		thread, err = p.postsRepository.GetThreadAndForumById(id)
		if err.Code == 404 {
			return nil, err
		}
	} else {
		thread, err = p.postsRepository.GetThreadAndForumBySlug(slug)
		if err.Code == 404 {
			return nil, err
		}
	}

	return p.postsRepository.CreatePosts(thread, posts)
}
