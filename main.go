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
		var events v1.EventList
		err = json.Unmarshal(body, &events)
		if err != nil {
			http.Error(w, "failed to unmarshal audit events", http.StatusBadRequest)
			return
		}

		for _, event := range events.Items {
			if event.ObjectRef != nil {
				log.Printf("##### %s --- %s --- %s\n", event.ObjectRef.Namespace, event.ObjectRef.Resource, event.Verb)
			}
		}

	})
	log.Println("Starting server at :3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
