package main

import (
	"crypto/sha256"
	"encoding/hex"
)

func hashString(stringToHash string)(string, error){
	hash := sha256.New()
	hash.Write([]byte(stringToHash))
	hashedBytes := hash.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)
	return hashedString, nil
}