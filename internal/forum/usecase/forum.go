package usecase

import (
	"DBproject/internal/forum"
	"DBproject/models"
	"github.com/google/uuid"
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
	checkSlug := true
	if thread.Slug == "" {
		thread.Slug = uuid.New().String()
		checkSlug = false
	}
	threadNew, err := f.forumRepository.CreateThread(slug, thread)
	if err.Code == 409 {
		return f.forumRepository.GetThreadBySlug(thread.Slug)
	}

	if !checkSlug {
		threadNew.Slug = ""
	}
	return threadNew, err
}

func (f forumUsecase) GetUsers(slug string, params models.ParseParams) ([]models.User, models.Error) {
	return f.forumRepository.GetUsers(slug, params)
}

func (f forumUsecase) GetThreads(slug string, params models.ParseParams) ([]models.Thread, models.Error) {
	threads, err := f.forumRepository.GetThreads(slug, params)
	if err.Code == 404 {
		_, errr := f.forumRepository.GetForum(slug)
		if errr.Code == 200 {
			return threads, errr
		}
	}

	return threads, err
}


