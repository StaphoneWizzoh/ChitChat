package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

func hashString(stringToHash string)(string, error){
	hash := sha256.New()
	hash.Write([]byte(stringToHash))
	hashedBytes := hash.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)
	return hashedString, nil
}

// Fetches token from local storage
func getTokenFromLocalStorage(r *http.Request) string{
	cookie, err := r.Cookie("auth")
	if err != http.ErrNoCookie{
		return cookie.Value
	}

	return ""
}

// Sets the authentication token in the local storage
func setTokenInLocalStorage(w http.ResponseWriter, username string){
	// TODO: Generate a secure and unique token
	token := generateAuthToken(username)
	
	// Set the auth cookie
	http.SetCookie(w, &http.Cookie{
		Name: "auth",
		Value: token,

		// TODO: Set other cookie parameters if need be
	})

	// Use javascript to set auth token
	// TODO: Look for an alternative way
	setTokenScript := `
		<script>
			localStorage.setItem("` + localStorageTokenKey + `", "`+token+`");
		</script>
	`
	fmt.Fprint(w, setTokenScript)
}

// Generates secure authentication token
func generateAuthToken(username string) string{
	// TODO: use aJWT library to implement

	// token := jwt.NewWithClaims(jwt.SigninMethodHS256, jwt.MapClaims{
	// 	"username": username,
	// })
	// signedToken, err := token.SignedString([]byte("your-secret-key"))
	// return signedToken 

	return generateUniqueString(username)
}

// Amateur implemetation to generate unique string
// For testing purposes
func generateUniqueString(input string) string{
	timestamp := time.Now().UnixNano()
	uniqueString := fmt.Sprintf("%s_%d", input, timestamp)
	return uniqueString
}