package receive

import "backend/internal/app/model"

type PackReceive struct {
	Code int          `json:"code"`
	Data _PackReceive `json:"data"`
	Msg  string       `json:"msg"`
}

type _PackReceive struct {
	Tutorial  model.Tutorial `json:"tutorial"`
	PackLog   model.PackLog  `json:"pack_log"`
	StartPage string         `json:"start_page"`
	FileName  string         `json:"file_name"`
	Message   string         `json:"message"`
}
