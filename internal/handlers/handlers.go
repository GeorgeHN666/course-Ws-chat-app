package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/CloudyKit/jet/v6"
	"github.com/GeorgeHN666/course-ws-chat-app/internal/models"
	"github.com/gorilla/websocket"
)

/* Instructions

1.- func home is the connection to the path
2.- we upgrade the http to websocket
var Clients = make(map[models.WsConn]string)
3 .- after it connects we initialize listening with the websocket in the function ListenFormWs()

*/

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./html"),
	jet.InDevelopmentMode(),
)

var WsChan = make(chan models.WsPayload)

// To Activate the browser for websocket we need to use a variable that upgrade the browser
var upgradeConnectio = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var Clients = make(map[models.WsConn]string)

func Home(w http.ResponseWriter, r *http.Request) {

	err := renderPage(w, "home.jet", nil)
	if err != nil {
		log.Println(err)
	}

	log.Println(r)

}

// This is the WS endpoint
func WsEndpoint(w http.ResponseWriter, r *http.Request) {
	// Here the user connects with our Page
	// Here we upgrade the connection
	ws, err := upgradeConnectio.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("client connected to endpoint")

	// Here we define a model
	var response models.WsJSONresponse

	// Here we set the message
	response.Message = `<em><small>Connected to Server</small></em>`

	con := models.WsConn{Conn: ws}

	// Here we add a user To the map
	Clients[con] = ""

	// Here we start to listening, and every time the client send a json message we send it to a channel
	err = ws.WriteJSON(response)
	if err != nil {
		log.Println("error::", err)
	}

	go ListenForWs(&con)

}

func ListenForWs(conn *models.WsConn) {
	// This function its gonna recover us if something fails
	defer func() {
		if r := recover(); r != nil {
			log.Println("error:::", fmt.Sprintf("%v", r))
		}
	}()

	var payload models.WsPayload

	// This loop forever
	// Listen for a payload
	// And send to a channel
	// Here we start to listening, and every time the client send a json message we send it to a channel
	for {
		err := conn.ReadJSON(&payload)
		if err != nil {

		} else {
			payload.Conn = *conn
			WsChan <- payload
		}
	}

}

func ListenToWsChannel() {
	var response models.WsJSONresponse

	// We create a for loop
	for {

		// Here we take the data from the channel and store it in a variable
		e := <-WsChan

		// Here we populate some fields form the response JSON
		// response.Actions = "Got Here"
		// response.Message = fmt.Sprintf("Some message, action was %s", e.Action)

		switch e.Action {
		case "username":
			Clients[e.Conn] = e.UserN
			users := getUserLists()
			response.Actions = "list_users"
			response.UsersConnected = users
			BroadcastToAll(response)
		case "left":
			response.Actions = "list_users"
			delete(Clients, e.Conn)
			users := getUserLists()
			// Here we set it in to a model
			response.UsersConnected = users
			// This is How We send A response
			BroadcastToAll(response)

		case "broadcast":
			response.Actions = "broadcast"
			response.Message = fmt.Sprintf("<strong>%s</strong>: %s", e.UserN, e.Message)
			BroadcastToAll(response)

		}

		// Here we broadcast the message to all the users that are in the chat

	}

}

func getUserLists() []string {

	var UserList []string

	for _, v := range Clients {
		if v != "" {
			UserList = append(UserList, v)
		}
	}

	sort.Strings(UserList)

	return UserList

}

func BroadcastToAll(response models.WsJSONresponse) {
	// Here it says, For every client that we know about send this response
	for client := range Clients {
		err := client.WriteJSON(response)
		// If we get an error it means that the client has leave and we close the connection and remove them from the map
		if err != nil {
			log.Println("websocket err:::")
			_ = client.Close()
			delete(Clients, client)
		}
	}
}

// Here we display the page
func renderPage(w http.ResponseWriter, tmpl string, data jet.VarMap) error {

	view, err := views.GetTemplate(tmpl)
	if err != nil {
		log.Panicln(err)
		return err
	}

	err = view.Execute(w, data, nil)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil

}
