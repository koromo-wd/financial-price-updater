package updater

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var headerRow = []interface{}{"Pair", "Price", "Updated Time"}

type GoogleSheet struct {
	Option     option.ClientOption
	SheetID    string
	WriteRange string
}

func NewGoogleSheet(serviceAccountTokenPath, sheetID, writeRange string) *GoogleSheet {
	return &GoogleSheet{
		Option:     option.WithCredentialsFile(serviceAccountTokenPath),
		SheetID:    sheetID,
		WriteRange: writeRange,
	}
}

func (updater GoogleSheet) UpdatePrice(ctx context.Context, tradingPairs []TradingPair) error {
	svc, err := sheets.NewService(ctx, updater.Option)
	if err != nil {
		return err
	}

	writeVal := [][]interface{}{}
	writeVal = append(writeVal, headerRow)

	for _, pair := range tradingPairs {
		writeVal = append(writeVal, []interface{}{
			fmt.Sprintf("%s/%s", pair.BaseSymbol, pair.QuoteSymbol),
			pair.Price,
			pair.UpdatedTime.Local().Format(time.RFC1123),
		})
	}

	_, err = svc.Spreadsheets.Values.Update(updater.SheetID, updater.WriteRange, &sheets.ValueRange{Values: writeVal}).ValueInputOption("USER_ENTERED").Do()
	if err != nil {
		return fmt.Errorf("unable to write data to sheet: %w", err)
	}

	return nil
}
