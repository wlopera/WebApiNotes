Proyecto de prueba WEb API Notes en GO

```
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
```

![image](https://github.com/user-attachments/assets/2fcc4eb4-5b3d-4db9-93c6-e24d0a509cf2)

![image](https://github.com/user-attachments/assets/d330c995-d504-4508-86b4-f6726cec405d)

![image](https://github.com/user-attachments/assets/2c53cd20-63a2-4478-9ed3-55f0a9ca7827)

** Actuzliar registro
![image](https://github.com/user-attachments/assets/a3313b86-6756-4715-8c4a-0a553325170a)

![image](https://github.com/user-attachments/assets/0a38fceb-8069-41a8-b421-85e53cd8855f)

** En caso de error
![image](https://github.com/user-attachments/assets/a39e1c1f-1ca3-4484-bd90-e73a1bfdc12b)

** Borrar registro
![image](https://github.com/user-attachments/assets/e2157be6-2204-4415-a03c-3eb490d961bb)

![image](https://github.com/user-attachments/assets/46680715-edad-4421-aa6c-d6eba24624c8)

** En caso de error
![image](https://github.com/user-attachments/assets/d08b5d0b-a477-496b-bf8b-a5a4d650dc42)

** Servidor
![image](https://github.com/user-attachments/assets/afa29a93-ff64-4c72-9ef6-845dfb2fa72f)







