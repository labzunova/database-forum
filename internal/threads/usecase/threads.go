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
	var err models.Error
	fmt.Println("slugorid:", slugOrId)
	if id != 0 {
		fmt.Println("id:", id)
		id, err = t.threadsRepository.GetThreadIDBySlug(slugOrId, id)
		if err.Code != 200 {
			err.Message = fmt.Sprintf("Can't find thread by id: %d", id)
			return nil, err
		}
	} else {
		id, err = t.threadsRepository.GetThreadIDBySlug(slugOrId, 0)
		fmt.Println("get by slug")
		if err.Code != 200 {
			err.Message = "Can't find thread by slug: " + slugOrId
			return nil, err
		}
	}
	return t.threadsRepository.GetThreadPostsById(id, slugOrId, params)
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