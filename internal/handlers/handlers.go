package handlers

import (
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/LealKevin/list-compare/internal/utils"
	"github.com/go-chi/chi/v5"
)

type TableData struct {
	Label string
	Pairs []utils.Pair
}

type PageData struct {
	Tables []TableData
	Error  string
}

func CompareHandler(w http.ResponseWriter, r *http.Request) {
	// Create the template
	tableTemplate := template.Must(template.ParseFiles("static/template/table.html"))

	// Get the list1 and list2 from the form
	list1JSON := r.FormValue("list1")
	list2JSON := r.FormValue("list2")

	// Scrap the lists
	list1Scrapped, player1, err1 := utils.ScrapList(list1JSON)

	list2Scrapped, player2, err2 := utils.ScrapList(list2JSON)
	if err1 != nil || err2 != nil {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		data := PageData{Error: "Something went wrong, be sure is a valid URL from MTGTop8"}
		tableTemplate.Execute(w, data)
		return
	}

	// Compare the lists
	result := utils.Compare(list1Scrapped, list2Scrapped)

	// Create the data to pass to the template
	tableData := PageData{
		Tables: []TableData{
			{Label: fmt.Sprintf("Difference in %s", player1.Name), Pairs: result[0]},
			{Label: fmt.Sprintf("Difference in %s", player2.Name), Pairs: result[1]},
			{Label: "Common Elements", Pairs: result[2]},
		},
	}
	// Execute the template
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	tableTemplate.Execute(w, tableData)
}

func GetImage(w http.ResponseWriter, r *http.Request) {
	imageName := chi.URLParam(r, "name")
	fmt.Println("🔥 Requête reçue pour:", imageName)

	imageURL := "https://api.scryfall.com/cards/named?fuzzy=" + imageName
	resp, err := http.Get(imageURL)
	if err != nil {
		http.Error(w, "Error downloading image", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != 200 {
		http.Error(w, "Error downloading image", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error downloading image", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Json: %s", string(body))
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}
