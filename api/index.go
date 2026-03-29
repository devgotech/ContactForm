package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
)

type FormData struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	// Allow CORS for your Framer domain
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var data FormData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := sendEmail(data); err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func sendEmail(data FormData) error {
	gmailUser := os.Getenv("GMAIL_USER")
	gmailPass := os.Getenv("GMAIL_APP_PASS")

	auth := smtp.PlainAuth("", gmailUser, gmailPass, "smtp.gmail.com")

	body := fmt.Sprintf(
		"Subject: New Portfolio Message from %s\r\n\r\nName: %s\r\nEmail: %s\r\nMessage: %s",
		data.Name, data.Name, data.Email, data.Message,
	)

	return smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		gmailUser,
		[]string{gmailUser},
		[]byte(body),
	)
}
