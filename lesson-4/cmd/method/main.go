package main

import (
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Handler struct {
	filesList *FilesList
}

type UploadHandler struct {
	HostAddr  string
	UploadDir string
	filesList *FilesList
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		ext := r.FormValue("ext")
		err := h.filesList.refresh(ext) // TODO Cделано так в учебных целях
		if err != nil {
			http.Error(w, "Unable to get file list", http.StatusInternalServerError)
			return
		}
		enc := json.NewEncoder(w)
		enc.Encode(h.filesList.list)

	case http.MethodPost:
		var employee Employee

		contentType := r.Header.Get("Content-Type")

		switch contentType {
		case "application/json":
			err := json.NewDecoder(r.Body).Decode(&employee)
			if err != nil {
				http.Error(w, "Unable to unmarshal JSON", http.StatusBadRequest)
				return
			}

		case "application/xml":
			err := xml.NewDecoder(r.Body).Decode(&employee)
			if err != nil {
				http.Error(w, "Unable to unmarshal XML", http.StatusBadRequest)
				return
			}

		default:
			http.Error(w, "Unknown content type", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(w, "Got a new employee!\nName: %s\nAge: %dy.o.\nSalary %0.2f\n",
			employee.Name,
			employee.Age,
			employee.Salary,
		)

	}
}

func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusBadRequest)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		http.Error(w, "Unable to read file", http.StatusBadRequest)
		return
	}

	// ДЗ #3
	id := uuid.New()
	var b []byte
	for i := 0; i < 16; i++ {
		b[i] = id[i]
	}
	sid := hex.EncodeToString(b)
	filePath := h.UploadDir + "/" + sid + "-" + header.Filename

	err = ioutil.WriteFile(filePath, data, 0777)
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	fileLink := h.HostAddr + "/" + header.Filename
	fmt.Fprintln(w, fileLink)
}

func main() {
	handler := &Handler{}
	http.Handle("/", handler)

	uploadHandler := &UploadHandler{
		UploadDir: "upload",
	}
	http.Handle("/upload", uploadHandler)

	handler.filesList = &FilesList{
		list: nil,
		dir:  "upload",
	}
	uploadHandler.filesList = handler.filesList
	handler.filesList.refresh("")

	srv := &http.Server{
		Addr:         ":80",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Panic(err)
		}
	}()

	dirToServe := http.Dir(uploadHandler.UploadDir)

	fs := &http.Server{
		Addr:         ":8080",
		Handler:      http.FileServer(dirToServe),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	fs.ListenAndServe()
}
