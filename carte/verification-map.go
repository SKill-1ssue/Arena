package carte

func ValidMap(grille [][]rune) string {

	if len(grille[0]) < 4 || len(grille[1]) < 12 || len(grille[2]) < 7 || len(grille[3]) < 10 {
    	return "map invalide"
	}

	//verification nom
	if string(grille[0][0:4]) != "NOM "  {
		return "map invalide"
	}

	//verification difficulte
	if string(grille[1][0:11]) != "DIFFICULTE "  {
		return "map invalide"
	}
	for i := 11; i < len(grille[1]); i++ {
		if grille[1][i] < '0' || grille[1][i] > '9' {
			return "map invalide"
		}
	}
	
	//verification tours
	if string(grille[2][0:6]) != "TOURS "  {
		return "map invalide"
	}
	for i := 6; i < len(grille[2])-1; i++ {
		if grille[2][i] < '0' || grille[2][i] > '9' {
			return "map invalide"
		}
	}

	//verification taille
	if string(grille[3][0:7]) != "TAILLE "  {
		return "map invalide"
	}
	/*
	for i := 6; i < len(grille[2])-1; i++ {
		if i == ' ' {
			continue
		}
		if grille[1][i] < '0' || grille[1][i] > '9' {
			return "map invalide"
		}
		 

	} */

	return "map valide"
}