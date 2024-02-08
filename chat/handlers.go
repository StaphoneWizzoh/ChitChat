package main

import (
	"log"
	"net/http"
)

const localStorageTokenKey = "authToken"

func signupHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodPost{
		// Retrieve form data
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Create a new user
		err := CreateUser(username, email, password)
		if err != nil {
			log.Println("Error creating user:", err)
			// TODO: Handle error
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		// Redirect the user to the login page after successful registration
		http.Redirect(w,r,"/login", http.StatusSeeOther)
		return
	}else{
		// Rendering the signup form template
		signupTemplate := &TemplateHandler{filename: "signup.html"}
		signupTemplate.ServeHTTP(w,r)
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodPost{
		// Retrieve form data
		username := r.FormValue("username")
		password := r.FormValue("password")

		// check if the user is already authenticated via local storage
		if token := getTokenFromLocalStorage(r); token != "" {
			// User is authenticated, redirect to chat page
			http.Redirect(w, r, "/chat", http.StatusSeeOther)
			return
		}

		var user User
		result := db.Where("username = ?", username).First(&user)
		if result.RowsAffected == 0 {
			// User not found 
			// http.Error(w, "User does not Exist", http.StatusNotFound)
			log.Println("User does not Exist", http.StatusNotFound)
			http.Redirect(w,r,"/signup", http.StatusSeeOther)
			return
		}

		// Check if the provided password is correct
		if err := user.CheckPassword(password); err != nil{
			// Password does not match
			http.Error(w,"Incorrect password", http.StatusUnauthorized)
			return
		}

		// Set the authenticationtoken in local storage
		// setTokenInLocalStorage(w, username)

		// Redirect the user to the chat page after successful login
		http.Redirect(w,r,"/chat", http.StatusSeeOther)
	}else{
		// Render the login form template
		loginTemplate := &TemplateHandler{filename: "login.html"}
		loginTemplate.ServeHTTP(w,r)
	}
}


