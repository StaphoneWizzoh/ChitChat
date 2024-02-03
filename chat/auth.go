package main

import (
	"log"
	"net/http"
)

type authHandler struct{
	next http.Handler
}

func (h * authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	// _, err := r.Cookie("auth")
	// if err == http.ErrNoCookie{
	// 	// User not authenticated
	// 	log.Println("Redirecting since user is not authenticated.")
	// 	w.Header().Set("Location", "/login")
	// 	w.WriteHeader(http.StatusTemporaryRedirect)
	// 	return
	// }

	// if err != nil {
	// 	// Occurence of some other error
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// check if the user is already authenticated via local storage
	if token := getTokenFromLocalStorage(r); token == "" {
		// User is not authenticated, redirect to chat page
		log.Println("Redirecting since user is not authenticated")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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
// func loginHandler(w http.ResponseWriter, r * http.Request){
// 	// TODO: handle case where the handler is called with few segments 
// 	segs := strings.Split(r.URL.Path, "/")
// 	action := segs[2]
// 	provider := segs[3]
// 	if provider == ""{
// 		log.Println("The url path triggered didn't provide a provider service.")
// 	}
// 	switch action{
// 	case "login":
// 		log.Println("TODO handle login for", provider)
// 	case "signup":

// 	default:
// 		w.WriteHeader(http.StatusNotFound)
// 		fmt.Fprintf(w, "Auth action %s not supported", action)
// 	}
// }

func getAuthenticatedId(r *http.Request)string{
	// Implement the logic to retrieve the user ID from the session or token.
    // This may involve checking cookies, headers, or any other authentication mechanism.
    // For simplicity, let's assume a cookie named "userID" is used.
    userIDCookie, err := r.Cookie("auth")
	if err == http.ErrNoCookie{
		return ""
	}
    return userIDCookie.Value
}