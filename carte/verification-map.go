package carte

import "fmt"

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
		if grille[1][i] < '0' || grille[1][i] > '9' {
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
	if len(grille[4:]) != largeur {
		return "map invalide"
	}
	for i := 4; i < len(grille); i++ {
		if len(grille[i]) != longueur {
			return "map invalide"
		}
	}

	print(longueur)
	fmt.Println()
	print(largeur)
	return "\nmap valide"
}