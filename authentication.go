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
	fmt.Print("Enter username: ")
	username := readInput()

	fmt.Print("Enter password: ")
	password := readInput()

	fmt.Print("Enter email: ")
	email := readInput()

	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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

	err = RegisterUser(db, username, password, email)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Enter verification code: ")
	verificationCode := readInput()

	err = VerifyUser(db, email, verificationCode)
	if err != nil {
		log.Fatal(err)
	}
}

func RegisterUser(db *sql.DB, username, password, email string) error {
	var existingUsername string
	err := db.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&existingUsername)
	if err == nil {
		return fmt.Errorf("This user already exists.")
	}

	var existingEmail string
	err = db.QueryRow("SELECT email FROM users WHERE email = ?", email).Scan(&existingEmail)
	if err == nil {
		return fmt.Errorf("This user already exists.")
	}

	verificationCode := generateVerificationCode()

	_, err = db.Exec("INSERT INTO users (username, password, email, verification_code) VALUES (?, ?, ?, ?)",
		username, password, email, verificationCode)
	if err != nil {
		return err
	}

	sendVerificationCode(email, verificationCode)

	return nil
}

func VerifyUser(db *sql.DB, email, verificationCode string) error {
	var storedVerificationCode string
	err := db.QueryRow("SELECT verification_code FROM users WHERE email = ?", email).Scan(&storedVerificationCode)
	if err != nil {
		return err
	}

	if verificationCode != storedVerificationCode {
		return fmt.Errorf("verification code does not match")
	}

	_, err = db.Exec("UPDATE users SET verified = 1 WHERE email = ?", email)
	if err != nil {
		return err
	}

	return fmt.Errorf("User successfully registered!")
}

func readInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

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

func generateVerificationCode() string {
	characters := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	code := make([]byte, 6)
	for i := range code {
		code[i] = characters[rand.Intn(len(characters))]
	}

	return string(code)
}
