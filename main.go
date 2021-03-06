package main

import (
	"archive/zip"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "7777"
	}

	http.ListenAndServe(":"+port, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		fileName := r.FormValue("filename")
		if fileName == "" {
			fileName = "download.zip"
		}
		w.Header().Add("Content-Disposition", "attachment; filename=\""+fileName+"\"")
		log.Println("Generating zip for " + fileName)

		archive := zip.NewWriter(w)

		for path, urls := range r.PostForm {
			log.Println("Downloading key: " + path + " URL: " + urls[0])
			entryHeader := &zip.FileHeader{
				Name: path,
			}
			entryHeader.SetModTime(time.Now())
			entryHeader.SetMode(os.ModePerm)
			entry, err := archive.CreateHeader(entryHeader)
			if err != nil {
				http.Error(w, "error", http.StatusBadRequest)
				log.Println(err)
				return
			}
			download, err := http.Get(urls[0])
			if err != nil {
				http.Error(w, "error", http.StatusBadRequest)
				log.Println(err)
				return
			}
			io.Copy(entry, download.Body)
		}
		err := archive.Close()
		if err != nil {
			http.Error(w, "error", http.StatusBadRequest)
			log.Println(err)
			return
		}
	}))
}
