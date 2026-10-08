package carte

import (
	"fmt"
	"slices"
)

func ContientSeulement(grille [][]rune, autorises []rune) bool {
    for _, ligne := range grille {
        for _, char := range ligne {
            // Si une case n'est PAS dans la liste des autorisés
            if !slices.Contains(autorises, char) {
                return false 
            }
        }
    }
    return true
}

func ValidMap(grille [][]rune) string {
	if len(grille[0]) < 4 || len(grille[1]) < 12 || len(grille[2]) < 7 || len(grille[3]) < 10 {
		return "map invalide"
	}

	//verification nom
	if string(grille[0][0:4]) != "NOM " {
		return "map invalide"
	}

	//verification difficulte
	if string(grille[1][0:11]) != "DIFFICULTE " {
		return "map invalide"
	}
	for i := 11; i < len(grille[1]); i++ {
		if grille[1][i] < '1' || grille[1][i] > '3' {
			return "map invalide"
		}
	}

	//verification tours
	if string(grille[2][0:6]) != "TOURS " {
		return "map invalide"
	}
	for i := 6; i < len(grille[2])-1; i++ {
		if grille[2][i] < '0' || grille[2][i] > '9' {
			return "map invalide"
		}
	}

	//verification taille et prise des dimentions de la map dans TAILLE
	if string(grille[3][0:7]) != "TAILLE " {
		return "map invalide"
	}

	longueur := 0
	largeur := 0
	flag := 0

	for i := 7; i < len(grille[3]); i++ {
		if grille[3][i] == ' ' {
			if grille[3][i-1] == ' ' {
				return "map invalide"
			}
			flag = 1
			continue
		}

		if grille[3][i] < '0' || grille[3][i] > '9' {
			return "map invalide"
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
		return "map invalide"
	}
	for i := start; i < len(grille); i++ {

		if len(grille[i]) != longueur {
			return "map invalide"
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

	if !ContientSeulement(grilleSymbole,symbolesValides) {
		return "map invalide"
	}

	//verification que les bordures de la map soient toujours des #
	bordure := ""
	for i := 0; i < longueur; i++ {
		bordure += "#"
	}
	if string(grille[start]) != bordure {
		return "map invalide"
	} else {
		for i := start+1; i < len(grille); i++ {
			if grille[i][0] != '#' && grille[i][longueur-1] != '#' {
				return "map invalide"
			}
			if i == len(grille) && string(grille[i]) != bordure {
				return "map invalide"
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

	if nbrPerso1 != 1 || nbrPerso2 != 1 {
		return "map inavalide"
	}




	print(longueur)
	fmt.Println()

	print(largeur)
	return "\nmap valide"
}
