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
	nbEsquives   int = 0
	munitions2   int = 10
	posAEviter   Point = Point{X: -1, Y: -1} // Case évitee temporairement après une esquive

	// Système anti-collision et detection d'immobilite
	derniereCaseVisee Point = Point{X: -1, Y: -1}
	caseBloquee       Point = Point{X: -1, Y: -1}
	nbCollisions      int   = 0
)
func Decider(d protocol.Donnees) string {
	if !carteChargee {
		chargerCarteDepuisFlag()
		carteChargee = true
	}

	monPos := Point{X: d.X1, Y: d.Y1}
	ennemiPos := Point{X: d.X2, Y: d.Y2}

	// 0. MISE À JOUR DES ÉTATS ET DÉTECTION DU BLOCAGE
	munitionBot2(d)
	actualiserObjetsSurGrille(d)

	if len(d.Position1) >= 2 {
		posPrc := Point{X: d.Position1[len(d.Position1)-2].X, Y: d.Position1[len(d.Position1)-2].Y}
		if monPos == posPrc && derniereCaseVisee != (Point{X: -1, Y: -1}) {
			nbCollisions++
			if nbCollisions >= 2 {
				caseBloquee = derniereCaseVisee
			}
		} else {
			nbCollisions = 0
			caseBloquee = Point{X: -1, Y: -1}
		}
	}

	if abs(monPos.X-posAEviter.X)+abs(monPos.Y-posAEviter.Y) > 1 {
		posAEviter = Point{X: -1, Y: -1}
	}

	// 1. URGENCE SOIN : PV <= 40
	if d.Vie1 <= 40 {
		if posSoin, trouve := trouverRunePlusProche(monPos, '+'); trouve {
			if ordre := ChercherPoint(monPos, posSoin, d); ordre != "" {
				return memoriserEtRenvoyer(monPos, ordre)
			}
		}
	}

	// 2 URGENCE MUNITIONS : 0 balle restante
	if d.Mun1 == 0 {
		if posMun, trouve := trouverRunePlusProche(monPos, '*'); trouve {
			if ordre := ChercherPoint(monPos, posMun, d); ordre != "" {
				return memoriserEtRenvoyer(monPos, ordre)
			}
		}
	}

	// 3 COMBAT ET ESQUIVE
	if tir := peutTirer(d.X1, d.Y1, d.X2, d.Y2); tir != "" {
		if munitions2 > 0 && nbEsquives < 5 {
			if esq := esquive(d); esq != "" {
				nbEsquives++
				posAEviter = monPos
				return memoriserEtRenvoyer(monPos, esq)
			}
		}

		if d.Mun1 > 0 {
			posAEviter = Point{X: -1, Y: -1}
			return memoriserEtRenvoyer(monPos, tir)
		}
	}

	// 3.5. TRAQUE opport: Ennemi sans munition et sans recharge accessible
	if munitions2 == 0 {
		posMun, trouve := trouverRunePlusProche(ennemiPos, '*')
		if !trouve {
			if ordre := ChercherPoint(monPos, ennemiPos, d); ordre != "" {
				return memoriserEtRenvoyer(monPos, ordre)
			}
		} else {
			distMoi := abs(monPos.X-posMun.X) + abs(monPos.Y-posMun.Y)
			distEnnemi := abs(ennemiPos.X-posMun.X) + abs(ennemiPos.Y-posMun.Y)

			if distMoi <= distEnnemi {
				if ordre := ChercherPoint(monPos, ennemiPos, d); ordre != "" {
					return memoriserEtRenvoyer(monPos, ordre)
				}
			}
		}
	}

	// 4. PRIORITÉ SCORE : Trésors '$' (+50 pts)
	if posTresor, trouve := trouverRunePlusProche(monPos, '$'); trouve {
		if ordre := ChercherPoint(monPos, posTresor, d); ordre != "" {
			return memoriserEtRenvoyer(monPos, ordre)
		}
	}

	// 5. AUTRES OBJETS
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
			if ordre := ChercherPoint(monPos, posBonus, d); ordre != "" {
				return memoriserEtRenvoyer(monPos, ordre)
			}
		}
	}

	// 6. TRAQUE PAR DÉFAUT
	if ordre := ChercherPoint(monPos, ennemiPos, d); ordre != "" {
		return memoriserEtRenvoyer(monPos, ordre)
	}

	return memoriserEtRenvoyer(monPos, "ATTENDS")
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

