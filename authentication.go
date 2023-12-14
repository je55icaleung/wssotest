package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Prompt for user input
	fmt.Print("Enter username: ")
	username := readInput()

	fmt.Print("Enter password: ")
	password := readInput()

	fmt.Print("Enter email: ")
	email := readInput()

	// Open the database connection
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create the users table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		password TEXT,
		email TEXT,
		verification_code TEXT,
		verified INTEGER DEFAULT 0
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// Register the user
	err = RegisterUser(db, username, password, email)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Enter verification code: ")
	verificationCode := readInput()

	// Verify the user
	err = VerifyUser(db, email, verificationCode)
	if err != nil {
		log.Fatal(err)
	}
}

// RegisterUser registers a new user account
func RegisterUser(db *sql.DB, username, password, email string) error {
	// Generate a random 6-digit verification code
	verificationCode := generateVerificationCode()

	// Save the user details and verification code to the database
	_, err := db.Exec("INSERT INTO users (username, password, email, verification_code) VALUES (?, ?, ?, ?)",
		username, password, email, verificationCode)
	if err != nil {
		return err
	}

	// Send the verification code to the user's email
	sendVerificationCode(email, verificationCode)

	return nil
}

// VerifyUser verifies a user's email by matching the verification code
func VerifyUser(db *sql.DB, email, verificationCode string) error {
	// Check if the verification code matches the one stored in the database
	var storedVerificationCode string
	err := db.QueryRow("SELECT verification_code FROM users WHERE email = ?", email).Scan(&storedVerificationCode)
	if err != nil {
		return err
	}

	if verificationCode != storedVerificationCode {
		return fmt.Errorf("verification code does not match")
	}

	// Update the user's account to mark it as verified
	_, err = db.Exec("UPDATE users SET verified = 1 WHERE email = ?", email)
	if err != nil {
		return err
	}

	return fmt.Errorf("User successfully registered!")
}

// readInput reads user input from the console
func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// sendVerificationCode sends the verification code to the user's email
// sendVerificationCode sends the verification code to the user's email
func sendVerificationCode(body string, verificationCode string) {
	from := "skinnywsso@gmail.com"
	pass := "bzei uxxz ecef sdmi"
	to := "jessicaemail1@gmail.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		"Verification Code: " + verificationCode + "\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
	log.Println("Successfully sent to " + to)
}

// generateVerificationCode generates a random 6-digit verification code
func generateVerificationCode() string {
	// Define the characters to be used in the verification code
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

	// Generate a random verification code with 6 characters
	code := make([]byte, 6)
	for i := range code {
		code[i] = characters[rand.Intn(len(characters))]
	}

	return string(code)
}
