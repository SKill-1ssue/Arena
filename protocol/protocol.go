package protocol

import (
	"bufio"
	"strconv"
	"strings"
)

type Objet struct {
	X, Y int
	Type string // "SOIN", "MUNITION", "BOUCLIER", "TRESOR"
}

type Point struct {
	X, Y int
}

type Donnees struct {
	Tour      int
	Vie1      int
	X1        int
	Y1        int
	Mun1      int
	Bouclier1 int

	Vie2      int
	X2        int
	Y2        int
	Mun2      int
	Bouclier2 int

	Objets    []Objet
	Position1 []Point
	Position2 []Point
}

func Extraction(scanner *bufio.Scanner) Donnees {
	var d Donnees
	d.Objets = make([]Objet, 0)

	for scanner.Scan() {
		ligne := strings.TrimSpace(scanner.Text())
		if ligne == "FIN" {
			break
		}

		mots := strings.Fields(ligne)
		if len(mots) == 0 {
			continue
		}

		switch mots[0] {
		case "TOUR":
			if len(mots) >= 2 {
				d.Tour, _ = strconv.Atoi(mots[1])
			}

		case "MOI":
			if len(mots) >= 5 {
				// Lecture de X, Y et VIE (valable en v1 et v2)
				d.X1, _ = strconv.Atoi(mots[1])
				d.Y1, _ = strconv.Atoi(mots[2])
				d.Vie1, _ = strconv.Atoi(mots[4])

				d.Position2 = append(d.Position2, Point{X:d.X1, Y:d.Y1})

				// En v2 : MOI x y VIE v MUN m BOUCLIER b (len = 9)
				if len(mots) >= 9 {
					d.Mun1, _ = strconv.Atoi(mots[6])
					d.Bouclier1, _ = strconv.Atoi(mots[8])
				}
			}

		case "ENNEMI":
			if len(mots) >= 5 {
				d.X2, _ = strconv.Atoi(mots[1])
				d.Y2, _ = strconv.Atoi(mots[2])
				d.Vie2, _ = strconv.Atoi(mots[4])
				
				d.Position2 = append(d.Position2, Point{X:d.X2, Y:d.Y2})
			}

		case "OBJET":
			// Format v2 : OBJET x y TYPE
			if len(mots) >= 4 {
				x, _ := strconv.Atoi(mots[1])
				y, _ := strconv.Atoi(mots[2])
				d.Objets = append(d.Objets, Objet{
					X:    x,
					Y:    y,
					Type: mots[3],
				})
			}
		}
	}
	return d
}
