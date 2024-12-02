// {
// 	"subject": "O titulo",
// 	"body_html": "<p>O corpo do email</p>",
// 	"sender": "igor@gtrinvestimentos.com.br",
// 	"recipients": [
// 		"igor@gtrinvestimentos.com.br"
// 	]
// }

package external

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "net/mail"
    "os"
    "time"

    "github.com/mailersend/mailersend-go"
)

type EmailRequest struct {
    Subject    string   `json:"subject"`
    BodyHTML   string   `json:"body_html"`
    Sender     string   `json:"sender"`
    Recipients []string `json:"recipients"`
}

func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
    // Only allow POST requests
    if r.Method != http.MethodPost {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }

    // Parse JSON request body
    var emailReq EmailRequest
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&emailReq)
    if err != nil {
        http.Error(w, "Bad Request: Invalid JSON", http.StatusBadRequest)
        return
    }

    // Validate required fields
    if emailReq.Subject == "" || emailReq.BodyHTML == "" || emailReq.Sender == "" || len(emailReq.Recipients) == 0 {
        http.Error(w, "Bad Request: Missing required fields", http.StatusBadRequest)
        return
    }

    // Validate sender email
    _, err = mail.ParseAddress(emailReq.Sender)
    if err != nil {
        http.Error(w, "Bad Request: Invalid sender email address", http.StatusBadRequest)
        return
    }

    // Create an instance of the MailerSend client
    apiKey := os.Getenv("MAILERSEND_API_KEY")
    if apiKey == "" {
        http.Error(w, "Server Error: Missing MailerSend API key", http.StatusInternalServerError)
        return
    }
    ms := mailersend.NewMailersend(apiKey)

    // Prepare the sender
    from := mailersend.From{
        Email: emailReq.Sender,
        Name:  "", // No name provided
    }

    // Initialize a slice to hold message IDs
    var messageIDs []string

    // Set a context with timeout
    ctx := context.Background()
    ctx, cancel := context.WithTimeout(ctx, time.Minute)
    defer cancel()

    // Loop over recipients and send individual emails
    for _, recipientEmail := range emailReq.Recipients {
        // Validate recipient email
        _, err = mail.ParseAddress(recipientEmail)
        if err != nil {
            // Skip invalid email addresses
            continue
        }

        // Create the email message
        message := ms.Email.NewMessage()
        message.SetFrom(from)
        message.SetSubject(emailReq.Subject)
        message.SetHTML(emailReq.BodyHTML)

        // Set the recipient
        recipient := mailersend.Recipient{
            Email: recipientEmail,
            Name:  "", // No name provided
        }
        message.SetRecipients([]mailersend.Recipient{recipient})

        // Send the email
        res, err := ms.Email.Send(ctx, message)
        if err != nil {
            // Handle send error (log, continue, or return error)
            fmt.Printf("Failed to send email to %s: %v\n", recipientEmail, err)
            continue
        }

        // Append the message ID to the list
        messageID := res.Header.Get("X-Message-Id")
        messageIDs = append(messageIDs, messageID)
    }

    // Check if any emails were sent
    if len(messageIDs) == 0 {
        http.Error(w, "Failed to send emails to any recipients", http.StatusInternalServerError)
        return
    }

    // Respond with success and the list of message IDs
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    response := map[string]interface{}{
        "message":     "Emails sent successfully",
        "message_ids": messageIDs,
    }
    json.NewEncoder(w).Encode(response)
}
