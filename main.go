package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	clog "github.com/barelyhuman/commitlog/log"
	git "github.com/go-git/go-git/v5"
)

// Repo - construct containing the details of the cloned repo
type Repo struct {
	Name      string `json:"name"`
	Changelog string `json:"changelog"`
}

var templates *template.Template

func main() {

	var err error
	templates, err = parseTemplates()
	if err != nil {
		log.Println("Error parsing templates, starting only API")
	}

	http.HandleFunc("/generate.json", handleGenerateRequest)
	http.HandleFunc("/generate", handleGenerateRequest)
	http.HandleFunc("/", viewPage)
	http.HandleFunc("/about", viewPage)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	port := ":3000"
	if os.Getenv("PORT") != "" {
		port = ":" + os.Getenv("PORT")
	}
	fmt.Println("Listening on" + port)
	http.ListenAndServe(port, nil)
}

func parseTemplates() (templates *template.Template, err error) {
	const directory = "templates"
	var allFiles []string

	dirContents, _ := ioutil.ReadDir(directory)
	for _, file := range dirContents {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") {
			filePath := filepath.Join(directory, filename)
			allFiles = append(allFiles, filePath)
		}
	}

	templates, err = template.New("").ParseFiles(allFiles...)
	return
}

func viewPage(rw http.ResponseWriter, r *http.Request) {
	pathFile := r.URL.Path[len("/"):]
	if pathFile == "" {
		pathFile = "home"
	}

	rw.Header().Set("Content-Type", "text/html")
	rw.Header().Set("Cache-Control", "no-cache")
	if err := templates.ExecuteTemplate(rw, pathFile+"HTML", nil); err != nil {
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

	rw.Header().Set("Cache-Control", "no-cache")
	http.Redirect(rw, r, "/generate?repo="+shortened, http.StatusSeeOther)
}

func viewGeneratePage(rw http.ResponseWriter, r *http.Request) {
	repo := Repo{}
	asJSON := false

	queryValues := r.URL.Query()

	if r.URL.Path == "/generate.json" {
		asJSON = true
	}

	if len(queryValues["repo"]) > 0 {
		repo.Name = queryValues["repo"][0]
	} else {
		http.Error(rw, fmt.Errorf("no repo received").Error(), http.StatusInternalServerError)
		return
	}

	if !(strings.Contains(repo.Name, "https:") || strings.Contains(repo.Name, "http:")) {
		repo.Name = "https://github.com/" + repo.Name
	}

	os.RemoveAll(path.Join("tmp", repo.Name));

	_, err := git.PlainClone(path.Join("tmp", repo.Name), false, &git.CloneOptions{
		URL: repo.Name,
	})

	if err != nil {
		if !strings.Contains(err.Error(), "repository already exists") {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	rw.Header().Set("Cache-Control", "no-cache")

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

	if asJSON {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(repo)
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
