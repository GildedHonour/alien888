package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func updateGoogleSheets__Test() {
	spreadsheetID := "your-spreadsheet-id"
	readRange := "Sheet1!A:C" // Assumes the data is in the first sheet and spans columns A to C

	ctx := context.Background()

	// Load the credentials JSON file for the Google Sheets API
	credentialsFile := "path/to/credentials.json"
	credentials, err := os.ReadFile(credentialsFile)
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// Create a new Google Sheets service client
	srv, err := sheets.NewService(ctx, option.WithCredentialsJSON(credentials))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	// Make the API call to retrieve data from the spreadsheet
	resp, err := srv.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet: %v", err)
	}

	// Process the response and print the values
	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
	} else {
		fmt.Println("Data:")

		for _, row := range resp.Values {
			for _, value := range row {
				fmt.Printf("%s\t", value)
			}
			fmt.Println()
		}
	}
}
