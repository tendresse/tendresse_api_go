package models

type JsonResponse struct {
	Code  int    `json:"status_code,omitempty"`
	Err   string `json:"err,omitempty"`
	Msg   string `json:"msg,omitempty"`
	Token string `json:"token,omitempty"`
}
