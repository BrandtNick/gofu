package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

const (
	maxSize     = 2000000 // Maximum size in bytes (2MB)
	storagePath = "./"    // Where to save the files
)

// Allowed mime types
var allowedTypes = []string{"image/jpeg", "image/gif", "image/png"}

func assertTypeIsAllowed(mimeType string) bool {
	for _, v := range allowedTypes {
		if mimeType == v {
			return true
		}
	}
	return false
}

func assertIsAllowedSize(size int64) bool {
	if size < maxSize {
		return true
	}
	return false
}

func upload(w http.ResponseWriter, r *http.Request) {
	uploadedFile, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error ", err)
		return
	}
	defer uploadedFile.Close()

	mimeType := handler.Header["Content-Type"][0]
	allowedType := assertTypeIsAllowed(mimeType)
	if !allowedType {
		w.Write([]byte("File type not supported"))
		return
	}

	allowedSize := assertIsAllowedSize(handler.Size)
	if !allowedSize {
		w.Write([]byte("File size too big"))
		return
	}

	file, err := os.Create(storagePath + handler.Filename)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(uploadedFile)
	if err != nil {
		fmt.Println(err)
	}

	file.Write(fileBytes)

	w.Write([]byte("Success, file uploaded!"))
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/upload", upload).Methods("POST")

	port := "8080"
	fmt.Println("listening on port: ", port)
	err := http.ListenAndServe(":"+port, r)
	if err != nil {
		fmt.Print(err)
	}
}
