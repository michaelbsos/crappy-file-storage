package main

import (
	_ "embed"
	"io"
	"log"
	"net/http"
	"os"
)

//go:embed index.html
var indexTemplate []byte

func main() {
	http.HandleFunc("/", indexHandler)

	log.Println("Listening on 0.0.0.0:6666")
	http.HandleFunc("/upload", uploadHandler)
	http.ListenAndServe("0.0.0.0:6666", nil)
}

// indexHandler handles requests to / and shows the contents of index.html
func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write(indexTemplate)
}

// uploadHandler fetches a file from the file field on a multipart/form-data form and save it
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(200000) // grab the multipart form
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Oops. Something went wrong"))
		log.Println(err)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Oops. Something went wrong"))
		log.Println(err)
		return
	}
	defer file.Close()

	resFile, err := os.Create("./uploads/" + handler.Filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Oops. Something went wrong"))
		log.Println(err)
		return
	}
	defer resFile.Close()

	io.Copy(resFile, file)
	defer resFile.Close()

	w.Write([]byte("Thanks!"))
}