// ChercherPoint calcule le premier mouvement vers la cible 
func ChercherPoint(start, goal Point, d protocol.Donnees) string {
	if start == goal {
		return ""
	}

	h := len(grille)
	if h == 0 {
		return ""
	}
	w := len(grille[0])

	dist := make([][]int, h)
	for y := 0; y < h; y++ {
		dist[y] = make([]int, w)
		for x := 0; x < w; x++ {
			dist[y][x] = 1e9
		}
	}

	queue := []Elem{
		{P: start, PremierMouv: "", Cout: 0},
	}
	dist[start.Y][start.X] = 0

	dirs := []struct {
		dx, dy int
		ordre  string
	}{
		{0, -1, "AVANCE N"},
		{0, 1, "AVANCE S"},
		{1, 0, "AVANCE E"},
		{-1, 0, "AVANCE O"},
	}

	ennemiPos := Point{X: d.X2, Y: d.Y2}

	for len(queue) > 0 {
		bestIdx := 0
		for i := 1; i < len(queue); i++ {
			if queue[i].Cout < queue[bestIdx].Cout {
				bestIdx = i
			}
		}
		curr := queue[bestIdx]
		queue = append(queue[:bestIdx], queue[bestIdx+1:]...)

		if curr.P == goal {
			return curr.PremierMouv
		}

		if curr.Cout > dist[curr.P.Y][curr.P.X] {
			continue
		}

		for _, dStruct := range dirs {
			nx, ny := curr.P.X+dStruct.dx, curr.P.Y+dStruct.dy

			// 1. Limites de la carte
			if ny < 0 || ny >= h || nx < 0 || nx >= w {
				continue
			}

			voisin := Point{X: nx, Y: ny}

			// 2. Interdiction de la case d'esquive récente
			if voisin == posAEviter {
				continue
			}

			c := grille[ny][nx]

			// 3. Obstacles (murs et pièges)
			if c == '#' || c == '^' {
				continue
			}

			// 4. Anti-collision : ne pas foncer sur l'ennemi sauf si on le traque
			if voisin == ennemiPos && goal != ennemiPos {
				continue
			}

			// 5. Calcul du coût de base du terrain
			coutBase := 1
			if c == '~' {
				coutBase = 3 // Ralentissement dû à l'eau
			}

			// 6. Pénalité sur case verrouillée par collisions répétées (+200)
			if voisin == caseBloquee {
				coutBase += 200
			}

			// 7. Anti-ping-pong / anti-boucle basés sur l'historique
			if len(d.Position1) >= 2 {
				pos1Tour := Point{X: d.Position1[len(d.Position1)-2].X, Y: d.Position1[len(d.Position1)-2].Y}
				if voisin == pos1Tour {
					coutBase += 15 // Pénalité demi-tour immédiat
				}
			}
			if len(d.Position1) >= 4 {
				pos3Tours := Point{X: d.Position1[len(d.Position1)-4].X, Y: d.Position1[len(d.Position1)-4].Y}
				if voisin == pos3Tours {
					coutBase += 30 // Pénalité boucle sur 3 tours (A -> B -> C -> A)
				}
			}

			// 8. Surcoût massif (+100) pour contourner la ligne de tir ennemie
			if munitions2 > 0 && estDansLigneDeTir(voisin, ennemiPos) {
				coutBase += 100
			}

			nouveauCout := curr.Cout + coutBase

			if nouveauCout < dist[ny][nx] {
				dist[ny][nx] = nouveauCout

				mouv := curr.PremierMouv
				if mouv == "" {
					mouv = dStruct.ordre
				}

				queue = append(queue, Elem{
					P:           voisin,
					PremierMouv: mouv,
					Cout:        nouveauCout,
				})
			}
		}
	}

	return "" // Aucun chemin trouvé
}

