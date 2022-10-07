package api

import (
	"fmt"
	"net/http"
	"time"
)

func Data(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{ "now": %d }`, time.Now().UnixNano()/1000000)
}
