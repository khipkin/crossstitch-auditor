package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/khipkin/geddit"
)

const (
	subreddit      = "CrossStitch"
	redditClientID = "Kkfhbwt2W5C0Rw"
	redditUsername = "CrossStitchBot"
)

type auditor struct {
	redditSession *geddit.OAuthSession
}

func newAuditor() (*auditor, error) {
	// Authenticate with Reddit.
	redditClientSecret := os.Getenv("REDDIT_CLIENT_SECRET")
	if redditClientSecret == "" {
		log.Print("REDDIT_CLIENT_SECRET not set")
		return nil, errors.New("REDDIT_CLIENT_SECRET not set")
	}
	redditSession, err := geddit.NewOAuthSession(
		redditClientID,
		redditClientSecret,
		"gedditAgent v1 fork by khipkin",
		"redirect.url",
	)
	if err != nil {
		log.Printf("Failed to create new Reddit OAuth session: %v", err)
		return nil, err
	}
	redditPassword := os.Getenv("REDDIT_PASSWORD")
	if redditPassword == "" {
		log.Print("REDDIT_PASSWORD not set")
		return nil, errors.New("REDDIT_PASSWORD not set")
	}
	if err = redditSession.LoginAuth(redditUsername, redditPassword); err != nil {
		log.Printf("Failed to authenticate with Reddit: %v", err)
		return nil, err
	}

	// To prevent Reddit rate limiting errors, throttle requests.
	redditSession.Throttle(5 * time.Second)

	return &auditor{
		redditSession: redditSession,
	}, nil
}

func (a *auditor) auditUser(user string) string {
	return user
}

// AuditUser processes the JSON encoded "user" field in the body
// of the request or prints "No user given!" if there isn't one.
// Prints any errors encountered during execution.
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

	a, err := newAuditor()
	if err != nil {
		fmt.Fprintf(w, "Error setting up auditor! %v", err)
		return
	}
	fmt.Fprint(w, a.auditUser(html.EscapeString(d.User)))
}

// main processed the given REDDIT_USER.
// Prints any errors encountered during execution.
func main() {
	user := os.Getenv("REDDIT_USER")
	if user == "" {
		log.Fatal("REDDIT_USER not set")
	}

	a, err := newAuditor()
	if err != nil {
		log.Fatalf("Error setting up auditor! %v", err)
	}
	log.Printf(a.auditUser(user))
}
