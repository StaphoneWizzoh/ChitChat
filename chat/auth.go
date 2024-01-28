package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type authHandler struct{
	next http.Handler
}

func (h * authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	_, err := r.Cookie("auth")
	if err == http.ErrNoCookie{
		// User not authenticated
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	if err != nil {
		// Occurence of some other error
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Successful authentication - Calling next handler
	h.next.ServeHTTP(w,r)
}

func MustAuth(handler http.Handler) http.Handler{
	return &authHandler{next: handler}
}

// loginHandler handles the third-party login process.
// Format: /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r * http.Request){
	// TODO: handle case where the handler is called with few segments 
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	if provider == ""{
		log.Println("The url path triggered didn't provide a provider service.")
	}
	switch action{
	case "login":
		log.Println("TODO handle login for", provider)
	case "signup":

	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}