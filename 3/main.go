package main

import (
	"cyoa/data"
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
)

//go:embed templates/* static/*
var appFiles embed.FS

/*
 - go embed above tells Go to include the matching HTML files when compiling
 - They're stored in the compiled executable
 - ParseFS reads the stored templates from files.
 - ParseGlob will look for actual directories, and if they're not available where the executable is, things break

 - you can have separate embeds for templates & static if needed
*/

func main() {

	stories, err := data.ReadAndParseJson("gopher.json")
	if err != nil {
		log.Fatal("Error accessing JSON:", err)
	}

	app := &App{
		tmpl: template.Must(template.ParseFS(appFiles, "templates/*.html")),
		strs: &stories,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", app.homepageHandler)
	mux.HandleFunc("GET /arc/{id}", app.arcHandler)

	staticFS, err := fs.Sub(appFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

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
	if arc, ok := (*a.strs)[arcId]; ok {
		log.Printf("Loading /%s\n", arcId)

		a.tmpl.ExecuteTemplate(w, "arc.html", arc)
		return
	}

	http.Error(w, "An error occured", http.StatusNotFound)
}
