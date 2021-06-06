package models

type ParseParams struct {
	Limit int    `json:"limit"`
	Since string `json:"since"`
	Desc  bool   `json:"desc"`
}

type ParseParamsThread struct {
	Limit int    `json:"limit"`
	Since int    `json:"since"`
	Sort  string `lson:"since"`
	Desc  bool   `json:"desc"`
}

type PostsToCreate struct {
	Posts []Post
}
