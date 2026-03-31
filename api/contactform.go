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

	// Send notification email to you
	if err := sendNotification(data); err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	// Send auto-reply to the person
	sendAutoReply(data)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]bool{"success": true})
}

func sendNotification(data FormData) error {
	gmailUser := os.Getenv("GMAIL_USER")
	gmailPass := os.Getenv("GMAIL_APP_PASS")

	auth := smtp.PlainAuth("", gmailUser, gmailPass, "smtp.gmail.com")

	// MIME headers for HTML email + Reply-To
	headers := fmt.Sprintf(
		"From: Portfolio Contact <%s>\r\n"+
			"To: %s\r\n"+
			"Reply-To: %s <%s>\r\n"+
			"Subject: New message from %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n",
		gmailUser, gmailUser, data.Name, data.Email, data.Name,
	)

	body := fmt.Sprintf(`
	<div style="font-family: Inter, sans-serif; max-width: 600px; margin: 0 auto; padding: 40px 24px; color: #1a1a1a;">
		<h2 style="font-size: 20px; font-weight: 600; margin-bottom: 32px; border-bottom: 1px solid #EAEAEA; padding-bottom: 16px;">
			New Portfolio Message
		</h2>
		<table style="width: 100%%; border-collapse: collapse;">
			<tr>
				<td style="padding: 12px 0; color: #7a7a7a; font-size: 13px; width: 80px; vertical-align: top;">Name</td>
				<td style="padding: 12px 0; font-size: 14px; font-weight: 500;">%s</td>
			</tr>
			<tr style="border-top: 1px solid #EAEAEA;">
				<td style="padding: 12px 0; color: #7a7a7a; font-size: 13px; vertical-align: top;">Email</td>
				<td style="padding: 12px 0; font-size: 14px;">
					<a href="mailto:%s" style="color: #1a1a1a;">%s</a>
				</td>
			</tr>
			<tr style="border-top: 1px solid #EAEAEA;">
				<td style="padding: 12px 0; color: #7a7a7a; font-size: 13px; vertical-align: top;">Message</td>
				<td style="padding: 12px 0; font-size: 14px; line-height: 1.6;">%s</td>
			</tr>
		</table>
		<p style="margin-top: 32px; font-size: 12px; color: #7a7a7a;">
			Hit reply to respond directly to %s
		</p>
	</div>
	`, data.Name, data.Email, data.Email, data.Message, data.Name)

	return smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		gmailUser,
		[]string{gmailUser},
		[]byte(headers+body),
	)
}

func sendAutoReply(data FormData) {
	gmailUser := os.Getenv("GMAIL_USER")
	gmailPass := os.Getenv("GMAIL_APP_PASS")

	auth := smtp.PlainAuth("", gmailUser, gmailPass, "smtp.gmail.com")

	// Replace the name and sign-off below with your own name
	headers := fmt.Sprintf(
		"From: Gotech <%s>\r\n"+
			"To: %s <%s>\r\n"+
			"Subject: Thanks for reaching out!\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n\r\n",
		gmailUser, data.Name, data.Email,
	)

	body := fmt.Sprintf(`
	<div style="font-family: Inter, sans-serif; max-width: 600px; margin: 0 auto; padding: 40px 24px; color: #1a1a1a;">
		<h2 style="font-size: 20px; font-weight: 600; margin-bottom: 24px;">
			Hi %s,
		</h2>
		<p style="font-size: 14px; line-height: 1.8; color: #1a1a1a;">
			Thanks for reaching out! I've received your message and will get back to you as soon as possible.
		</p>
		<p style="font-size: 14px; line-height: 1.8; color: #1a1a1a;">
			In the meantime, feel free to check out my work or connect with me on LinkedIn.
		</p>
		<div style="margin: 32px 0; padding: 20px 24px; background: #EAEAEA; border-left: 3px solid #1a1a1a;">
			<p style="margin: 0; font-size: 13px; color: #7a7a7a; margin-bottom: 8px;">Your message</p>
			<p style="margin: 0; font-size: 14px; line-height: 1.6;">%s</p>
		</div>
		<p style="font-size: 14px; line-height: 1.8; margin-top: 32px;">
			Best regards,<br/>
			<strong>Gotech</strong>
		</p>
		<p style="font-size: 12px; color: #7a7a7a; margin-top: 32px; border-top: 1px solid #EAEAEA; padding-top: 16px;">
			This is an automated reply. Please do not reply to this email.
		</p>
	</div>
	`, data.Name, data.Message)

	// We don't return the error here — if auto-reply fails,
	// the main notification already succeeded so we still return success
	smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		gmailUser,
		[]string{data.Email},
		[]byte(headers+body),
	)
}