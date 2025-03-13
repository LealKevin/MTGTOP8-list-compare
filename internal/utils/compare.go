package utils

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Pair struct {
	Key   string
	Value int
}

type Card struct {
	Name     string
	Quantity int
}

type List struct {
	Line map[string]int
}

type Player struct {
	Name string
}

func ScrapList(url string) (List, Player) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatalf("Erreur lors de la récupération de l'URL : %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("Erreur : status code %d", res.StatusCode)
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("Erreur lors du parsing HTML : %v", err)
	}
	cards := make(map[string]int)
	doc.Find("div.deck_line").Each(func(i int, s *goquery.Selection) {
		fullText := s.Text()
		fields := strings.Fields(fullText)
		if len(fields) < 1 {
			return
		}
		quantity, err := strconv.Atoi(fields[0])
		if err != nil {

			log.Printf("Unable to convert: %v", err)
			return
		}
		cardName := s.Find("span.L14").Text()
		cardName = strings.TrimSpace(cardName)
		if cardName != "" && quantity != 0 {
			cards[cardName] = quantity
		}
	})

	playerSelection := doc.Find("a.player_big")
	player := Player{Name: playerSelection.Text()}
	log.Printf("Player name: %s", player)
	return List{Line: cards}, player
}

func Compare(list1, list2 List) [][]Pair {
	dif1 := make(map[string]int)
	dif2 := make(map[string]int)
	same := make(map[string]int)

	for name, quantity1 := range list1.Line {
		if quantity2, exists := list2.Line[name]; exists {
			if quantity1 == quantity2 {
				same[name] = quantity1
			} else if quantity1 > quantity2 {
				same[name] = quantity2
				dif1[name] = quantity1 - quantity2
			} else {
				same[name] = quantity2
				dif2[name] = quantity2 - quantity1
			}
		} else {
			dif1[name] = quantity1
		}
	}

	for name, quantity2 := range list2.Line {
		if _, exists := list1.Line[name]; !exists {
			dif2[name] = quantity2
		}
	}

	dif1sorted := SortAlphabetically(List{Line: dif1})
	dif2sorted := SortAlphabetically(List{Line: dif2})
	samesorted := SortAlphabetically(List{Line: same})

	return [][]Pair{dif1sorted, dif2sorted, samesorted}
}

func SortAlphabetically(list List) []Pair {
	pairs := make([]Pair, 0, len(list.Line))
	for key, value := range list.Line {
		pairs = append(pairs, Pair{Key: key, Value: value})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Key < pairs[j].Key
	})
	return pairs
}
