package server

import (
	"fmt"
	"net/http"

	"github.com/LealKevin/list-compare/internal/utils"
	"github.com/go-chi/chi/v5"
)

func CompareHandler(w http.ResponseWriter, r *http.Request) {

	list1JSON := r.FormValue("list1")
	list2JSON := r.FormValue("list2")

	list1Scrapped, player1 := utils.ScrapList(list1JSON)
	list2Scrapped, player2 := utils.ScrapList(list2JSON)

	result := utils.Compare(list1Scrapped, list2Scrapped)

	htmlResponse := `<div class="overflow-x-auto flex items-start justify-center w-full">`

	labels := []string{fmt.Sprintf("Difference in %s", player1.Name), fmt.Sprintf("Difference in %s", player2.Name), "Common Elements"}

	for i, resultList := range result {
		htmlResponse += "<table class='table table-s m-4 p-2 table-zebra border-separate border border-solid border-primary border-base-content/8 bg-base-100'><tbody>"
		htmlResponse += "<thead> <tr> <th></th> <th colspan='2' class='text-accent'>" + labels[i] + "</th> </tr> </thead>"
		for _, pair := range resultList {
			htmlResponse += "<tr><th>" + fmt.Sprintf("%d", pair.Value) + "</th><td>" + pair.Key + "</td></tr>"
		}
		htmlResponse += "</tbody></table>"
	}
	htmlResponse += ` </div>`

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlResponse))

}

func HomePage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}
func Router() *chi.Mux {
	mux := chi.NewRouter()

	mux.Get("/", HomePage)
	mux.Post("/compare", CompareHandler)

	return mux
}
