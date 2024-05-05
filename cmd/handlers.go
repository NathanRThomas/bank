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
	"strconv"
	"time"
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

func (this *app) updateLedger(w http.ResponseWriter, r *http.Request) {

	ledgerSlug := chi.URLParam(r, "ledger")
	if len(ledgerSlug) == 0 {
		http.Error(w, "Ledger is invalid", http.StatusNotFound)
		return 
	}

	// find it from file
	fileLoc := opts.Ledgers + "/" + ledgerSlug + ".json"
	jsFile, err := os.ReadFile(fileLoc)
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

	// get the url param for how much we're paying
	// convert it to an int
	cents, _ := strconv.Atoi(r.URL.Query().Get("cents"))

	if cents == 0 {
		http.Error(w, "cents was zero", http.StatusBadRequest)
		return
	}

	lineItem := r.URL.Query().Get("line")

	if cents < 0 {
		lineItem = "Payment"
	} else if len(lineItem) == 0 {
		http.Error(w, "line was empty", http.StatusBadRequest)
		return
	}

	// add this line item
	ledger.LineItems = append (ledger.LineItems, models.LineItem {
		Date: time.Now(),
		Note: lineItem,
		Amount: cents,
	})

	// now save it back to the file

	jstr, err := json.Marshal(ledger)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Bad json", http.StatusConflict)
		return 
	}

	err = os.WriteFile (fileLoc, jstr, 0666)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "saving file", http.StatusConflict)
		return 
	}
}
