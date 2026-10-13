package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	// Remplace "Arena" par le nom exact inscrit dans ton fichier go.mod si ce n'est pas Arena
	"github.com/votre-pseudo/arena/carte"
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
	for {
		// Étape A : Ton package protocol lit l'arbitre jusqu'à "FIN"
		donnees := protocol.Extraction(scanner)

		// Étape B : Ton package robot (fichier ia.go) prend la décision
		// Assure-toi que dans robot/ia.go tu as "package robot" et une fonction "Decider"
		action := robot.Decider(donnees)

		// Étape C : Envoi de la commande à l'arbitre sur os.Stdout
		fmt.Println(action)
	}
}