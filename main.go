package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//Post struct
type Post struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedOn   time.Time `json:"createon"`
}

//Store de la colección Post
var postStore = make(map[string]Post)
var id int = 0

//PostCreate crea un post
func PostCreate(w http.ResponseWriter, r *http.Request) {
	var post Post

	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		log.Fatal(err)
	}

	post.CreatedOn = time.Now()
	id++
	p := strconv.Itoa(id) //convierte de int a string
	postStore[p] = post

	j, err := json.Marshal(post)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}

//PostGetAll obtiene todos los posts creados
func PostGetAll(w http.ResponseWriter, r *http.Request) {
	var posts []Post

	for _, v := range postStore {
		posts = append(posts, v)
	}

	j, err := json.Marshal(posts)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "aplication/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

//PostUpdate actualiza un Post
func PostUpdate(w http.ResponseWriter, r *http.Request) {
	var postToUp Post
	params := mux.Vars(r)
	p := params["id"]

	err := json.NewDecoder(r.Body).Decode(&postToUp)
	if err != nil {
		log.Fatal(err)
	}

	if post, ok := postStore[p]; ok {
		postToUp.CreatedOn = post.CreatedOn
		delete(postStore, p) //borramos si existe el item y actualizamos el item
		postStore[p] = postToUp
	} else {
		log.Printf("No se pudo encontrar la clave del post %s para eliminar", p)
	}
	w.WriteHeader(http.StatusNoContent)
}

//PostDelete elimina un Post
func PostDelete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	p := params["id"]

	if _, ok := postStore[p]; ok {
		delete(postStore, p)
	} else {
		log.Printf("No se pudo encontrar el key del post para eliminnar %s", p)
	}
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/api/posts", PostGetAll).Methods("GET")
	r.HandleFunc("/api/posts/create", PostCreate).Methods("POST")
	r.HandleFunc("/api/posts/update/{id}", PostUpdate).Methods("PUT")
	r.HandleFunc("/api/posts/delete/{id}", PostDelete).Methods("DELETE")

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	log.Println("Servidor corriendo en http://localhost:8080")
	log.Println(server.ListenAndServe())
}
