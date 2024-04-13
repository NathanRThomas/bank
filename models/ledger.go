
package models 

import (
	
	"time"
	"math"
	"fmt"
)

type LineItem struct {
	Date time.Time 
	Note string 
	Amount int // in cents 
}

func (this *LineItem) TemplateDate () string {
	return this.Date.Format("Jan 2, 2006")
}

func (this *LineItem) TemplateAmount () string {
	return fmt.Sprintf("$%0.2f", float64(this.Amount) / 100)
}

type Ledger struct {
	Rake float64 // apr
	LineItems []LineItem 

	// for the template
	MinPayment, PayoffAmount, Name string 
}

// figures out the payout amount of the loan as it is right now
// returns the min payment as well as the total outstanding
func (this *Ledger) Payoff () (string, string) {
	if len(this.LineItems) == 0 { 
		this.MinPayment, this.PayoffAmount = "$0", "$0"
		return this.MinPayment, this.PayoffAmount // we're clean
	}

	var cnt, minPay float64

	// loop through the line items
	for idx, li := range this.LineItems {
		// add to our running total, only if positive
		if li.Amount > 0 {
			cnt += float64(li.Amount) / 100
		}

		if idx == 0 { continue } // nothing to compare it against

		// figure out how many days we're in we are
		hoursBack := li.Date.Sub(this.LineItems[idx-1].Date).Hours()

		// compare it to the last date we had
		final := cnt * math.Pow(math.E, this.Rake * ((hoursBack / 24) / 365))
		minPay += final - cnt 

		cnt = final // keep the running total accurate

		// now subtract the payment, if any
		if li.Amount < 0 {
			cnt += float64(li.Amount) / 100 
			minPay = 0 // reset the min
		}
	}

	// now do it to the current date
	hoursBack := time.Now().Sub(this.LineItems[len(this.LineItems) - 1].Date).Hours()

	// compare it to the last date we had
	// we want this value as it's how we calculate the min payment

	final := cnt * math.Pow(math.E, this.Rake * ((hoursBack / 24) / 365))

	minPay += final - cnt

	this.MinPayment = fmt.Sprintf("$%0.2f", minPay + 1)
	this.PayoffAmount = fmt.Sprintf("$%0.2f", final)

	return this.MinPayment, this.PayoffAmount 
}