// Vérifie si une case est sur l'axe de tir direct de l'ennemi à portée (<= 5 cases, sans mur)
func estDansLigneDeTir(caseTest, ennemi Point) bool {
	if caseTest.X != ennemi.X && caseTest.Y != ennemi.Y {
		return false
	}

	// Même colonne
	if caseTest.X == ennemi.X {
		minY, maxY := caseTest.Y, ennemi.Y
		if minY > maxY {
			minY, maxY = maxY, minY
		}
		if maxY-minY > 5 {
			return false
		}
		for y := minY + 1; y < maxY; y++ {
			if grille[y][caseTest.X] == '#' {
				return false
			}
		}
		return true
	}

	// Même ligne
	if caseTest.Y == ennemi.Y {
		minX, maxX := caseTest.X, ennemi.X
		if minX > maxX {
			minX, maxX = maxX, minX
		}
		if maxX-minX > 5 {
			return false
		}
		for x := minX + 1; x < maxX; x++ {
			if grille[caseTest.Y][x] == '#' {
				return false
			}
		}
		return true
	}

	return false
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

func esquive(d protocol.Donnees) string {
	menace := peutTirer(d.X1, d.Y1, d.X2, d.Y2)

	switch menace {
	case "TIRE E", "TIRE O":
		// Vraie esquive perpendiculaire (sort du champ de tir)
		if estCaseSûre(d.X1, d.Y1-1) {
			return "AVANCE N"
		}
		if estCaseSûre(d.X1, d.Y1+1) {
			return "AVANCE S"
		}
		// Secours uniquement si on a ZERO munition (reculer pour gagner du temps)
		if d.Mun1 == 0 {
			if menace == "TIRE E" && estCaseSûre(d.X1-1, d.Y1) {
				return "AVANCE O"
			}
			if menace == "TIRE O" && estCaseSûre(d.X1+1, d.Y1) {
				return "AVANCE E"
			}
		}

	case "TIRE S", "TIRE N":
		// Vraie esquive perpendiculaire (sort du champ de tir)
		if estCaseSûre(d.X1+1, d.Y1) {
			return "AVANCE E"
		}
		if estCaseSûre(d.X1-1, d.Y1) {
			return "AVANCE O"
		}
		// Secours uniquement si on a ZERO munition
		if d.Mun1 == 0 {
			if menace == "TIRE S" && estCaseSûre(d.X1, d.Y1-1) {
				return "AVANCE N"
			}
			if menace == "TIRE N" && estCaseSûre(d.X1, d.Y1+1) {
				return "AVANCE S"
			}
		}
	}

	return "" // Impossible de sortir du champ de tir -> Decider va déclencher le tir
}

func estCaseSûre(x, y int) bool {
	if y < 0 || y >= len(grille) || x < 0 || x >= len(grille[y]) {
		return false
	}
	c := grille[y][x]
	return c != '#' && c != '^'
}

func munitionBot2(d protocol.Donnees) {
	if len(d.Position2) < 2 {
		return
	}

	posPrc := d.Position2[len(d.Position2)-2]
	posAct := d.Position2[len(d.Position2)-1]

	// 1. DÉTECTION DU TIR : L'ennemi reste immobile ET a une ligne de vue
	if posPrc.X == posAct.X && posPrc.Y == posAct.Y && peutTirer(d.X1, d.Y1, d.X2, d.Y2) != "" {
		if munitions2 > 0 {
			munitions2--
		}
	}

	// 2. RAMASSAGE DE MUNITIONS '*'
	if posAct.Y >= 0 && posAct.Y < len(grille) && posAct.X >= 0 && posAct.X < len(grille[posAct.Y]) {
		if grille[posAct.Y][posAct.X] == '*' {
			munitions2 += 5
		}
	}
}

func memoriserEtRenvoyer(monPos Point, ordre string) string {
	switch ordre {
	case "AVANCE N":
		derniereCaseVisee = Point{X: monPos.X, Y: monPos.Y - 1}
	case "AVANCE S":
		derniereCaseVisee = Point{X: monPos.X, Y: monPos.Y + 1}
	case "AVANCE E":
		derniereCaseVisee = Point{X: monPos.X + 1, Y: monPos.Y}
	case "AVANCE O":
		derniereCaseVisee = Point{X: monPos.X - 1, Y: monPos.Y}
	default:
		derniereCaseVisee = Point{X: -1, Y: -1}
	}
	return ordre
}