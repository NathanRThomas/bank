/** ****************************************************************************************************************** **
	General handlers

** ****************************************************************************************************************** **/

package main

import (
	"bank/models"
	"github.com/go-chi/chi/v5"
	
	"html/template"
	"net/http"
	"log"
	"os"
	"encoding/json"
)

  //-------------------------------------------------------------------------------------------------------------------------//
 //----- CONST -------------------------------------------------------------------------------------------------------------//
//-------------------------------------------------------------------------------------------------------------------------//

  //-------------------------------------------------------------------------------------------------------------------------//
 //----- HANDLERS ----------------------------------------------------------------------------------------------------------//
//-------------------------------------------------------------------------------------------------------------------------//

func (this *app) getLedger(w http.ResponseWriter, r *http.Request) {

	ledgerSlug := chi.URLParam(r, "ledger")
	if len(ledgerSlug) == 0 {
		http.Error(w, "Ledger is invalid", http.StatusNotFound)
		return 
	}

	// find it from file
	jsFile, err := os.ReadFile(opts.Ledgers + "/" + ledgerSlug + ".json")
	if err != nil {
		http.Error(w, "Ledger not found", http.StatusNotFound)
		return 
	}

	ledger := &models.Ledger{} // this is our object

	err = json.Unmarshal(jsFile, ledger)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Ledger is corrupt", http.StatusConflict)
		return 
	}

	tmpl := template.Must(template.ParseFiles(opts.Templates + "/ledger.html"))

	ledger.Payoff() // calculate things
	ledger.Name = ledgerSlug

	err = tmpl.Execute(w, ledger)
	if err != nil {
		log.Println("template error", err)
	}
}
