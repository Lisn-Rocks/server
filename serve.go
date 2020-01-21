package main

import (
	"database/sql"
    "os"
    "net/http"
	"log"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/sharpvik/lisn-backend/config"
)



var logr *log.Logger

var mux *http.ServeMux

var db *sql.DB



type mainHandler struct {}

// ServeHTTP function is the entry point for server's routing mechanisms.
// It uses mux to delegate request to a proper handler function.
func (*mainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logr.Printf( "URL: %s", r.URL.String() )
	mux.ServeHTTP(w, r)
}



func serveRandomSong(w http.ResponseWriter, r *http.Request) {
	rows, _ := db.Query("SELECT id FROM songs")

	var id int

	rows.Next()
	rows.Scan(&id)

	path := fmt.Sprintf("store/%d.mp3", id)
	http.ServeFile(w, r, path)
}

func serveSongByID(w http.ResponseWriter, r *http.Request) {
	id := "1" // in production id must be read from request body (type string)
	path := fmt.Sprintf("store/%s.mp3", id)
	http.ServeFile(w, r, path)
}



func initDB(db *sql.DB) {
	stmt, _ := db.Prepare("CREATE TABLE IF NOT EXISTS songs (id INTEGER PRIMARY KEY, title TEXT, duration FLOAT, genre TEXT, author TEXT, album TEXT NULL)")
	stmt.Exec()
}

func insertSongs(db *sql.DB) {
	stmt, _ := db.Prepare("INSERT INTO songs (title, duration, genre, author, album) VALUES (?, ?, ?, ?, ?)")
	stmt.Exec("Some Title", 3.483, "Jazz", "Viky is Kinky", "Void Functions")
}



func main() {
	db, _ = sql.Open("sqlite3", "./songs.db")

	if config.InitRequired {
		initDB(db)
		insertSongs(db)
	}


	logr = log.New(os.Stdout, "", log.Ltime)


	server := http.Server{
		Addr:		config.Port,
		Handler:	&mainHandler{},
		ErrorLog:	logr,
	}


	mux = http.NewServeMux()
	mux.HandleFunc("/", serveRandomSong)
	mux.HandleFunc("/id", serveSongByID)


	logr.Printf("Serving at localhost:%s", config.Port)
	server.ListenAndServe()
}
