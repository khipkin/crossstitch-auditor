package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/khipkin/geddit"
)

const (
	subreddit      = "CrossStitch" // TODO: paremterize this
	redditClientID = "Kkfhbwt2W5C0Rw"
	redditUsername = "CrossStitchBot" // TODO: paremterize this
)

type auditor struct {
	redditSession *geddit.OAuthSession
}

type comment struct {
	Link string `json:"link"`
	Body string `json:"body"`
}

type post struct {
	Link     string     `json:"link"`
	Title    string     `json:"title"`
	Comments []*comment `json:"comments"`
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

func (a *auditor) auditUser(user string) ([]*post, error) {
	params := geddit.ListingOptions{
		Limit: 100, // the max value
	}
	posts, err := a.redditSession.UserPosts(subreddit, user, geddit.NewSubmissions, params)
	if err != nil {
		return nil, err
	}
	comments, err := a.redditSession.UserComments(subreddit, user, geddit.NewSubmissions, params)
	if err != nil {
		return nil, err
	}

	// Build the auditMap, keyed by link fullID. "" is a placeholder for comments without a matching post.
	auditMap := map[string]*post{
		"": &post{Comments: []*comment{}},
	}
	// TODO: filter out posts that have been removed by a moderator
	for _, p := range posts {
		auditMap[p.FullID] = &post{
			Link:     p.FullPermalink(),
			Title:    p.Title,
			Comments: []*comment{},
		}
	}
	// TODO: filter out comments on posts that have been removed by a moderator
	for _, c := range comments {
		parent := auditMap[""]
		if p, ok := auditMap[c.LinkID]; ok {
			parent = p
		}
		parent.Comments = append(parent.Comments, &comment{
			Link: c.FullPermalink(),
			Body: c.Body,
		})
	}

	ret := []*post{}
	for _, p := range auditMap {
		ret = append(ret, p)
	}
	return ret, nil
}

func buildAuditString(postData []*post) string {
	audit := ""
	for _, p := range postData {
		if p.Link == "" {
			audit += "\n**(comments on other people's posts)**\n\n"
		} else {
			audit += fmt.Sprintf("\n[**%s**](%s)\n\n", p.Title, p.Link)
		}

		for _, c := range p.Comments {
			audit += fmt.Sprintf(" *  [`%s`](%s)\n", strings.Replace(c.Body, "\n", " ", -1 /*unlimited*/), c.Link)
		}
	}
	return audit
}

// AuditUser is the function that should be invoked if this code is
// hosted as a Cloud Function on Google Cloud.
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
	posts, err := a.auditUser(html.EscapeString(d.User))
	if err != nil {
		fmt.Fprintf(w, "Error building post data! %v", err)
		return
	}

	result, err := json.Marshal(posts)
	if err != nil {
		fmt.Fprintf(w, "Error marshalling response to JSON! %v", err)
	}
	fmt.Fprint(w, result)
}

// main is the function invoked if this code is run from the command line.
// main processes the given REDDIT_USER.
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
	postData, err := a.auditUser(user)
	if err != nil {
		log.Fatalf("Error building post data! %v", err)
	}
	log.Printf(buildAuditString(postData))
}
