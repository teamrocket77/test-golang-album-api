package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"testing"
)

var albums []album

type album struct {
	Id     int    `json:"id"`
	Artist string `json:"artist"`
	Album  string `json:"album"`
}

func printAlbums() {
	// prints the albums in json format
	for _, a := range albums {
		j, err := json.Marshal(a)
		if err != nil {
			log.Println("Wasn't able to encode a for some reason")
		}
		log.Printf("%x \n", j)
	}
}

func defaultAlbums() {
	// creates default album list
	albums = append(
		albums,
		album{
			Id:     1,
			Artist: "Rick James",
			Album:  "Rick Album 1",
		})
	albums = append(
		albums,
		album{
			Id:     2,
			Artist: "James Rick",
			Album:  "James Album 2",
		})
	albums = append(
		albums,
		album{
			Id:     3,
			Artist: "Person Next",
			Album:  "PND1",
		})
}

func rootHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Recieved Request at / ")
	switch httpMethod := req.Method; httpMethod {
	case "PUT":
		io.WriteString(w, "Method was a PUT")
	case "DELETE":
		io.WriteString(w, "Method was a DELETE")
	case "POST":
		io.WriteString(w, "Method was a POST")
	default:
		io.WriteString(w, "Method was a Get")
	}
}

func addAlbum(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("500 - Unsupported method for this endpoint")
		return
	}
	var album album
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&album)
	if err != nil {
		panic(err)
	}
	albums = append(albums, album)
	io.WriteString(w, "Upload Successful\n")
}

func getAlbums(w http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("500 - Unsupported method for this endpoint")
		return
	}
	w.Header().Set("Content-Type", "application/json")

	io.WriteString(w, req.URL.RequestURI())

	json.NewEncoder(w).Encode(albums)
}

func getAlbum(w http.ResponseWriter, req *http.Request) {
	log.Printf("url: %s", req.URL)
	if req.Method != "GET" {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode("500 - Unsupported method for this endpoint")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	url := req.URL
	str_value := url.Query().Get("id")
	value, err := strconv.Atoi(str_value)

	if err != nil {
		panic(err)
	}
	for _, album := range albums {
		if album.Id == value {
			json.NewEncoder(w).Encode(album)
			return
		}
	}
	json.NewEncoder(w).Encode(fmt.Sprintf("Unable to find Album id %d", value))
}

func main() {
	defaultAlbums()
	http.HandleFunc("/getAlbums", getAlbums)
	http.HandleFunc("/addAlbum", addAlbum)
	http.HandleFunc("/getAlbum", getAlbum)
	log.Fatal(
		http.ListenAndServe(":8080", nil))
}

func testgetAlbum(t *testing.T) {
}
