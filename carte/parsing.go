package carte

import (
	"bufio"
	"os"
)

func Parse(fichier string) ([][]rune,error) {
    file, err := os.Open(fichier)
	if err != nil {
		return nil,err
	}
	defer file.Close()

	var grille [][]rune
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		ligneTexte := scanner.Text()

		grille = append(grille, []rune(ligneTexte))

	}

	if err := scanner.Err(); err != nil {

	}

	return grille, nil

}


