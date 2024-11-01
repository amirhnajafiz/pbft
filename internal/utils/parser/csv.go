package parser

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/f24-cse535/pbft/pkg/models"
)

// CSVInput is used to parse the test-case files.
func CSVInput(path string) ([]*models.TestSet, error) {
	list := make([]*models.TestSet, 0)

	// open CSV file
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open the CSV file %s: %v", path, err)
	}
	defer file.Close()

	// create CSV reader
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1 // allow variable number of fields per row

	// read all data
	data, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV file %s: %v", path, err)
	}

	// loop variables
	var (
		index        = ""
		servers      []string
		byzantine    []string
		transactions []*models.Transaction
	)

	// read row by row
	for i, row := range data {
		// skip the headers
		if i == 0 {
			continue
		}

		if row[0] == "" { // old row
			tmp := strings.Split(strings.Replace(strings.Replace(row[1], "(", "", -1), ")", "", -1), ", ")
			transactions = append(transactions, &models.Transaction{
				Sender:   tmp[0],
				Receiver: tmp[1],
				Amount:   tmp[2],
			})
		} else {
			// save the current values
			if index != "" {
				list = append(list, &models.TestSet{
					Index:            index,
					LiveServers:      servers,
					ByzantineServers: byzantine,
					Transactions:     transactions,
				})
			}

			// reset values
			index = row[0]
			servers = make([]string, 0)
			byzantine = make([]string, 0)
			transactions = make([]*models.Transaction, 0)

			// set servers and byzantine servers
			servers = append(servers, strings.Split(strings.Replace(strings.Replace(row[2], "[", "", -1), "]", "", -1), ", ")...)
			byzantine = append(byzantine, strings.Split(strings.Replace(strings.Replace(row[3], "[", "", -1), "]", "", -1), ", ")...)

			// process the first row transactions
			tmp := strings.Split(strings.Replace(strings.Replace(row[1], "(", "", -1), ")", "", -1), ", ")
			transactions = append(transactions, &models.Transaction{
				Sender:   tmp[0],
				Receiver: tmp[1],
				Amount:   tmp[2],
			})
		}

		// save the last set
		if i == len(data)-1 {
			list = append(list, &models.TestSet{
				Index:            index,
				LiveServers:      servers,
				ByzantineServers: byzantine,
				Transactions:     transactions,
			})
		}
	}

	return list, nil
}
