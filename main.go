package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type wrappable func(http.ResponseWriter, *http.Request)

type MiddleWare struct {
	handler http.Handler
}

type album struct {
	Id     int     `json:"id"`
	Artist string  `json:"artist"`
	Album  string  `json:"album"`
	Price  float32 `json:"price"`
}

func printEncodedAlbums() {
	for _, a := range albums {
		j, err := json.Marshal(a)
		if err != nil {
			log.Println("Wasn't able to encode a for some reason")
		}
		log.Printf("%x \n", j)
	}
}

func getConfigString() string {
	// fmt.Sprintf("docker:docker@tcp(127.0.0.1:3308)/docker",
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)
}

func insertIntoTable() {
	db, err := sql.Open("mysql", getConfigString())

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
		fmt.Printf("Inserting album %d right now\n", id)
		if _, err := stmt.Exec(album.Artist, album.Album, album.Price); err != nil {
			log.Fatal(err)
		}
	}
}

var albums = []album{
	{
		Id:     1,
		Artist: "Rick James",
		Album:  "Rick Album 1",
		Price:  11.06,
	},
	{
		Id:     2,
		Artist: "James Rick",
		Album:  "James Album 2",
		Price:  12.57,
	},
	{
		Id:     3,
		Artist: "Person Next",
		Album:  "PND1",
		Price:  13.87,
	},
}

// querying the base url and base line understanding
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

func addAlbumToDb(w http.ResponseWriter, req *http.Request) {
	var album album
	// albums = append(albums, album)
	db, err := sql.Open("mysql", getConfigString())

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
	decoder := json.NewDecoder(req.Body)
	err = decoder.Decode(&album)

	if err != nil {
		log.Fatal(err)
	}
	stmt.Exec(album.Artist, album.Album, album.Price)
	io.WriteString(w, "Upload Successful\n")

}

func addAlbum(w http.ResponseWriter, req *http.Request) {
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
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(albums)
}

func getAlbum(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	url := req.URL
	str_value := url.Query().Get("id")

	if str_value == "" {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Invalid Query")
		io.WriteString(w, "try something like URL:?id=<id_num>")
		return
	}

	value, err := strconv.Atoi(str_value)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Invalid Query")
		panic(err)
	}

	for _, album := range albums {
		if album.Id == value {
			json.NewEncoder(w).Encode(album)
			return
		}
	}
	w.WriteHeader(http.StatusInternalServerError)
	io.WriteString(w, "Invalid Query")
	json.NewEncoder(w).Encode(fmt.Sprintf("Unable to find Album id %d", value))
}

func ApiLogger(fn wrappable) wrappable {
	// same as the python wrapper we are creating a function that accepts the
	//args and then we are returning that function with additional steps wrapped around it
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Recieved request at: ", time.Now())
		fmt.Println("Endpoint being serviced: ", req.URL)
		fn(w, req)
		fmt.Println("Done handling request at: ", time.Now())
	}
}

func NewLoggerMiddleware(handle http.Handler) *MiddleWare {
	return &MiddleWare{handle}
}

func (m *MiddleWare) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("Recieved request at: ", time.Now())
	fmt.Println("Endpoint being serviced: ", req.URL)
	m.handler.ServeHTTP(w, req)
	fmt.Println("Done handling request at: ", time.Now())

}

func main() {
	// creating wrapped functions
	// getAlbumsWrapped := NewLoggerMiddleware(http.HandlerFunc(getAlbums))
	getAlbumsWrapped := ApiLogger(getAlbums)
	getAlbumWrapped := ApiLogger(getAlbum)
	addAlbumDBWrapper := ApiLogger(addAlbumToDb)

	http.HandleFunc("GET /getAlbums", getAlbumsWrapped)
	http.HandleFunc("POST /addAlbum", addAlbumDBWrapper)
	http.HandleFunc("GET /getAlbum", getAlbumWrapped)
	fmt.Println("We have registered all handles")
	log.Fatal(
		http.ListenAndServe(":8080", nil))
}
