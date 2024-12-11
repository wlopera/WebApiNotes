package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type Note struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

var noteStore = make(map[string]Note)
var id int

// GetNotesHandler - GET - api/notes
func GetNotesHandler(w http.ResponseWriter, r *http.Request) {
	var notes []Note
	for _, v := range noteStore {
		notes = append(notes, v)
	}
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(notes)
	if err != nil {
		log.Printf("Error al convertir la nota a Json")
		http.Error(w, "Error al convertir las notas a JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// PostNotesHandler - POST - api/notes
func PostNotesHandler(w http.ResponseWriter, r *http.Request) {
	var note Note
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	note.CreatedAt = time.Now()
	id++
	k := strconv.Itoa(id)
	noteStore[k] = note

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(note)
	if err != nil {
		log.Printf("Error al convertir la nota a Json")
		http.Error(w, "Error al convertir la nota a JSON", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

// PutNotesHandler - PUT - api/notes/{id}
func PutNotesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]

	var noteUpdate Note
	err := json.NewDecoder(r.Body).Decode(&noteUpdate)
	if err != nil {
		log.Printf("Error al decodificar la solicitud %s:", k)
		http.Error(w, "Error al decodificar la solicitud", http.StatusBadRequest)
		return
	}

	if _, ok := noteStore[k]; ok {
		noteUpdate.CreatedAt = time.Now()
		noteStore[k] = noteUpdate
		w.WriteHeader(http.StatusNoContent)
	} else {
		log.Printf("Identificador %s no encontrado", k)
		http.Error(w, "Nota no encontrada", http.StatusNotFound)
	}
}

// DeleteNotesHandler - DELETE - api/notes/{id}
func DeleteNotesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k := vars["id"]

	if _, ok := noteStore[k]; ok {
		delete(noteStore, k)
		w.WriteHeader(http.StatusNoContent)
	} else {
		log.Printf("Nota no encontrada: %s", k)
		http.Error(w, "Nota no encontrada", http.StatusNotFound)
	}
}

func main() {
	r := mux.NewRouter().StrictSlash(false)

	r.HandleFunc("/api/notes", GetNotesHandler).Methods("GET")
	r.HandleFunc("/api/notes", PostNotesHandler).Methods("POST")
	r.HandleFunc("/api/notes/{id}", PutNotesHandler).Methods("PUT")
	r.HandleFunc("/api/notes/{id}", DeleteNotesHandler).Methods("DELETE")

	server := &http.Server{
		Addr:           ":8080",
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Println("Servidor levantado en http://localhost:8080...")
	log.Fatal(server.ListenAndServe())
}
