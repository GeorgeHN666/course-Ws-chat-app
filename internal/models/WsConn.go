package models

import "github.com/gorilla/websocket"

type WsConn struct {
	*websocket.Conn
}
