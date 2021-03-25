package usecase

import (
	"DBproject/internal/threads"
	"DBproject/internal/threads/delivery/http"
	"DBproject/models"
)

type threadsUsecase struct {
	threadsRepository threads.ThreadsRepo
}

func NewThreadsUsecase(repo threads.ThreadsRepo) threads.ThreadsUsecase {
	return &threadsUsecase{
		threadsRepository: repo,
	}
}

func (t threadsUsecase) CreateThread(slug string, thread models.Thread) (models.Thread, models.Error) {
	return t.threadsRepository.CreateThread(slug, thread)
}

func (t threadsUsecase) GetThread(slug string) (models.Thread, models.Error) {
	return t.threadsRepository.GetThread(slug)
}

func (t threadsUsecase) UpdateThread(slug string, thread models.Thread) (models.Thread, models.Error) {
	return t.threadsRepository.UpdateThread(slug, thread)
}

func (t threadsUsecase) GetThreadPosts(slug string, params http.ThreadsParse) ([]models.Post, models.Error) {
	return t.threadsRepository.GetThreadPosts(slug, params)
}

func (t threadsUsecase) VoteThread(slug string, vote models.Vote) (models.Thread, models.Error) {
	return t.threadsRepository.VoteThread(slug, vote)
}
