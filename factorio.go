package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/websocket"
	"github.com/gtaylor/factorio-rcon"
)

var (
	rconServerAddress  = os.Getenv("RCON_SERVER_ADDRESS")
	rconServerPassword = os.Getenv("RCON_SERVER_PASSWORD")
	serverAddress      = os.Getenv("SERVER_ADDRESS")
)

var upgrader = websocket.Upgrader{} // use default options

type rconHandler struct {
	client *rcon.RCON
}

func (rh *rconHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	connection, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}
	defer connection.Close()
	for {
		messageType, message, err := connection.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		trimmedText := strings.TrimSpace(string(message))

		response, err := rh.client.Execute(trimmedText)
		if err != nil {
			log.Println("execute error:", err)
			break
		}

		err = connection.WriteMessage(messageType, []byte(response.Body))
		if err != nil {
			log.Println("write error:", err)
			break
		}
	}
}

func home(writer http.ResponseWriter, request *http.Request) {
	homeTemplate.Execute(writer, "ws://"+request.Host+"/socket")
}

func main() {
	log.SetFlags(0)

	rconClient, err := rcon.Dial(rconServerAddress)
	if err != nil {
		panic(err)
	}
	err = rconClient.Authenticate(rconServerPassword)
	if err != nil {
		panic(err)
	}

	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			err = rconClient.Close()
			if err != nil {
				panic(err)
			}
			os.Exit(0)
		}
	}()

	http.Handle("/socket", &rconHandler{client: rconClient})
	http.HandleFunc("/", home)

	err = http.ListenAndServe(serverAddress, nil)
	if err != nil {
		panic(err)
	}
}

var data, err = Asset("public/index.html")
var homeTemplate = template.Must(template.New("").Parse(string(data)))
