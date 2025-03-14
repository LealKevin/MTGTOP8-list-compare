package handlers

import (
	"fmt"
	"github.com/LealKevin/list-compare/internal/utils"
	"html/template"
	"net/http"
)

type TableData struct {
	Label string
	Pairs []utils.Pair
}

type PageData struct {
	Tables []TableData
}

func CompareHandler(w http.ResponseWriter, r *http.Request) {

	// Get the list1 and list2 from the form
	list1JSON := r.FormValue("list1")
	list2JSON := r.FormValue("list2")

	// Scrap the lists
	list1Scrapped, player1 := utils.ScrapList(list1JSON)
	list2Scrapped, player2 := utils.ScrapList(list2JSON)

	// Compare the lists
	result := utils.Compare(list1Scrapped, list2Scrapped)

	// Create the template
	tableTemplate := template.Must(template.ParseFiles("static/template/table.html"))

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

func HomePageHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}
