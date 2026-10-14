package robot

import (
	"bufio"
	"flag"
	"os"
	"strings"

	"github.com/votre-pseudo/arena/protocol"
)

type Point struct {
	X, Y int
}

type Elem struct {
	P           Point
	PremierMouv string
	Cout        int
}

var (
	grille       [][]rune
	carteChargee bool
)

func Decider(d protocol.Donnees) string {
	if !carteChargee {
		chargerCarteDepuisFlag()
		carteChargee = true
	}

	monPos := Point{X: d.X1, Y: d.Y1}
	ennemiPos := Point{X: d.X2, Y: d.Y2}

	// Synchronisation dynamique des objets en v2 sur la grille
	actualiserObjetsSurGrille(d)

	// 1. URGENCE SOIN : PV <= 40
	if d.Vie1 <= 40 {
		if posSoin, trouve := trouverRunePlusProche(monPos, '+'); trouve {
			if ordre := ChercherPoint(monPos, posSoin); ordre != "" {
				return ordre
			}
		}
	}

	// 2. URGENCE MUNITIONS : 0 balle restante
	if d.Mun1 == 0 {
		if posMun, trouve := trouverRunePlusProche(monPos, '*'); trouve {
			if ordre := ChercherPoint(monPos, posMun); ordre != "" {
				return ordre
			}
		}
	}

	// 3. COMBAT : Tirer uniquement si l'ennemi est à portée ET qu'il reste des munitions
	
		
	
	if d.Mun1 > 0 {
		if tir := peutTirer(d.X1, d.Y1, d.X2, d.Y2); tir != "" {
			return tir
		}
	}
	if tir := peutTirer(d.X1, d.Y1, d.X2, d.Y2); tir != "" {
			return tir
	}
	
	

	// 4. PRIORITÉ SCORE : Trésors '$' (+50 pts)
	if posTresor, trouve := trouverRunePlusProche(monPos, '$'); trouve {
		if ordre := ChercherPoint(monPos, posTresor); ordre != "" {
			return ordre
		}
	}

	// 5. AUTRES OBJETS : Boucliers '!', Munitions '*' (si < 10), Soins '+' (si < 100 PV)
	var objetsRecherches []rune
	objetsRecherches = append(objetsRecherches, '!')
	if d.Mun1 < 10 {
		objetsRecherches = append(objetsRecherches, '*')
	}
	if d.Vie1 < 80 {
		objetsRecherches = append(objetsRecherches, '+')
	}

	if len(objetsRecherches) > 0 {
		if posBonus, trouve := trouverRunePlusProche(monPos, objetsRecherches...); trouve {
			if ordre := ChercherPoint(monPos, posBonus); ordre != "" {
				return ordre
			}
		}
	}

	// 6. TRAQUE : Chasser l'adversaire
	if ordre := ChercherPoint(monPos, ennemiPos); ordre != "" {
		return ordre
	}

	return "ATTENDS"
}

// Nettoie les anciens objets et place ceux envoyés par l'arbitre au tour actuel
func actualiserObjetsSurGrille(d protocol.Donnees) {
	for y := 0; y < len(grille); y++ {
		for x := 0; x < len(grille[y]); x++ {
			c := grille[y][x]
			if c == '$' || c == '+' || c == '*' || c == '!' {
				grille[y][x] = '.'
			}
		}
	}

	for _, obj := range d.Objets {
		var r rune
		switch obj.Type {
		case "SOIN":
			r = '+'
		case "MUNITION":
			r = '*'
		case "BOUCLIER":
			r = '!'
		case "TRESOR":
			r = '$'
		}
		if r != 0 && obj.Y < len(grille) && obj.X < len(grille[obj.Y]) {
			grille[obj.Y][obj.X] = r
		}
	}
}

