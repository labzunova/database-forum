package usecase

import (
	"DBproject/internal/threads"
	"DBproject/models"
	"fmt"
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
	id, _ = t.threadsRepository.GetThreadIDBySlug(slugOrId)
	return t.threadsRepository.GetThreadPostsById(id, params)
}

func (t threadsUsecase) VoteThread(slugOrId string, vote models.Vote) (models.Thread, models.Error) {
	id := t.SlugOrID(slugOrId)
	fmt.Println("slug or id ",slugOrId, id)
	if id != 0 {
		fmt.Println("vote by id")
		err := t.threadsRepository.VoteThreadById(id, vote)
		if err.Code != 200 {
			return models.Thread{}, err
		}
		return t.threadsRepository.GetThread("", id)
	}

	fmt.Println("vote by slug")
	err := t.threadsRepository.VoteThreadBySlug(slugOrId, vote)
	if err.Code != 200 {
		return models.Thread{}, err
	}

	return t.threadsRepository.GetThread(slugOrId, 0)
}

// extra

func (t threadsUsecase) SlugOrID(slugOrId string) int {
	id, errID := strconv.Atoi(slugOrId)
	if errID == nil {
		return id
	}

	return 0
}