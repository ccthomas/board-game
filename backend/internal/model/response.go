package model

type healthResponse struct {
	Service   string      `json:"service"`
	Database  string      `json:"database"`
	Timestamp interface{} `json:"timestamp"`
}
