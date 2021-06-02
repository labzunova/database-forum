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

func (t threadsUsecase) GetThread(slugOrId string) (models.Thread, models.Error) {
	id := t.SlugOrID(slugOrId)
	return t.threadsRepository.GetThread(slugOrId, id)
}

func (t threadsUsecase) UpdateThread(slugOrId string, thread models.Thread) (models.Thread, models.Error) {
	id := t.SlugOrID(slugOrId)
	if id != 0 {
		return t.threadsRepository.UpdateThreadById(id, thread)
	}
	return t.threadsRepository.UpdateThreadBySlug(slugOrId, thread)
}

func (t threadsUsecase) GetThreadPosts(slugOrId string, params models.ParseParamsThread) ([]models.Post, models.Error) {
	id := t.SlugOrID(slugOrId)
	if id != 0 {
		return t.threadsRepository.GetThreadPostsById(id, params)
	}
	return t.threadsRepository.GetThreadPostsBySlug(slugOrId, params)
}

func (t threadsUsecase) VoteThread(slugOrId string, vote models.Vote) (models.Thread, models.Error) {
	id := t.SlugOrID(slugOrId)
	if id != 0 {
		return t.threadsRepository.VoteThreadById(id, vote)
	}
	return t.threadsRepository.VoteThreadBySlug(slugOrId, vote)
}

// extra

func (t threadsUsecase) SlugOrID(slugOrId string) int {
	id, errID := strconv.Atoi(slugOrId)
	if errID != nil {
		return id
	}

	return 0
}