package service

import (
	"log"
)

func sendNotification(profileName string, removed, added, currentStocks []string) {
	log.Printf("Sendin notification for profile [%v], removed [%+v], added [%+v], final [%+v]\n", profileName, removed, added, currentStocks)
}
