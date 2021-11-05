package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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

func NewGoogleSheetOAuth(credentialPath, tokenStoredPath, sheetID, writeRange string) (*GoogleSheet, error) {
	b, err := ioutil.ReadFile(credentialPath)
	if err != nil {
		return nil, fmt.Errorf("fail to read google sheet oauth credential from path")
	}
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, fmt.Errorf("fail to get google sheet config from json")
	}

	client, err := getClient(config, tokenStoredPath)
	if err != nil {
		return nil, fmt.Errorf("fail to get client: %w", err)
	}

	return &GoogleSheet{
		Option:     option.WithHTTPClient(client),
		SheetID:    sheetID,
		WriteRange: writeRange,
	}, nil
}

func (updater GoogleSheet) UpdatePrice(ctx context.Context, tradingPairs []TradingPair) error {
	svc, err := sheets.NewService(ctx, updater.Option)
	if err != nil {
		return err
	}

	if err := deleteExistingCells(svc, updater.SheetID, updater.WriteRange); err != nil {
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

func deleteExistingCells(svc *sheets.Service, sheetID, clearRange string) error {
	if _, err := svc.Spreadsheets.Values.Clear(sheetID, clearRange, &sheets.ClearValuesRequest{}).Do(); err != nil {
		return err
	}
	return nil
}

func getClient(config *oauth2.Config, tokenStoredPath string) (*http.Client, error) {
	ctx := context.Background()

	tok, err := tokenFromFile(tokenStoredPath)
	if err != nil {
		tok, err = getTokenFromWeb(ctx, config)
		if err != nil {
			return nil, err
		}

		if err = saveToken(tokenStoredPath, tok); err != nil {
			return nil, err
		}
	}

	return config.Client(ctx, tok), nil
}

func getTokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		return nil, fmt.Errorf("unable to read authorization code: %w", err)
	}

	tok, err := config.Exchange(ctx, authCode)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve token from web: %w", err)
	}

	return tok, nil
}

func tokenFromFile(filePath string) (*oauth2.Token, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)

	return tok, err
}

func saveToken(path string, token *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %w", err)
	}
	defer f.Close()

	json.NewEncoder(f).Encode(token)

	return nil
}
