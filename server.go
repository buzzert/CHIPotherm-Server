package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type commandControlServer struct {
	commands    chan string
	stateChange chan bool
	state       ThermostatState
}

// The CHIP responds to these messages
// Will return either "refresh", in which the server expects a call to updateState, or
// the string representation of the new state (i.e., "enabled 80")
func (s *commandControlServer) poll(w http.ResponseWriter, r *http.Request) {
	log.Print("Waiting for command...")
	context := r.Context()
	select {
	case command := <-s.commands:
		fmt.Fprintf(w, "%s", command)
		break
	case <-context.Done():
		log.Print("Connection closed unexpectedly")
	}
}

// Returns the cached state on the server
func (s *commandControlServer) getCachedState(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Encode as JSON
	b, err := json.Marshal(s.state)
	if err == nil {
		w.Write(b)
	} else {
		log.Printf("Error encoding as JSON: %s", err)
	}
}

// Just sends the "state" command to the CHIP to ask it to send us a new state via
// updateState()
func (s *commandControlServer) refreshState(w http.ResponseWriter, r *http.Request) {
	s.commands <- "refresh"
	<-s.stateChange
	s.getCachedState(w, r)
}

// Sets the server's cached state and notifies the CHIP
func (s *commandControlServer) setState(w http.ResponseWriter, r *http.Request) {
	requestBytes, err := ioutil.ReadAll(r.Body)
	requestBody := string(requestBytes)

	if err == nil {
		s.state = StateFromString(requestBody)
		log.Printf("Setting new state: %s", s.state.String())

		// Optionally message observers
		select {
		case s.commands <- s.state.String():
		default:
			// No message sent
			// TODO: Should this get queued instead?
		}
	} else {
		log.Printf("Error parsing state string: %s", requestBody)
	}
}

// From CHIP: Update the server's cached state
func (s *commandControlServer) updateState(w http.ResponseWriter, r *http.Request) {
	responseBytes, err := ioutil.ReadAll(r.Body)

	if err == nil {
		// Response is JSON
		err := json.Unmarshal(responseBytes, &s.state)
		if err != nil {
			log.Printf("Error parsing JSON: %s", err)
		} else {
			log.Printf("Got new state: %s", s.state.String())
		}
	} else {
		log.Printf("Error parsing state string: %s", string(responseBytes))
	}

	// Notify state change
	select {
	case s.stateChange <- true:
	default:
		break
	}
}

func main() {
	commandChannel := make(chan string)
	stateChangeChannel := make(chan bool)

	server := commandControlServer{
		commands:    commandChannel,
		stateChange: stateChangeChannel,
	}

	router := mux.NewRouter()
	router.HandleFunc("/poll", server.poll)
	router.HandleFunc("/updateState", server.updateState)
	router.HandleFunc("/refreshState", server.refreshState)

	router.HandleFunc("/getCachedState", server.getCachedState)
	router.HandleFunc("/setState", server.setState)

	listenAddress := ":43001"
	log.Print("Listening on " + listenAddress)

	http.Handle("/", http.TimeoutHandler(router, time.Duration(1*time.Hour), "replaceme"))
	http.ListenAndServe(":43001", nil)
}
