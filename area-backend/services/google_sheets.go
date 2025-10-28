package services

import (
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type GoogleSheetsService struct {
	client *sheets.Service
}

func NewGoogleSheetsService() (*GoogleSheetsService, error) {
	apiKey := os.Getenv("GOOGLE_SHEETS_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_SHEETS_API_KEY is not configured")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	service, err := sheets.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to initialise google sheets client: %w", err)
	}

	return &GoogleSheetsService{client: service}, nil
}

func (s *GoogleSheetsService) FetchValues(spreadsheetID, readRange string) ([][]string, error) {
	if spreadsheetID == "" {
		return nil, fmt.Errorf("spreadsheetID is required")
	}
	if readRange == "" {
		return nil, fmt.Errorf("range is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := s.client.Spreadsheets.Values.Get(spreadsheetID, readRange).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch sheet values: %w", err)
	}

	rows := make([][]string, len(resp.Values))
	for i, row := range resp.Values {
		rowValues := make([]string, len(row))
		for j, cell := range row {
			rowValues[j] = fmt.Sprintf("%v", cell)
		}
		rows[i] = rowValues
	}

	return rows, nil
}
