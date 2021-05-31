package usecase

import (
	"DBproject/internal/forum"
	"DBproject/models"
)

type forumUsecase struct {
	forumRepository forum.ForumRepo
}

func NewForumUsecase(repo forum.ForumRepo) forum.ForumUsecase {
	return &forumUsecase{
		forumRepository: repo,
	}
}

func (f forumUsecase) CreateNewForum(forum models.Forum) (models.Forum, models.Error) {
	return f.forumRepository.CreateNewForum(forum)
}

func (f forumUsecase) GetForum(id string) (models.Forum, models.Error) {
	return f.forumRepository.GetForum(id)
}

func (f forumUsecase) CreateThread(slug string, thread models.Thread) (models.Thread, models.Error) {
	return f.forumRepository.CreateThread(slug, thread)
}

func (f forumUsecase) GetUsers(slug string, params models.ParseParams) ([]models.User, models.Error) {
	return f.forumRepository.GetUsers(slug, params)
}

func (f forumUsecase) GetThreads(slug string, params models.ParseParams) ([]models.Thread, models.Error) {
	return f.forumRepository.GetThreads(slug, params)
}


