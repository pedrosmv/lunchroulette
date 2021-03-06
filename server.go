package main

import (
	"fmt"
	"net/http"

	"github.com/pedrosmv/lunchroulette/handlers"

	"goji.io/pat"

	goji "goji.io"

	"github.com/rs/cors"
	mgo "gopkg.in/mgo.v2"
)

/*Index is a struct from the mongo package that groups settings for the DB. This
function makes sure that all of them are true */
func Index(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	context := session.DB("store").C("locations")

	index := mgo.Index{
		Key:        []string{"id"}, //Index key fields; prefix name with dash (-) for descending order
		Unique:     true,           //Prevent two documents from having the same index key
		DropDups:   true,           //Drop documents with the same index key as a previously indexed one
		Background: true,           //Build index in background and return immediately
		Sparse:     true,           //Only index documents containing the Key fields
	}

	err := context.EnsureIndex(index)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		fmt.Println(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	Index(session)

	multiplex := goji.NewMux()
	multiplex.HandleFunc(pat.Post("/locations"), handlers.CreateWrapper(session))
	multiplex.HandleFunc(pat.Get("/locations/:city"), handlers.FetchAll(session))
	multiplex.HandleFunc(pat.Get("/locations/:id"), handlers.ReadWrapper(session))
	multiplex.HandleFunc(pat.Put("/locations/:id"), handlers.UpdateWrapper(session))
	multiplex.HandleFunc(pat.Delete("/locations/:id"), handlers.DeleteWrapper(session))
	handler := cors.Default().Handler(multiplex)
	http.ListenAndServe("localhost:8080", handler)
}
