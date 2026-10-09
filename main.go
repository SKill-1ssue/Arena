package main

import (
	"fmt"

	"github.com/votre-pseudo/arena/carte"
	// Remplace par le bon chemin d'import de ton package
)

func main() {
	// 1. On récupère la grille et l'erreur
	grille, err := carte.Parse(`C:\Users\USER\Arena\carte\iles.map`)
	
	// 2. On vérifie s'il y a une erreur
	if err != nil {
		fmt.Println("Erreur de lecture :", err)
		return
	}

	// 3. On affiche chaque ligne sous forme de texte
	for _, ligne := range grille {
		fmt.Println(string(ligne))
	}

	print(carte.ValidMap(grille))
}