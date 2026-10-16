package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	
	"github.com/votre-pseudo/arena/carte"
	"github.com/votre-pseudo/arena/journal"
	"github.com/votre-pseudo/arena/protocol"
	"github.com/votre-pseudo/arena/robot"
)

func main() {
	// 1. Récupération des arguments envoyés par l'arbitre (--carte et --joueur)
	cheminCarte := flag.String("carte", "", "fichier .map de la partie")
	joueur := flag.Int("joueur", 0, "numéro du joueur (1 ou 2)")
	flag.Parse()

	// 2. Affichage de contrôle sur os.Stderr (os.Stdout est réservé à l'arbitre !)
	fmt.Fprintf(os.Stderr, "[ROBOT] Joueur %d démarré avec la carte : %s\n", *joueur, *cheminCarte)

	// TODO : Charger/valider la carte avec ton package carte si besoin
	grille , _ := carte.Parse(*cheminCarte)
	if *cheminCarte != "" {
    	_ = carte.ValidMap(grille) // ou carte.Valider(*cheminCarte)
	}
	// ex: _ = carte.Charger(*cheminCarte)

	// 3. Initialisation du scanner unique sur l'entrée standard
	scanner := bufio.NewScanner(os.Stdin)

	// 4. Boucle principale du jeu
	tour := 1
	journal.SupLogs()
	for {
		
		// package protocol lit l'arbitre jusqu'à "FIN"
		donnees := protocol.Extraction(scanner)

		//  package robot (fichier ia.go) prend la décision
		action := robot.Decider(donnees)

		//  Envoi de la commande à l'arbitre sur os.Stdout
		fmt.Println(action)

		journal.EcrireLog(tour,action)
		tour++
	}
}
