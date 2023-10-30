package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	v1 "k8s.io/apiserver/pkg/apis/audit/v1"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/audit", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "can't read body", http.StatusBadRequest)
			return
		}
		var event v1.Event
		err = json.Unmarshal(body, &event)
		if err != nil {
			http.Error(w, "failed to unmarshal audit events", http.StatusBadRequest)
			return
		}

		log.Printf("####### %s %s %s %s #######\n", event.Kind, event.Verb, event.ObjectRef.Resource, event.Stage)

	})
	log.Println("Starting server at :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
