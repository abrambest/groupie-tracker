package main

import (
	"fmt"
	"gtrack/pkg"
	"html/template"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
)

type Error struct {
	Error string
	Code  int
}

func errorPage(w http.ResponseWriter, error string, code int) {
	htmlFiles := []string{
		"./templates/error.html",
		"./templates/base.layout.html",
		"./templates/footer.html",
	}
	tmpl, err := template.ParseFiles(htmlFiles...)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
	w.WriteHeader(code)
	err = tmpl.ExecuteTemplate(w, "error.html", Error{
		Error: error,
		Code:  code,
	})
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		return
	}
}

func artistPage(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		errorPage(w, "Page not found", 404)
		return
	}
	htmlFiles := []string{
		"./templates/artist.html",
		"./templates/base.layout.html",
		"./templates/footer.html",
	}
	if r.Method != http.MethodGet {
		errorPage(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	tmpl, err := template.ParseFiles(htmlFiles...)
	if err != nil {
		errorPage(w, "Internal Server Error", 500)
		return
	}
	data, err := pkg.Parser()

	if err != nil {
		fmt.Println("error pars")
	}
	err = pkg.ParsRelation(id)
	if err != nil {
		fmt.Println("error pars dates and location")
		return
	}
	err = tmpl.Execute(w, data[id-1])
	if err != nil {
		errorPage(w, "Internal Server Error", 500)
		return
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorPage(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	htmlFiles := []string{
		"./templates/home.html",
		"./templates/base.layout.html",
		"./templates/footer.html",
	}
	if r.Method != http.MethodGet {
		errorPage(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	tmpl, err := template.ParseFiles(htmlFiles...)
	if err != nil {
		errorPage(w, "Internal Server Error", 500)
		return
	}
	data, err := pkg.Parser()
	if err != nil {
		errorPage(w, "Error read JSON", 500)
		return
	}
	err = tmpl.Execute(w, data)

	if err != nil {

		errorPage(w, "Internal Server Error", 500)
		return
	}
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/artist/", artistPage)
	fileServer := http.FileServer(neuteredFileSystem{http.Dir("./static/")})
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	log.Println("Server start: http://localhost:8080")
	err := http.ListenAndServe("localhost:8080", mux)
	log.Fatal(err)
}

type neuteredFileSystem struct {
	fs http.FileSystem
}

func (nfs neuteredFileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := f.Stat()
	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := nfs.fs.Open(index); err != nil {
			closeErr := f.Close()
			if closeErr != nil {
				return nil, closeErr
			}
			return nil, err
		}
	}
	return f, nil
}
