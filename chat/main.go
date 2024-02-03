package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/StaphoneWizzoh/ChitChat/trace"
)

// templ represents a single template
type TemplateHandler struct {
	once sync.Once
	filename string
	templ *template.Template

	// Current user field
	currentUser *User
   }
   // ServeHTTP handles the HTTP request.
func (t *TemplateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("../templates", t.filename)))
	})

	// Getting the current user ID from the session
	userId := getAuthenticatedId(r)

	// Retrieve user information based on the user id
	user, exists := GetUser(userId)
	if exists{
		t.currentUser = &user
	}

	t.templ.Execute(w, r)
   }

func main(){
	var addr = flag.String("addr", ":8000", "The addr of the application.")
	flag.Parse()

	// Initializations
	InitDatabase()
	r := newRoom()
		

	// For tracing websocket flow in the application
	r.tracer = trace.New(os.Stdout)

	// http.Handle("/", MustAuth(&TemplateHandler{filename: "index.html"}))
	http.Handle("/chat", MustAuth(&TemplateHandler{filename: "chat.html"}))
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/signup", signupHandler)
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