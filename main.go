package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const maxUploadSize = 32 << 20 // 32 mb
const uploadPath = "./tmp"

func main() {
	http.HandleFunc("/upload", uploadFileHandler())

	fs := http.FileServer(http.Dir(uploadPath))
	http.Handle("/files/", http.StripPrefix("/files", fs))

	log.Print("Server started on localhost:8080, use /upload for uploading files and /files/{fileName} for downloading")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			t, _ := template.ParseFiles("upload.gtpl")
			t.Execute(w, nil)
			return
		}
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			fmt.Printf("Could not parse multipart form: %v\n", err)
			renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
			return
		}

		fhs := r.MultipartForm.File["uploadFiles"]
		index := 1
		if len(fhs) > 0 {
			fmt.Println("Uploading files: ")
			w.Write([]byte("Upload files:\n"))

		}
		for _, fileHeader := range fhs {
			// parse and validate file and post parameters
			file, err := fileHeader.Open()
			if err != nil {
				renderError(w, "INVALID_FILE", http.StatusBadRequest)
				return
			}
			defer file.Close()

			fileSize := fileHeader.Size
			fileName := fileHeader.Filename

			fileBytes, err := io.ReadAll(file)
			if err != nil {
				renderError(w, "INVALID_FILE", http.StatusBadRequest)
				return
			}

			newPath := filepath.Join(uploadPath, fileName)

			// write file
			newFile, err := os.Create(newPath)
			if err != nil {
				renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
				return
			}
			defer newFile.Close() // idempotent, okay to call twice
			if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
				renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
				return
			}

			fmt.Printf("%2d. %s (%d bytes)\n", index, fileName, fileSize)
			w.Write([]byte(fmt.Sprintf("%2d. %s (%d bytes) [SUCCESS]\n", index, fileName, fileSize)))

			index++
		}

	})
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

/*func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}*/
