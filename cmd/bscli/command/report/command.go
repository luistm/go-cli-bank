package report

import (
	"encoding/csv"
	"os"

	"github.com/luistm/banksaurus/banksauruslib/usecases/listtransactions"
	"github.com/luistm/banksaurus/banksauruslib/usecases/listtransactionsgrouped"
	"github.com/luistm/banksaurus/cmd/bscli/adapter/cgdgateway"
	"github.com/luistm/banksaurus/cmd/bscli/adapter/presenterlisttransactions"
	"github.com/luistm/banksaurus/cmd/bscli/adapter/transactiongateway"
	"github.com/luistm/banksaurus/cmd/bscli/application"
)

// Command handles reports
type Command struct{}

// Execute the report command
func (rc *Command) Execute(arguments map[string]interface{}) error {
	var grouped bool

	// TODO: Lots of code in function, refactor

	hasFile := arguments["--input"].(bool)

	if arguments["--grouped"].(bool) {
		grouped = true
	}

	var lines [][]string
	if hasFile {

		filePath := arguments["<file>"].(string)
		_, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close()

		reader := csv.NewReader(file)
		reader.Comma = ';'
		reader.FieldsPerRecord = -1

		lines, err = reader.ReadAll()
		if err != nil {
			return err
		}
	}

	p, err := presenterlisttransactions.NewPresenter()
	if err != nil {
		return err
	}

	if !hasFile {
		db, err := application.Database()
		if err != nil {
			return err
		}

		repository, err := transactiongateway.NewTransactionRepository(db)
		if err != nil {
			return err
		}

		i, err := listtransactions.NewInteractor(p, repository)
		if err != nil {
			return err
		}

		err = i.Execute()
		if err != nil {
			return err
		}
	}

	if hasFile && grouped {
		repository, err := cgdgateway.New(lines)
		if err != nil {
			return err
		}

		i, err := listtransactionsgrouped.NewInteractor(repository, p)
		if err != nil {
			return err
		}

		err = i.Execute()
		if err != nil {
			return err
		}
	}

	if hasFile && !grouped {
		repository, err := cgdgateway.New(lines)
		if err != nil {
			return err
		}

		i, err := listtransactions.NewInteractor(p, repository)
		if err != nil {
			return err
		}

		err = i.Execute()
		if err != nil {
			return err
		}
	}

	vm, err := p.ViewModel()
	if err != nil {
		return err
	}

	vm.Write(os.Stdout)

	return nil
}
