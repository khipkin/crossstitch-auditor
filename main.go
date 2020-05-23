package main

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
)

// AuditUser prints the JSON encoded "user" field in the body
// of the request or "No user given!" if there isn't one.
func AuditUser(w http.ResponseWriter, r *http.Request) {
	var d struct {
		User string `json:"user"`
	}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprintf(w, "Error decoding request body! %v", err)
		return
	}
	if d.User == "" {
		fmt.Fprint(w, "No user given!")
		return
	}
	fmt.Fprint(w, html.EscapeString(d.User))
}

func main() {
	// Nothing yet.
}
