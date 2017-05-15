package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type appError struct {
	Error   error
	Message string
	Code    int
}

type player struct {
	id   uint64
	Name string
}

var nrOfRequests = 0
var playerMap = make(map[uint64]player)

func main() {
	http.HandleFunc("/player/", servePlayer)
	log.Print("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func servePlayer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Print("GET /player")
		getPlayer(w, r)
	case "POST":
		log.Print("POST /player")
		postPlayer(w, r)
	default:
		http.Error(w, "Only GET and POST are supported", http.StatusMethodNotAllowed)
	}
	nrOfRequests++
}

func getPlayer(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.RequestURI()
	idString := strings.Replace(uri, "/player/", "", 1)
	var err *appError
	if len(idString) > 0 {
		log.Print("getting player by id ", idString)
		var p player
		p, err = getPlayerByID(idString)
		if err == nil {
			fmt.Fprint(w, p)
		}
	} else {
		log.Print("getting all players ")
		var ps []player
		ps, err = getAllPlayers()
		if err == nil {
			fmt.Fprint(w, ps)
		}
	}
	if err != nil {
		http.Error(w, err.Message, err.Code)
	}

}

func getPlayerByID(idString string) (player, *appError) {
	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		errorString := fmt.Sprintf("Not valid format of player ID %v: %v", idString, err)
		return player{}, &appError{err, errorString, http.StatusBadRequest}
	}
	p, ok := playerMap[id]
	if !ok {
		errorString := fmt.Sprintf("No player exists with id: %v", id)
		return player{}, &appError{errors.New(errorString), errorString, http.StatusNotFound}
	}
	log.Printf("Player found: %v", p)
	return p, nil
}

func getAllPlayers() ([]player, *appError) {
	allPlayers := make([]player, 0, len(playerMap))
	for _, player := range playerMap {
		allPlayers = append(allPlayers, player)
	}
	log.Print("all players are ", allPlayers)
	return allPlayers, nil
}

func postPlayer(w http.ResponseWriter, r *http.Request) {
	pID, err := createPlayer(r)
	if err != nil {
		http.Error(w, err.Message, err.Code)
		return
	}
	fmt.Fprint(w, pID)
}

func createPlayer(r *http.Request) (uint64, *appError) {
	var newPlayer player
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return 0, &appError{err, err.Error(), http.StatusInternalServerError}
	}
	err = json.Unmarshal(jsonData, &newPlayer)
	if err != nil {
		return 0, &appError{err, err.Error(), http.StatusBadRequest}
	}
	newPlayerID := uint64(len(playerMap))
	newPlayer.id = newPlayerID
	playerMap[newPlayerID] = newPlayer
	return newPlayerID, nil
}
