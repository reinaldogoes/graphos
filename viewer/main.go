package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type msg struct {
	Num int
}

type clientStatus struct {
	Subscriptions map[string]bool
}

type Data struct {
	ID    string    `json:"id"`   // pwr.v
	Type  string    `json:"type"` // data
	Value dataValue `json:"value"`
}

type dataValue struct {
	Timestamp int64       `json:"timestamp"` // 1474891025120
	Value     interface{} `json:"value"`     // 30.00046680219344
}

type handleFn func()

var clientList map[*websocket.Conn]clientStatus
var systemStatus map[string]dataValue
var handleList map[string]handleFn

//var lock = sync.RWMutex{}

func init() {
	log.Println("Initializing Telemetry Hub")

	clientList = make(map[*websocket.Conn]clientStatus)
	systemStatus = make(map[string]dataValue)
	handleList = make(map[string]handleFn)
}

// MakeTimestamp creates an int64 with current timestamp in the format expected by OpenMCT
func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func makeData(id string, timestamp int64, value interface{}) (r Data) {
	r = Data{}

	r.ID = id
	r.Type = "data"
	r.Value.Timestamp = timestamp
	r.Value.Value = value

	return
}

// CloseAll close all websocket connections
func CloseAll() {
	for c := range clientList {
		c.Close()
	}
}

// HandleFunc set a subsystem handle function
func HandleFunc(subsystemIdentifier string, handleFunc handleFn) {
	handleList[subsystemIdentifier] = handleFunc
}

// SetDataValue set a value to subsystem
func SetDataValue(identifier string, timeStamp int64, value interface{}) {
	dv := dataValue{}

	if timeStamp == -1 {
		timeStamp = MakeTimestamp()
	}

	dv.Timestamp = timeStamp
	dv.Value = value

	//lock.Lock()
	systemStatus[identifier] = dv
	//lock.Unlock()
}

// ListenAndServe websocket
func ListenAndServe(port int, timerInterval int) {

	log.Println("websocket port", port)

	go func() {
		for {
			time.Sleep(time.Millisecond * time.Duration(timerInterval))

			for subsystem := range handleList {
				if f, ok := handleList[subsystem]; ok && f != nil {
					f()
				}
			}

			for id, status := range systemStatus {
				SendValue(id, status.Timestamp, status.Value)
			}

		}
	}()

	http.HandleFunc("/", wsHandler)
	panic(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

// SendValue send value to subsystem by identifier
func SendValue(id string, timestamp int64, value interface{}) {
	for conn, c := range clientList {
		//lock.Lock()
		//defer lock.Unlock()

		if v, ok := c.Subscriptions[id]; ok && v {

			var p = makeData(id, timestamp, value)

			if err := conn.WriteJSON(p); err != nil {
				log.Println(err)
			}
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	//if r.Header.Get("Origin") != "http://"+r.Host {
	//	http.Error(w, "Origin not allowed", 403)
	//	return
	//}
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
	}

	log.Println("Client connection")
	clientList[conn] = clientStatus{Subscriptions: make(map[string]bool)}

	log.Println("Client tot:", len(clientList))

	go telemetryWs(conn)
}

func telemetryWs(conn *websocket.Conn) {

	defer func() {
		conn.Close()
		delete(clientList, conn)
	}()

	for {

		messageType, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Println("Socket close normal")
				return
			}
			log.Printf("error: %v", err)
			return
		}

		msg := string(p)
		log.Printf("msg:[%s] type:%d\r\n", msg, messageType)
		msgArray := strings.Split(msg, " ")
		for k, v := range msgArray {
			log.Printf("%d\t[%s]\r\n", k, v)
		}

		switch msgArray[0] {
		case "e":
			for k, _ := range clientList[conn].Subscriptions {
				clientList[conn].Subscriptions[k] = false
			}
			err = conn.Close()
			if err != nil {
				log.Println(err.Error())
			}
		case "subscribe":
			if len(msgArray) < 2 {
				log.Println("error: no subscribe parameter")
			}
			clientList[conn].Subscriptions[msgArray[1]] = true
			log.Println("client subscribe", msgArray[1])
		case "unsubscribe":
			if len(msgArray) < 2 {
				log.Println("error: no unsubscribe parameter")
			}
			clientList[conn].Subscriptions[msgArray[1]] = false
		case "history":

		case "list":
			var s []string
			for k, v := range clientList[conn].Subscriptions {
				if v {
					s = append(s, k)
				}
				if err = conn.WriteJSON(s); err != nil {
					log.Println(err)
				}
			}
		}

		//err = conn.WriteMessage(messageType, p)
		//if err != nil {
		//	return
		//}

		//if err = conn.WriteJSON(m); err != nil {
		//	fmt.Println(err)
		//}
	}
}

func main() {

	go func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, os.Interrupt)
		<-sc
		// close all connections
		CloseAll()

		fmt.Print("\n")
		log.Println("Have a nice day!")
		os.Exit(0)
	}()

	HandleFunc("graphos", graphosFunc)
	ListenAndServe(8081, 1000) // 20 = 50 fps

}

var x float64

func graphosFunc() {
	x += 0.01
	y := math.Sin(x) * (x / 2.0 * math.Pi)
	SetDataValue("graphos", MakeTimestamp(), y)
}
