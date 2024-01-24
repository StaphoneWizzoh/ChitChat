package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

// templ represents a single template
type TemplateHandler struct {
	once sync.Once
	filename string
	templ *template.Template
   }
   // ServeHTTP handles the HTTP request.
func (t *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("../templates", t.filename)))
	})
	t.templ.Execute(w, r)
   }

func main(){
	var addr = flag.String("addr", ":8000", "The addr of the application.")
	flag.Parse()

	r := newRoom()
	// r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&TemplateHandler{filename: "chat.html"}))
	http.Handle("/login", &TemplateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	// Starting the room
	go r.run()

	// Serving static files from the "static" subdirectory
	staticHandler := http.FileServer(http.Dir("../static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))

	// Starting the web server
	log.Printf("Server started at port http://localhost%s",*addr )
	if err := http.ListenAndServe(*addr, nil); err != nil{
		log.Fatal("ListenAndServe:", err)
	}
}