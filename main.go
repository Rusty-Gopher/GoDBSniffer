// main.go
package main

import (
	"GoDBSniffer/config"
	"GoDBSniffer/pkg/database"
	"GoDBSniffer/pkg/health"
	"GoDBSniffer/pkg/performance"
	"GoDBSniffer/pkg/security"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func main() {
	// Step 1: Letting the user know that
	asciiArt := figure.NewFigure("GoDBSniffer is starting...", "", true)
	asciiArt.Print()
	fmt.Println()

	// Step 2: Ask for DB connection details
	printSeparator()
	fmt.Println("Please enter the database connection details:")
	dbConfig := promptDBDetails()

	// Step 3: Store the details for later use
	printSeparator()
	color.Cyan("Storing database connection details...")
	config.SetDBConfig(dbConfig)

	// Step 4: Ask user for a sniff with a warning about potential time consumption
	printSeparator()
	warning := color.New(color.FgHiYellow, color.Bold)
	warning.Println("WARNING: The sniffing process may take some time, depending on the database size and performance.")

	if askForSniff() {
		// Step 5: Perform the sniffing process
		fmt.Println("Starting the sniffing process...")
		executeSniff()
	}
}

// promptDBDetails prompts user for DB connection details and returns config.DBConfig
func promptDBDetails() config.DBConfig {
	var dbConfig config.DBConfig
	prompt := []*survey.Question{
		{
			Name: "host",
			Prompt: &survey.Input{
				Message: "What is your database host?",
				Default: "localhost",
			},
		},
		{
			Name: "port",
			Prompt: &survey.Input{
				Message: "What is your database port?",
				Default: "3306",
			},
		},
		{
			Name: "user",
			Prompt: &survey.Input{
				Message: "What is your database user?",
				Default: "root",
			},
		},
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: "What is your database password?",
			},
		},
		{
			Name: "database",
			Prompt: &survey.Input{
				Message: "What is your database name?",
				Default: "mydatabase",
			},
		},
	}

	survey.Ask(prompt, &dbConfig)
	return dbConfig
}

// askForSniff prompts the user to confirm if they want to perform a database sniff
func askForSniff() bool {
	confirm := false // This will store the user's confirmation

	// Define the prompt
	prompt := &survey.Confirm{
		Message: "Do you want to perform a database sniff? It may take some time.",
	}

	// Ask the user for confirmation
	err := survey.AskOne(prompt, &confirm)
	if err != nil {
		log.Fatalf("Prompt failed: %v\n", err)
	}

	return confirm // Return the user's confirmation
}

// Ths will start the process by performing health, security, and performance checks
func executeSniff() {
	db, err := database.ConnectDatabase(config.GetDBConfig())
	if err != nil {
		fmt.Printf("Error connecting to the database: %v\n", err)
		return // Stop if it cannot connect to the database
	}
	fmt.Println("Successfully connected to the database.")
	defer db.Close()

	// once connected, Check the existing tables and print them, using fancy tablewriter package
	displayDatabaseInfo(db)

	// Step 6: Check that health...
	printSeparator()
	if err := health.CheckDatabaseUptime(db); err != nil {
		fmt.Printf("Health checks failed: %v\n", err)
	}

	// step 7: Check that security
	printSeparator()
	if err := security.BasicSecurityCheck(db); err != nil {
		fmt.Printf("Security checks failed: %v\n", err)
	}

	// step 7: Check that basic perormance, at this point i am just adding some comments
	printSeparator()
	if err := performance.BasicPerformanceCheck(db); err != nil {
		fmt.Printf("Performance checks failed: %v\n", err)
	}

	// want to make it very clear to the user
	color.Set(color.FgHiYellow, color.Bold)
	fmt.Println("\nAll checks completed successfully but for more details please check their respective tables, this is where you will find the nuggets.")
	color.Unset()

}

// Utility functions for styling
func printSeparator() {
	color.Blue("================================================================================")
}

func displayDatabaseInfo(db *sql.DB) {
	tables, err := database.GetTables(db)
	if err != nil {
		log.Fatalf("Error retrieving tables: %v", err)
	}

	fmt.Println("\nTables in the database:")
	displayCount := 0 // this is a counter to make sure we arent displaying more than 5

	for _, table := range tables {
		if displayCount >= 5 { // display only 5, yes i said 5, thats totally random
			break
		}
		fmt.Printf("\n%s\n", table)
		columns, err := database.GetTableSchema(db, table)
		if err != nil {
			log.Fatalf("Error retrieving schema for table %s: %v", table, err)
		}

		tw := tablewriter.NewWriter(os.Stdout)
		tw.SetHeader([]string{"Column", "Type", "Null", "Key", "Default", "Extra"})
		columnCount := 0 // same counter for column brahh

		for _, col := range columns {
			if columnCount >= 10 { // My majesty, dont you get it, this is the column counter till 10
				break
			}
			tw.Append([]string{col.Field, col.Type, col.Null, col.Key, col.Default.String, col.Extra})
			columnCount++
		}
		tw.Render() // This is the time to render datt, oh yeah renderrrr, insert office meme dancing
		displayCount++
	}

	if len(tables) > 5 {
		fmt.Printf("\nDisplayed %d of %d tables. For more details, check the database directly.\n", displayCount, len(tables))
	}
}
