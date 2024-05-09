package models 

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"encoding/json"
)

func TestQALedgerPayoff1a (t *testing.T) {
	ledge := &Ledger{
		Rake: 0.07,
	}

	minPay, pay := ledge.Payoff()

	assert.Equal (t, "$0", minPay)
	assert.Equal (t, "$0", pay)

	ledge.LineItems = append(ledge.LineItems, LineItem {
		Date: time.Now().AddDate(-1, 0, 0),
		Note: "regression",
		Amount: 300000,
	})

	minPay, pay = ledge.Payoff()
	assert.Equal (t, "$3218.14", pay)
	assert.Equal (t, "$219.14", minPay)
}

// adds another 3k half way through
func TestQALedgerPayoff1b (t *testing.T) {
	ledge := &Ledger{
		Rake: 0.07,
	}

	ledge.LineItems = append(ledge.LineItems, LineItem {
		Date: time.Now().AddDate(-1, 0, 0),
		Note: "regression",
		Amount: 300000,
	})

	ledge.LineItems = append(ledge.LineItems, LineItem {
		Date: time.Now().AddDate(0, -6, 0),
		Note: "regression",
		Amount: 300000,
	})

	minPay, pay := ledge.Payoff()
	assert.Equal (t, "$6436.28", pay)
	assert.Equal (t, "$437.28", minPay)
}

// adds a payment in there
func TestQALedgerPayoff1c (t *testing.T) {
	ledge := &Ledger{
		Rake: 0.07,
	}

	ledge.LineItems = append(ledge.LineItems, LineItem {
		Date: time.Now().AddDate(-1, 0, 0),
		Note: "regression",
		Amount: 300000,
	})

	ledge.LineItems = append(ledge.LineItems, LineItem {
		Date: time.Now().AddDate(0, -6, 0),
		Note: "regression",
		Amount: 300000,
	})

	ledge.LineItems = append(ledge.LineItems, LineItem {
		Date: time.Now().AddDate(0, -5, 0),
		Note: "regression",
		Amount: -30000,
	})


	minPay, pay := ledge.Payoff()
	assert.Equal (t, "$6127.41", pay)
	assert.Equal (t, "$177.04", minPay)
}


func TestQALedgerPayoff1d (t *testing.T) {
	ledge := &Ledger{
		Rake: 0.07,
	}

	ledge.LineItems = append(ledge.LineItems, LineItem {
		Date: time.Now().AddDate(-1, 0, 0),
		Note: "regression",
		Amount: 300000,
	})

	ledge.LineItems = append(ledge.LineItems, LineItem {
		Date: time.Now().AddDate(0, -6, 0),
		Note: "regression",
		Amount: 300000,
	})

	ledge.LineItems = append(ledge.LineItems, LineItem {
		Date: time.Now().AddDate(0, -5, 0),
		Note: "regression",
		Amount: -30000,
	})

	ledge.LineItems = append(ledge.LineItems, LineItem {
		Date: time.Now().AddDate(0, -4, 0),
		Note: "regression",
		Amount: 100000,
	})

	ledge.LineItems = append(ledge.LineItems, LineItem {
		Date: time.Now().AddDate(0, -3, 0),
		Note: "regression",
		Amount: -30000,
	})

	ledge.LineItems = append(ledge.LineItems, LineItem {
		Date: time.Now().AddDate(0, -2, 0),
		Note: "regression",
		Amount: -30000,
	})

	minPay, pay := ledge.Payoff()

	assert.Equal (t, "$6548.24", pay)
	assert.Equal (t, "$75.92", minPay)
}

// actual data
func TestQALedgerPayoff2a (t *testing.T) {
	ledge := &Ledger{}

	// err := json.Unmarshal([]byte(`{"Rake":0.1,"LineItems":[{"Date":"2024-04-02T18:00:00Z","Note":"tie rod ends","Amount":12957},{"Date":"2024-04-05T18:00:00Z","Note":"Steering rack","Amount":350654},{"Date":"2024-05-09T10:25:20.912300258Z","Note":"Personal Loan","Amount":120000}],"MinPayment":"","PayoffAmount":"","Name":""}`), ledge)
	err := json.Unmarshal([]byte(`{"Rake":0.1,"LineItems":[{"Date":"2024-04-02T18:00:00Z","Note":"tie rod ends","Amount":12957},{"Date":"2024-04-05T18:00:00Z","Note":"Steering rack","Amount":350654}],"MinPayment":"","PayoffAmount":"","Name":""}`), ledge)
	if err != nil { t.Fatal(err) }

	minPay, pay := ledge.Payoff()

	t.Logf("min: %s :: payoff: %s\n", minPay, pay)


}
