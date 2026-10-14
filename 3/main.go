package main

import (
	"cyoa/data"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func main() {

	stories, err := data.ReadAndParseJson("gopher.json")
	if err != nil {
		fmt.Println("Error accessing JSON:", err)
	}

	// TODO: look into template.ParseFS
	app := &App{
		tmpl: template.Must(template.ParseGlob("templates/*.html")),
		strs: &stories,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.homepageHandler)
	mux.HandleFunc("GET /arc/{id}", app.arcHandler)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Starting the server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

type App struct {
	tmpl *template.Template
	strs *data.Arcs
}

func (a *App) homepageHandler(w http.ResponseWriter, r *http.Request) {
	a.tmpl.ExecuteTemplate(w, "index.html", a.strs)
}

func (a *App) arcHandler(w http.ResponseWriter, r *http.Request) {
	arcId := r.PathValue("id")
	arc := (*a.strs)[arcId]

	log.Printf("Loading /%s\n", arcId)

	a.tmpl.ExecuteTemplate(w, "arc.html", arc)
}
