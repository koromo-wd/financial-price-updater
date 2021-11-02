package main

import (
	"context"
	"log"

	"example.com/financial-price-updater/oracle"
	"example.com/financial-price-updater/updater"
	"gopkg.in/alecthomas/kingpin.v2"
)

const version = "1.0.0"
const usd = "USD"

var (
	cmcAPIKey                     = kingpin.Flag("cmc-apikey", "CoinMarketCap API Key").Short('k').Envar("CMC_API_KEY").Required().String()
	targetCryptoSymbols           = kingpin.Flag("target-symbols", "List of target CryptoCurrency symbols").Short('t').Envar("TARGET_CRYPTO_SYMBOLS").Default("BTC", "ETH").Strings()
	googleSheetServiceAccountPath = kingpin.Flag("gsheet-sa-path", "Path to Google Sheet service account token").Short('s').Envar("GSHEET_SA_PATH").Required().String()
	googleSheetID                 = kingpin.Flag("gsheet-id", "Google Sheet ID").Short('i').Short('i').Envar("GSHEET_ID").Required().String()
	googleSheetRange              = kingpin.Flag("gsheet-range", "Google Sheet range to work on").Short('r').Envar("GSHEET_RANGE").Default("Sheet1!A1:B").String()
)

func main() {
	kingpin.Version(version)
	kingpin.Parse()

	ctx := context.Background()

	cryptoOracle := oracle.CMC{
		APIKey: *cmcAPIKey,
	}
	quoteItems, err := cryptoOracle.GetQuoteItems(ctx, *targetCryptoSymbols)
	if err != nil {
		log.Fatalf("Couldn't retrieve quote data from oracle: %s", err.Error())
	}

	priceUpdater := updater.NewGoogleSheet(
		*googleSheetServiceAccountPath,
		*googleSheetID,
		*googleSheetRange,
	)

	var tradingPairs []updater.TradingPair
	for _, v := range quoteItems {
		tradingPairs = append(tradingPairs, updater.TradingPair{
			BaseSymbol:  v.Symbol,
			QuoteSymbol: usd,
			Price:       v.Quote.USD.Price,
			UpdatedTime: v.LastUpdated,
		})
	}

	if err := priceUpdater.UpdatePrice(ctx, tradingPairs); err != nil {
		log.Fatalf("Couldn't update price: %s", err.Error())
	}

	log.Print("Finish updating price")
}
