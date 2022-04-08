package models

type WsJSONresponse struct {
	Actions        string   `json:"action"`
	Message        string   `json:"message"`
	MessType       string   `json:"message_type"`
	UsersConnected []string `json:"usersConnected"`
}
