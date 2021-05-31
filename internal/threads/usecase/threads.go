package usecase

import (
	"DBproject/internal/threads"
	"DBproject/models"
	"strconv"
)

type threadsUsecase struct {
	threadsRepository threads.ThreadsRepo
}

func NewThreadsUsecase(repo threads.ThreadsRepo) threads.ThreadsUsecase {
	return &threadsUsecase{
		threadsRepository: repo,
	}
}

func (t threadsUsecase) GetThread(slug string) (models.Thread, models.Error) {
	id := t.SlugOrID(slug)
	return t.threadsRepository.GetThread(slug, id)
}

func (t threadsUsecase) UpdateThread(slug string, thread models.Thread) (models.Thread, models.Error) {
	thread.ID = t.SlugOrID(slug)
	return t.threadsRepository.UpdateThread(slug, thread)
}

func (t threadsUsecase) GetThreadPosts(slug string, params models.ParseParamsThread) ([]models.Post, models.Error) {
	return t.threadsRepository.GetThreadPosts(slug, params)
}

func (t threadsUsecase) VoteThread(slug string, vote models.Vote) (models.Thread, models.Error) {
	return t.threadsRepository.VoteThread(slug, vote)
}

func (t threadsUsecase) SlugOrID(slug string) int {
	id, errID := strconv.Atoi(slug)
	if errID != nil {
		return id
	}

	return 0
}