func trouverRunePlusProche(depart Point, cibleRunes ...rune) (Point, bool) {
	var meilleurPoint Point
	minDist := 999999

	for y := 0; y < len(grille); y++ {
		for x := 0; x < len(grille[y]); x++ {
			c := grille[y][x]
			for _, cible := range cibleRunes {
				if c == cible {
					dist := abs(x-depart.X) + abs(y-depart.Y)
					if dist < minDist {
						minDist = dist
						meilleurPoint = Point{X: x, Y: y}
					}
				}
			}
		}
	}
	return meilleurPoint, minDist < 999999
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func ChercherPoint(depart Point, destination Point) string {
	if depart == destination {
		return ""
	}

	const INF = 999999
	dist := make([][]int, len(grille))
	for i := range dist {
		dist[i] = make([]int, len(grille[0]))
		for j := range dist[i] {
			dist[i][j] = INF
		}
	}

	dist[depart.Y][depart.X] = 0
	queue := []Elem{{P: depart, PremierMouv: "", Cout: 0}}

	dirs := []struct {
		dx, dy int
		nom    string
	}{
		{0, -1, "AVANCE N"},
		{0, 1, "AVANCE S"},
		{1, 0, "AVANCE E"},
		{-1, 0, "AVANCE O"},
	}

	for len(queue) > 0 {
		minIdx := 0
		for i := 1; i < len(queue); i++ {
			if queue[i].Cout < queue[minIdx].Cout {
				minIdx = i
			}
		}
		actuel := queue[minIdx]
		queue = append(queue[:minIdx], queue[minIdx+1:]...)

		if actuel.P == destination {
			return actuel.PremierMouv
		}

		for _, d := range dirs {
			nx := actuel.P.X + d.dx
			ny := actuel.P.Y + d.dy

			if ny >= 0 && ny < len(grille) && nx >= 0 && nx < len(grille[0]) {
				caseVoisine := grille[ny][nx]

				if caseVoisine == '#' {
					continue
				}

				coutDeplacement := 1
				switch caseVoisine {
				case '~':
					coutDeplacement = 2
				case '^':
					coutDeplacement = 10
				}

				nouveauCout := actuel.Cout + coutDeplacement

				if nouveauCout < dist[ny][nx] {
					dist[ny][nx] = nouveauCout

					mouv := actuel.PremierMouv
					if mouv == "" {
						mouv = d.nom
					}

					queue = append(queue, Elem{
						P:           Point{X: nx, Y: ny},
						PremierMouv: mouv,
						Cout:        nouveauCout,
					})
				}
			}
		}
	}

	return ""
}

func peutTirer(x1, y1, x2, y2 int) string {
	if y1 == y2 && x1 < x2 && x2-x1 <= 6 {
		for i := x1 + 1; i < x2; i++ {
			if grille[y1][i] == '#' {
				return ""
			}
		}
		return "TIRE E"
	}
	if y1 == y2 && x1 > x2 && x1-x2 <= 6 {
		for i := x2 + 1; i < x1; i++ {
			if grille[y1][i] == '#' {
				return ""
			}
		}
		return "TIRE O"
	}
	if x1 == x2 && y1 < y2 && y2-y1 <= 6 {
		for i := y1 + 1; i < y2; i++ {
			if grille[i][x1] == '#' {
				return ""
			}
		}
		return "TIRE S"
	}
	if x1 == x2 && y1 > y2 && y1-y2 <= 6 {
		for i := y2 + 1; i < y1; i++ {
			if grille[i][x1] == '#' {
				return ""
			}
		}
		return "TIRE N"
	}
	return ""
}

func chargerCarteDepuisFlag() {
	f := flag.Lookup("carte")
	if f == nil || f.Value.String() == "" {
		return
	}

	fichier, err := os.Open(f.Value.String())
	if err != nil {
		return
	}
	defer fichier.Close()

	scanner := bufio.NewScanner(fichier)
	dansGrille := false

	for scanner.Scan() {
		ligne := strings.TrimRight(scanner.Text(), "\r\n")
		if !dansGrille && strings.HasPrefix(ligne, "#") {
			dansGrille = true
		}
		if dansGrille && len(ligne) > 0 {
			grille = append(grille, []rune(ligne))
		}
	}
}

//func predirePosition() Point {
	
//}