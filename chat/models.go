package main

import (
	"errors"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

type User struct{
	gorm.Model
	Username 		string
	Email 			string
	Password 		string
	Messages 		[]Message 		`gorm:"foreignKey:SenderID"`
}

type Message struct{
	gorm.Model
	Content 		string 
	SenderID 		uint 
	Time 			time.Time
}

// Initialize database

func InitDatabase(){
	// Connect to SQLite database

	var err error
	db, err = gorm.Open(sqlite.Open("../chitchat.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	// Auto Migrate model structs
	db.AutoMigrate(&User{}, &Message{}) 
}

// User type related methods

func (user *User) SetPassword(password string) error{
	hashedPassword, err := hashString(password)
	if err != nil{
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

func (user *User) CheckPassword(password string) error{
	inputHash, err := hashString(password)
	if err != nil{
		log.Println("Error in hashing during password verification")
		return err
	}

	if user.Password != inputHash{
		return errors.New("passwords do not match")
	}
	return nil
}

// User type related functions

func CreateUser(username, email, password string) (registrationError error){
	var user User

	// Check if User data exists
	result := db.Where("username = ? AND email >= ?", username, email).First(&user)
	if result.RowsAffected == 0{
		// User does not exist
		// Create user
		hashPassword, err := hashString(password)
		if err != nil{
			registrationError = err
		}
		newUser := User{
			Username: username,
			Email: email,
			Password: hashPassword,
			Messages: []Message{},
		}
		db.Create(&newUser)
	}else{
		registrationError = errors.New("User already exists")
	}

	return registrationError
}

func GetUser(userId string)(User, bool){
	user, exists:= User{}, false
	return user, exists
}

// Message type related functions
func SaveMessage(content string, senderID uint){
	message := Message{
		Content: content,
		SenderID: senderID,
		Time: time.Now(),
	}
	db.Create(&message)
	// TODO: Handle possible errors
}