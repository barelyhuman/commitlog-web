package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"

	clog "github.com/barelyhuman/commitlog/log"
	git "github.com/go-git/go-git/v5"
)

type Repo struct {
	Name      string
	Changelog string
}

func main() {

	http.HandleFunc("/generate", handleGenerateRequest)
	http.HandleFunc("/", viewHomePage)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	port := ":3000"
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}
	fmt.Println("Listening on" + port)
	http.ListenAndServe(port, nil)
}

func viewHomePage(rw http.ResponseWriter, r *http.Request) {
	htmlFile := path.Join("templates", "home.html")
	tmpl, err := template.ParseFiles(htmlFile)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(rw, nil); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func handleGenerateRequest(rw http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		viewGeneratePage(rw, r)
		break
	case http.MethodPost:
		generateCommitlog(rw, r)
		break
	default:
		http.Error(rw, fmt.Errorf("Error").Error(), http.StatusNotFound)
		break
	}
}

func generateCommitlog(rw http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	form := r.Form

	urlToClone := form["github-url"][0]

	shortened := urlToClone
	if strings.Contains(urlToClone, "https://github.com/") {
		shortened = urlToClone[len("https://github.com/"):]
	}

	_, err := git.PlainClone(path.Join("tmp", shortened), false, &git.CloneOptions{
		URL: urlToClone,
	})

	if err != nil {
		if !strings.Contains(err.Error(), "repository already exists") {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	rw.Header().Set("Cache-Control", "no-cache")
	http.Redirect(rw, r, "/generate?repo="+shortened, http.StatusSeeOther)
}

func viewGeneratePage(rw http.ResponseWriter, r *http.Request) {
	repo := Repo{}

	queryValues := r.URL.Query()

	if len(queryValues["repo"]) > 0 {
		repo.Name = queryValues["repo"][0]
	} else {
		http.Error(rw, fmt.Errorf("no repo received").Error(), http.StatusInternalServerError)
		return
	}

	if _, err := os.Stat(path.Join("tmp", repo.Name)); os.IsNotExist(err) {
		http.Error(rw, fmt.Errorf("Couldn't find repository").Error(), http.StatusInternalServerError)
		return
	}
	var clogError clog.ErrMessage
	repo.Changelog, clogError = clog.CommitLog(path.Join("tmp", repo.Name), "", "", "ci,refactor,docs,fix,feat,test,chore,other", false)

	if clogError.Err != nil {
		http.Error(rw, clogError.Message, http.StatusInternalServerError)
		return
	}

	htmlFile := path.Join("templates", "generate.html")
	tmpl, err := template.ParseFiles(htmlFile)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(rw, repo); err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
