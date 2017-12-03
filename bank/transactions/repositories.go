package transactions

import (
	"strings"

	"github.com/luistm/go-bank-cli/bank"
	"github.com/luistm/go-bank-cli/infrastructure"
	"github.com/luistm/go-bank-cli/lib/customerrors"
	"github.com/luistm/go-bank-cli/lib/sellers"
)

// NewRepository creates a repository for transactions
func NewRepository(storage infrastructure.CSVStorage) *repository {
	return &repository{storage: storage}
}

type repository struct {
	storage      bank.CSVHandler
	transactions []*Transaction
}

func (r *repository) GetAll() ([]*Transaction, error) {

	if r.storage == nil {
		return []*Transaction{}, customerrors.ErrInfrastructureUndefined
	}

	lines, err := r.storage.Lines()
	if err != nil {
		return []*Transaction{}, &customerrors.ErrInfrastructure{Msg: err.Error()}
	}

	// TODO: Validate if Lines() output is the expected one
	// r.validateLines(lines)
	// if err != nil{}
	r.buildTransactions(lines[5 : len(lines)-2])
	// log.Println(len(r.transactions))

	return r.transactions, nil
}

func (r *repository) buildTransactions(lines [][]string) error {

	for _, line := range lines {
		slug := strings.TrimSuffix(line[2], " ")
		t := &Transaction{
			value:  line[3],
			seller: sellers.New(slug, slug),
		}
		r.transactions = append(r.transactions, t)
	}

	return nil
}
