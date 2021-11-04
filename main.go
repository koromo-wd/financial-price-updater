package main

import (
	"context"
	"log"

	"github.com/koromo-wd/priceupdater/oracle"
	"github.com/koromo-wd/priceupdater/updater"
	"gopkg.in/alecthomas/kingpin.v2"
)

const version = "1.2.0"
const coinGecko = "coingecko"
const coinMarketCap = "coinmarketcap"
const usd = "USD"

var (
	flagCryptoOracle         = kingpin.Flag("crypto-oracle", "Crypto oracle").PlaceHolder(coinGecko + "/" + coinMarketCap).Envar("CRYPTO_ORACLE").Default(coinGecko).String()
	coinGeckoTargetCryptoIDs = kingpin.Flag("coingecko-crypto-ids", "List of target Crypto IDs, used for CoinGecko").Envar("COINGECKO_CRYPTO_IDS").Default("bitcoin", "ethereum").Strings()
	cmcCryptoSymbols         = kingpin.Flag("crypto-symbols", "List of target Crypto symbols, used for CoinMarketCap").Envar("CMC_CRYPTO_SYMBOLS").Default("BTC", "ETH").Strings()
	cmcAPIKey                = kingpin.Flag("cmc-apikey", "CoinMarketCap API Key").Envar("CMC_API_KEY").String()
	googleSheetSAPath        = kingpin.Flag("gsheet-sa-path", "Path to Google Sheet service account token").Envar("GSHEET_SA_PATH").Default("/app/sa.json").String()
	googleSheetID            = kingpin.Flag("gsheet-id", "Google Sheet ID").Envar("GSHEET_ID").Required().String()
	googleSheetRange         = kingpin.Flag("gsheet-range", "Google Sheet range to work on").Envar("GSHEET_RANGE").Default("Sheet1!A1:B").String()
)

func main() {
	kingpin.Version(version)
	kingpin.Parse()

	ctx := context.Background()

	var cryptoOracle oracle.CryptoOracle
	var targetCryptos []string

	switch *flagCryptoOracle {
	case coinGecko:
		cryptoOracle = oracle.CoinGecko{}
		targetCryptos = *coinGeckoTargetCryptoIDs
	case coinMarketCap:
		cryptoOracle = oracle.CMC{APIKey: *cmcAPIKey}
		targetCryptos = *cmcCryptoSymbols
	default:
		log.Fatalf("Unmatched crypto oracle %s\n", *flagCryptoOracle)
	}

	quoteItems, err := cryptoOracle.GetQuoteItems(ctx, targetCryptos)
	if err != nil {
		log.Fatalf("Couldn't retrieve quote data from oracle: %s", err.Error())
	}

	priceUpdater := updater.NewGoogleSheet(
		*googleSheetSAPath,
		*googleSheetID,
		*googleSheetRange,
	)

	var tradingPairs []updater.TradingPair
	for _, v := range quoteItems {
		tradingPairs = append(tradingPairs, updater.TradingPair{
			BaseSymbol:  v.Symbol,
			QuoteSymbol: usd,
			Price:       v.USDPrice,
			UpdatedTime: v.LastUpdated,
		})
	}

	if err := priceUpdater.UpdatePrice(ctx, tradingPairs); err != nil {
		log.Fatalf("Couldn't update price: %s", err.Error())
	}

	log.Print("Finish updating price")
}
