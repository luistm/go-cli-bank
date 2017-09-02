package reports

import (
	"expensetracker/entities"
	"fmt"

	"github.com/shopspring/decimal"
)

// MonthlyReport builds a sum of expenses and credits
func MonthlyReport(records [][]string) error {

	// TODO: Create an entity report
	var report map[string]decimal.Decimal
	report = make(map[string]decimal.Decimal)
	var credit decimal.Decimal
	var expense decimal.Decimal

	// Read all transactions
	for lineCount, record := range records {

		r := entities.Record{Record: record}

		if !r.Valid() || lineCount < 4 {
			continue
		}

		t := entities.Transaction{}
		transaction := t.New(r)

		if transaction.IsFromThisMonth() {
			report[transaction.Description] = report[transaction.Description].Add(transaction.Value())
			if transaction.IsDebt() {
				expense = expense.Add(transaction.Value())
			} else {
				credit = credit.Add(transaction.Value())
			}
		}

	}

	for transactionDescription, transactionValue := range report {
		fmt.Printf("%24s %8s \n", transactionDescription, transactionValue.String())
	}

	fmt.Println("Expense is ", expense.String())
	fmt.Println("Credit is ", credit.String())

	return nil
}