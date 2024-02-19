package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"testing"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var albums []album

type album struct {
	Id     int     `json:"id"`
	Artist string  `json:"artist"`
	Album  string  `json:"album"`
	Price  float32 `json:"price"`
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

func getConfigString() string {
	return fmt.Sprintf("mysql://%s:%s@%s:%s/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

func insertIntoTable() {
	db, err := sql.Open("mysql", "docker:docker@tcp(127.0.0.1:3308)/docker")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	stmt, err := db.Prepare("INSERT INTO Albums( Artist, Album, Price) VALUES (?, ?, ?)")

	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()
	for _, album := range albums {
		id := album.Id
		fmt.Printf("Inserting album %d right now", id)
		if _, err := stmt.Exec(album.Artist, album.Album, album.Price); err != nil {
			log.Fatal(err)
		}
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
			Price:  11.06,
		})
	albums = append(
		albums,
		album{
			Id:     2,
			Artist: "James Rick",
			Album:  "James Album 2",
			Price:  12.57,
		})
	albums = append(
		albums,
		album{
			Id:     3,
			Artist: "Person Next",
			Album:  "PND1",
			Price:  13.87,
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
	insertIntoTable()
	http.HandleFunc("/getAlbums", getAlbums)
	http.HandleFunc("/addAlbum", addAlbum)
	http.HandleFunc("/getAlbum", getAlbum)
	fmt.Println("We have registered all handles")
	log.Fatal(
		http.ListenAndServe(":8080", nil))
}

func testgetAlbum(t *testing.T) {
}
