package models

type ParseParams struct {
	Limit int32  `json:"limit"`
	Since string `json:"since"`
	Desc  bool   `json:"desc"`
}

type ParseParamsThread struct {
	Limit int32  `json:"limit"`
	Since int64  `json:"since"`
	Sort  string `lson:"since"`
	Desc  bool   `json:"desc"`
}
