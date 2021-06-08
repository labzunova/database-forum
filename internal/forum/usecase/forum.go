package usecase

import (
	"DBproject/internal/forum"
	"DBproject/models"
	"fmt"
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
	threadNew, err := f.forumRepository.CreateThread(slug, thread)
	if err.Code == 409 {
		return f.forumRepository.GetThreadBySlug(thread.Slug)
	}

	return threadNew, err
}

func (f forumUsecase) GetUsers(slug string, params models.ParseParams) ([]models.User, models.Error) {
	users, err := f.forumRepository.GetUsers(slug, params)
	fmt.Println("       slug", slug)
	fmt.Println("       users", users)
	if len(users) == 0 {
		fmt.Println("no users was found")
		_, errr := f.forumRepository.GetForum(slug)
		fmt.Println(errr)
		if errr.Code == 404 {
			return users, errr
		}
	}

	return users, err
}

func (f forumUsecase) GetThreads(slug string, params models.ParseParams) ([]models.Thread, models.Error) {
	forumm, errr := f.forumRepository.GetForum(slug)
	if errr.Code != 200 {
		return []models.Thread{}, errr
	}

	threads, _ := f.forumRepository.GetThreads(forumm.Slug, params)
	fmt.Println("get threads done", errr)


	return threads, errr
}
