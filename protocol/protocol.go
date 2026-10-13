package protocol

import (
	"bufio"
	"strconv"
	"strings"
)

type Donnees struct {
	Tour	int
	Vie1 	int
	X1 		int
	Y1 		int
	Vie2 	int
	X2 		int
	Y2 		int
}


func Extraction(scanner *bufio.Scanner) Donnees {
	var d Donnees
	
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
				d.X1, _ = strconv.Atoi(mots[1])
				d.Y1, _ = strconv.Atoi(mots[2])
				d.Vie1, _ = strconv.Atoi(mots[4])
			}
		case "ENNEMI":
			if len(mots) >= 5 {
				d.X2, _ = strconv.Atoi(mots[1])
				d.Y2, _ = strconv.Atoi(mots[2])
				d.Vie2, _ = strconv.Atoi(mots[4])
			}
		}
	}
	return d	
}
