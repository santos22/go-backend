package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Bash to set ENV vars from file
// for line in $(cat github.env)
// do
// export $line
// done

var (
	accessToken = getenv("GITHUB_ACCESS_TOKEN")
)

// App comment
type App struct {
	client *github.Client
	ctx    context.Context
	// Use MUX?
}

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

func (app *App) handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Howdy! Visit one of the two routes: /closed or /open")
}

func (app *App) closedPullRequests(w http.ResponseWriter, r *http.Request) {
	opts := &github.SearchOptions{Sort: "created", Order: "asc"}
	issues, _, err := app.client.Search.Issues(app.ctx, "state:closed author:santos22 type:pr", opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Howdy! Here are your GitHub issues stringified: %s", github.Stringify(issues.Issues))
}

func (app *App) openPullRequests(w http.ResponseWriter, r *http.Request) {
	opts := &github.SearchOptions{Sort: "created", Order: "asc"}
	issues, _, err := app.client.Search.Issues(app.ctx, "state:open author:santos22 type:pr", opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Howdy! I love GitHub and there are currently %d pull requests opened.", len(issues.Issues))
}

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	app := &App{ctx: ctx, client: client}

	http.HandleFunc("/", app.handler)
	http.HandleFunc("/closed", app.closedPullRequests)
	http.HandleFunc("/open", app.openPullRequests)
	http.ListenAndServe(":8080", nil)
}
