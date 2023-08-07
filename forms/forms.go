package forms

import "clipboard/models"

type AddClipBoardRequest struct {
	Content string `json:"content"`
}

type AddClipBoardResponse struct {
	Code int `json:"code"`
}

type GetClipBoardListRequest struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type GetClipBoardListResponse struct {
	Code int                      `json:"code"`
	Data []models.ClipBoardEntity `json:"data"`
}

type GetLatestClipBoardResponse struct {
	Code int                    `json:"code"`
	Data models.ClipBoardEntity `json:"data"`
}
