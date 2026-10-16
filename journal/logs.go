package journal

import (
	"fmt"
	"os"
)

func EcrireLog(tour int, action string) {
	// O_APPEND : ajoute à la fin du fichier sans effacer
	// O_CREATE : crée le fichier s'il n'existe pas encore
	// O_WRONLY : ouvre en écriture seule
	fichier, err := os.OpenFile("journal.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer fichier.Close()

	// Écrit la nouvelle ligne à la suite des autres
	fmt.Fprintf(fichier, "Tour %d: Action %s\n", tour, action)
}

func SupLogs() {
	// Ramène la taille du fichier à 0 octet sans le supprimer
	err := os.Truncate("journal.txt", 0)
	if err != nil {
		fmt.Println("Erreur lors du vidage :", err)
	}
}