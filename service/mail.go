package service

import (
	"fmt"
	"log"
	"net/smtp"

	"github.com/nagymarci/stock-screener/config"
)

// smtpServer data to smtp server
type smtpServer struct {
	host string
	port string
}

// Address URI to smtp server
func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

func sendNotification(profileName string, removed, added, currentStocks []string) error {
	return sendNotification_watchlist(profileName, removed, added, currentStocks, config.Config.NotificationRecipient)
}

func sendNotification_watchlist(profileName string, removed, added, currentStocks []string, email string) error {
	log.Printf("Sendin notification for profile [%v], removed [%+v], added [%+v], final [%+v]\n", profileName, removed, added, currentStocks)

	// Sender data.
	from := config.Config.SmptSenderUsername
	password := config.Config.SmptSenderPassword
	// Receiver email address.
	to := []string{
		email,
	}
	// smtp server configuration.
	smtpServer := smtpServer{host: config.Config.SmptServerHost, port: config.Config.SmptServerPort}
	// Message.
	message := []byte(fmt.Sprintf("To: %v\n"+
		"Subject: %v changed!\n"+
		"\n"+
		"Recommendations in your %v profile has changed.\n\n"+
		"Removed stocks: %+v\n"+
		"Added stocks: %+v\n\n"+
		"Currently recommended stocks: %+v", to[0], profileName, profileName, removed, added, currentStocks))
	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpServer.host)
	// Sending email.
	err := smtp.SendMail(smtpServer.Address(), auth, from, to, message)
	if err == nil {
		log.Println("Email Sent!")
	}

	return err
}
