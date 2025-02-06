package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/websocket"

	"github.com/bensivo/streaming-analytics-example/webapp/service"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { // Allow all origins for demo purposes
		return true
	},
}

type HttpController struct {
	svc *service.Service
}

func NewHttpController(svc *service.Service) *HttpController {
	return &HttpController{
		svc: svc,
	}
}

func (c *HttpController) Start() {
	// Serve static files from the "static" directory
	fs := http.FileServer(http.Dir("/app/static"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", c.handleWebSocket)

	port := ":8080"
	log.Printf("Server starting on http://localhost%s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}

func (c *HttpController) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("New WebSocket connection from %s", conn.RemoteAddr().String())

	for {
		// Read message from browser
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		log.Printf("Received message: %s", msg)

		var rating service.Rating
		if err := json.Unmarshal([]byte(msg), &rating); err != nil {
			log.Printf("Error parsing JSON: %v", err)
			return
		}

		// Process the rating
		c.svc.SendRating(rating)

		// Echo the message back to the client
		// log.Printf("Received message: %s", string(msg))
		// if err := conn.WriteMessage(msgType, msg); err != nil {
		// 	log.Printf("Error writing message: %v", err)
		// 	break
		// }
	}
}