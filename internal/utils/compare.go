package utils

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
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

func ScrapList(url string) (List, Player, error) {

	res, err := http.Get(url)
	if err != nil {
		return List{}, Player{}, fmt.Errorf("Unable to get URL : %v", err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return List{}, Player{}, fmt.Errorf("Expected status code 200, but got %d", res.StatusCode)
	}

	reader, err := charset.NewReader(res.Body, res.Header.Get("Content-Type"))
	if err != nil {
		return List{}, Player{}, errors.New("Unable to get charset reader")
	}

	log.Printf("Content Type: %s", res.Header.Get("Content-Type"))

	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return List{}, Player{}, fmt.Errorf("Error parsing HTML: %v", err)
	}
	cards := make(map[string]int)

	// Find the div with the class deck_line
	// For each div found, get the Text
	// Split the text into Fields
	// Convert the first field to an integer
	// Get the card name
	// Add the card name and quantity to the cards map

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

	// Find the player name
	playerSelection := doc.Find("a.player_big")
	player := Player{Name: playerSelection.Text()}
	//log.Printf("List 1: %v", cards)
	//log.Printf("Player name: %s", player)
	return List{Line: cards}, player, nil
}

func Compare(list1, list2 List) [][]Pair {
	dif1 := make(map[string]int)
	dif2 := make(map[string]int)
	same := make(map[string]int)

	// For each card in list1, check if it exists in list2
	// If it exists, check if the quantity is the same
	// If the quantity is the same, add it to the same map
	// If the quantity is different, add it to the dif1 or dif2 map
	// If the card does not exist in list2, add it to the dif1 map
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

	// For each card in list2, check if it exists in list1
	// If it does not exist, add it to the dif2 map
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
