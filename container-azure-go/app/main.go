package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type AppResponse struct {
	Message string `json:"message"`
}

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := AppResponse{
			Message: "Hello, world! 👋",
		}

		s, err := json.Marshal(response)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, string(s))
	})

	if err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil); err != nil {
		log.Fatal(err)
	}
}
