package models

type WsPayload struct {
	UserN   string `json:"username"`
	Action  string `json:"action"`
	Message string `json:"message"`
	Conn    WsConn `json:"-"`
}
