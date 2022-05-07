/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// Variable to store the file location of the db.sqlite file
var file string

// Variable to store the output directory for the generated csv files
var output string

// YNABCmd represents the YNAB command
var YNABCmd = &cobra.Command{
	Use:   "YNAB",
	Short: "You Need A Budget CSV Export",
	Long:  `Exports the account data in a csv format that is consumable by YNAB web edition.`,
	Run: func(cmd *cobra.Command, args []string) {
		ynabExport(file, output)
	},
}

func init() {
	exportCmd.AddCommand(YNABCmd)

	// Add Flags for `file` and `output` to get user input on file and output locations
	YNABCmd.Flags().StringVarP(&file, "file", "f", "db.sqlite", "Path to the exported sqlite file from Actual")
	YNABCmd.Flags().StringVarP(&output, "output", "o", "", "Path to store the generated csv output")

	// Mark `file` and `output` as required
	YNABCmd.MarkFlagRequired("file")
	YNABCmd.MarkFlagRequired("output")

}

// ynabExport processes the inputs from the command and generates the required CSV files.
func ynabExport(file string, output string) {

	// Check that input file exists
	bf, _ := exists(file)

	// Throw error code 1 if missing
	if !bf {
		os.Exit(1)
	}

	// Try to open the sqlite file.
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		log.Fatalf("Unable to open %s, received error: %s", file, err)
	}

	// Defer closing of the db connection
	defer db.Close()

	// Query for all accounts
	row, err := db.Query("SELECT id, name, type from accounts")
	if err != nil {
		log.Fatal(err)
	}
	defer row.Close()

	// Create slice to store accounts
	var accounts []Account

	for row.Next() {
		var id string
		var name string
		var types string
		row.Scan(&id, &name, &types)
		a := Account{id: id, name: name, accountType: types}
		accounts = append(accounts, a)
		log.Printf("Account: %s - %s - %s\n", id, name, types)
	}

	log.Printf("A total of %v accounts were retrieved", len(accounts))

	for _, a := range accounts {
		log.Printf("Retrieving transactions for %s", a.name)
		row, err := db.Query(fmt.Sprintf("SELECT date, payees.name, notes, amount FROM v_transactions INNER JOIN payees on payee = payees.id WHERE account = \"%s\"", a.id))
		if err != nil {
			log.Fatal(err)
		}
		defer row.Close()

		var transactions []Transaction

		for row.Next() {
			var date int
			var name string
			var notes string
			var amount int
			row.Scan(&date, &name, &notes, &amount)
			t := Transaction{
				date:   date,
				name:   name,
				notes:  notes,
				amount: amount,
			}
			transactions = append(transactions, t)
			log.Printf("A total of %v transaction(s) were retrieved for %s", len(transactions), a.name)

			// Attempt to create the CSV file for the account

			path := filepath.Join(output, fmt.Sprintf("%s.csv", a.name))
			csvFile, err := os.Create(path)
			if err != nil {
				log.Fatal(err)
			}
			defer csvFile.Close()

			writer := csv.NewWriter(csvFile)
			defer writer.Flush()

			header := []string{"date", "name", "notes", "amount"}

			err = writer.Write(header)
			if err != nil {
				log.Fatal(err)
			}

			for _, t := range transactions {
				writer.Write(t.StringArray())
			}

			log.Printf("%s was created succesfully!", path)
		}
	}

}

// exists checks if a file or directory exists on the system
func exists(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		return true, nil
	} else if errors.Is(err, os.ErrNotExist) {
		log.Fatalf("The file %s was not found\n", file)
		return false, err
	} else {
		log.Fatalf("Unable to determine if file exists, please check the error: %s", err)
		return false, err
	}
}

// Account holds a single account entity from Actual
type Account struct {
	id          string
	name        string
	accountType string
}

type Transaction struct {
	date   int
	name   string
	notes  string
	amount int
}

func (t Transaction) String() string {
	return fmt.Sprintf("[%d, %s, %s, %d]", t.date, t.name, t.notes, t.amount)
}

func (t Transaction) StringArray() []string {
	d := strconv.Itoa(t.date)
	year := d[:4]
	month := d[4:6]
	day := d[6:]
	a := strconv.Itoa(t.amount)
	return []string{fmt.Sprintf("%s/%s/%s", year, month, day), t.name, t.notes, a}
}
