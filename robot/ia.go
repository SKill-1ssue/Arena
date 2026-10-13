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

	// Nettoyage de la case si on vient de ramasser un objet
	if len(grille) > d.Y1 && len(grille[d.Y1]) > d.X1 {
		c := grille[d.Y1][d.X1]
		if c == '$' || c == '+' || c == '*' || c == '!' {
			grille[d.Y1][d.X1] = '.'
		}
	}

	// 1. Priorité Combat : Tirer si l'ennemi est en ligne de vue (<= 5 cases)
	if tir := peutTirer(d.X1, d.Y1, d.X2, d.Y2); tir != "" {
		return tir
	}

	// 2. Priorité Objets : Chercher le bonus/soin/trésor/munitions le plus proche
	if posBonus, trouve := trouverBonusPlusProche(monPos); trouve {
		if ordre := ChercherPoint(monPos, posBonus); ordre != "" {
			return ordre
		}
	}

	// 3. Priorité Traque : Chasser l'ennemi en temps réel
	if ordre := ChercherPoint(monPos, ennemiPos); ordre != "" {
		return ordre
	}

	return "ATTENDS"
}

// ChercherPoint calcule le chemin optimal en attribuant un coût à chaque type de case
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

				// Seuls les murs sont totalement infranchissables
				if caseVoisine == '#' {
					continue
				}

				// PONDÉRATION DES TILES :
				// - Sol (.), Départs (1, 2), Objets (+, $, *, !) = 1 tour
				// - Eau (~) = 2 tours (déplacement + 1 tour gelé)
				// - Piège (^) = 10 tours (le robot contourne sauf si c'est la seule voie)
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

// Recherche le bonus ('$', '+', '*', '!') le plus proche sur la carte
func trouverBonusPlusProche(depart Point) (Point, bool) {
	var meilleurPoint Point
	minDist := 999999

	for y := 0; y < len(grille); y++ {
		for x := 0; x < len(grille[y]); x++ {
			c := grille[y][x]
			if c == '$' || c == '+' || c == '*' || c == '!' {
				dist := abs(x-depart.X) + abs(y-depart.Y)
				if dist < minDist {
					minDist = dist
					meilleurPoint = Point{X: x, Y: y}
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

func peutTirer(x1, y1, x2, y2 int) string {
	if y1 == y2 && x1 < x2 && x2-x1 <= 5 {
		for i := x1 + 1; i < x2; i++ {
			if grille[y1][i] == '#' {
				return ""
			}
		}
		return "TIRE E"
	}
	if y1 == y2 && x1 > x2 && x1-x2 <= 5 {
		for i := x2 + 1; i < x1; i++ {
			if grille[y1][i] == '#' {
				return ""
			}
		}
		return "TIRE O"
	}
	if x1 == x2 && y1 < y2 && y2-y1 <= 5 {
		for i := y1 + 1; i < y2; i++ {
			if grille[i][x1] == '#' {
				return ""
			}
		}
		return "TIRE S"
	}
	if x1 == x2 && y1 > y2 && y1-y2 <= 5 {
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
