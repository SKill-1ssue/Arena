package carte

import (
	"fmt"
	"slices"
	"strconv"
)

func ContientSeulement(grille [][]rune, autorises []rune) (bool,int,int) {
    for i, ligne := range grille {
        for j, char := range ligne {
            // Si une case n'est PAS dans la liste des autorisés
            if !slices.Contains(autorises, char) {
                return false, i, j 
            }
        }
    }
    return true, 0,0
}

func ValidMap(grille [][]rune) string {
	if len(grille[0]) < 4 || len(grille[1]) < 12 || len(grille[2]) < 7 || len(grille[3]) < 10 {
		return "map invalide"
	}

	//verification nom
	if string(grille[0][0:4]) != "NOM " {
		return "<fichier> ligne 1 : clé inconnue " + string(grille[0][0:4])
	}

	//verification difficulte
	if string(grille[1][0:11]) != "DIFFICULTE " {
		return "<fichier> ligne 2 : clé inconnue " + string(grille[1][0:11])
	}
	for i := 11; i < len(grille[1]); i++ {
		if grille[1][i] < '1' || grille[1][i] > '3' {
			return "<fichier> ligne 2 : DIFFICULTE attend une valeur comprise entre 1 et 3"
		}
	}

	//verification tours
	if string(grille[2][0:6]) != "TOURS " {
		return "<fichier> ligne 3 : clé inconnue " + string(grille[2][0:6])
	}
	for i := 6; i < len(grille[2])-1; i++ {
		if grille[2][i] < '0' || grille[2][i] > '9' {
			return "<fichier> ligne 3 : TOUR attend une valeur numérique"
		}
	}

	//verification taille et prise des dimentions de la map dans TAILLE
	if string(grille[3][0:7]) != "TAILLE " {
		return "<fichier> ligne 4 : clé inconnue " + string(grille[2][0:6])
	}

	longueur := 0
	largeur := 0
	flag := 0

	for i := 7; i < len(grille[3]); i++ {
		if grille[3][i] == ' ' {
			if grille[3][i-1] == ' ' {
				return "<fichier> ligne 4 : TAILLE attend des valeur numérique"
			}
			flag = 1
			continue
		}

		if grille[3][i] < '0' || grille[3][i] > '9' {
			return "<fichier> ligne 4 : TAILLE attend des valeur numérique"
		}

		if flag == 0 {
			longueur = longueur*10 + int(grille[3][i] - '0')
		}

		if flag == 1 {
			largeur = largeur*10 + int(grille[3][i] - '0')
		}
	}

	// verification du fait que les dimentions de la grille correspondent aux dimentions de TAILLE
	start := 0
	for s := range grille {
		if grille[s][0] == '#' {
			start = s
			break
		}
	}
	if start == 0 {
		return "map invalide"
	} 
	if len(grille[start:]) != largeur {
		return "<fichier> : " + strconv.Itoa(len(grille[start:])) + " lignes de grille, " + strconv.Itoa(largeur) + " attendues"
	}
	for i := start; i < len(grille); i++ {

		if len(grille[i]) != longueur {
			return "<fichier> ligne " + strconv.Itoa(i+1) + " : " + strconv.Itoa(len(grille[i])) + " cases, " + strconv.Itoa(longueur)+ " attendues"
		}
	}

	//verification que la map ne contient que les symboles autorises
	var symbolesValides = []rune{
		'#', // Mur
		'.', // Sol
		'~', // Eau
		'^', // Piège
		'+', // Soin
		'*', // Munitions
		'!', // Bouclier
		'$', // Trésor
		'1', // Départ joueur 1
		'2', // Départ joueur 2
	}
	grilleSymbole := grille[start:]

	boolean, x, y := ContientSeulement(grilleSymbole,symbolesValides)

	if !boolean {
		return "<fichier> ligne " + strconv.Itoa(x+5) + " colonne " + strconv.Itoa(y+1) + " : symbole inconnu: " + string(grille[x+4][y])
	}

	//verification que les bordures de la map soient toujours des #
	bordure := ""
	for i := 0; i < longueur; i++ {
		bordure += "#"
	}
	if string(grille[start]) != bordure {
		for i := range grille[start] {
			if grille[start][i] != '#' {
				return "<fichier> ligne " + strconv.Itoa(start+1) + " colonne " + strconv.Itoa(i+1) + " : la bordure doit être un mur"
			}
		}
		
	} else {
		for i := start+1; i < len(grille); i++ {
			if grille[i][0] != '#' {
				return "<fichier> ligne " + strconv.Itoa(i+1) + " colonne " + strconv.Itoa(1) + " : la bordure doit être un mur"
			}
			if grille[i][longueur-1] != '#' {
				return "<fichier> ligne " + strconv.Itoa(i+1) + " colonne " + strconv.Itoa(longueur) + " : la bordure doit être un mur"
			}	
		}

		for i := range grille[len(grille)-1] {
			if grille[len(grille)-1][i] != '#' {
				return "<fichier> ligne " + strconv.Itoa(len(grille)) + " colonne " + strconv.Itoa(i+1) + " : la bordure doit être un mur"
			}
		}
		
	}

	//verification des spawn
	nbrPerso1 := 0
	nbrPerso2 := 0

	for i := start; i < len(grille); i++ {
		for j := range grille[i] {
			if grille[i][j] == '1' {
				nbrPerso1 += 1
			}
			if grille[i][j] == '2' {
				nbrPerso2 += 1
			}
		}
	}
	switch {
	case nbrPerso1 > 1:
		return "départ du joueur 1 déjà défini"

	case nbrPerso2 > 1:
		return "départ du joueur 2 déjà défini"

	case nbrPerso1 == 0:
		return "<fivhier> : aucun départ pour le joueur 1"

	case nbrPerso2 == 0:
		return "<fichier> : aucun départ pour le joueur 2"
	}

	print(longueur)
	fmt.Println()

	print(largeur)
	return "\nmap valide"
}
