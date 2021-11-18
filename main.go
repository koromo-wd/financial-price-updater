package main

import (
	"context"
	"log"

	"github.com/koromo-wd/priceupdater/oracle"
	"github.com/koromo-wd/priceupdater/updater"
	"gopkg.in/alecthomas/kingpin.v2"
)

const version = "1.3.0"
const coinGecko = "coingecko"
const coinMarketCap = "coinmarketcap"
const gsheetUpdaterSa = "gsheet-sa"
const gsheetUpdaterOauth = "gsheet-oauth"

var (
	flagUpdater              = kingpin.Flag("updater", "updater to use").Envar("UPDATER").Default(gsheetUpdaterOauth).String()
	googleSheetSAPath        = kingpin.Flag("gsheet-sa-path", "Path to Google Sheet service account token").Envar("GSHEET_SA_PATH").Default("/app/sa.json").String()
	googleSheetOauthCredPath = kingpin.Flag("gsheet-oauth-cred-path", "Path to Google Sheet oauth credential").Envar("GSHEET_OAUTH_CRED_PATH").Default("/app/oauth-cred.json").String()
	googleSheetOauthTokPath  = kingpin.Flag("gsheet-oauth-token-path", "Path to Google Sheet stored token").Envar("GSHEET_OAUTH_TOKEN_PATH").Default("/tmp/oauth-token.json").String()
	googleSheetID            = kingpin.Flag("gsheet-id", "Google Sheet ID").Envar("GSHEET_ID").Required().String()
	googleSheetRange         = kingpin.Flag("gsheet-range", "Google Sheet range to work on").Envar("GSHEET_RANGE").Default("Sheet1!A1:B").String()

	cryptoCommand            = kingpin.Command("crypto", "Update crypto price")
	flagCryptoOracle         = cryptoCommand.Flag("crypto-oracle", "Crypto oracle").PlaceHolder(coinGecko + "/" + coinMarketCap).Envar("CRYPTO_ORACLE").Default(coinGecko).String()
	coinGeckoTargetCryptoIDs = cryptoCommand.Flag("coingecko-crypto-ids", "List of target Crypto IDs, used for CoinGecko").Envar("COINGECKO_CRYPTO_IDS").Default("bitcoin", "ethereum").Strings()
	cmcCryptoSymbols         = cryptoCommand.Flag("crypto-symbols", "List of target Crypto symbols, used for CoinMarketCap").Envar("CMC_CRYPTO_SYMBOLS").Default("BTC", "ETH").Strings()
	cmcAPIKey                = cryptoCommand.Flag("cmc-apikey", "CoinMarketCap API Key").Envar("CMC_API_KEY").String()

	fundCommand            = kingpin.Command("fund", "Update mutual fund price")
	thaiSecFundDailyAPIKey = fundCommand.Flag("thsec-fdaily-apikey", "Thai Sec Fund Daily Info API Key").Envar("THSEC_FDAILY_API_KEY").String()
	thaiSecFundFactAPIKey  = fundCommand.Flag("thsec-ffact-apikey", "Thai Sec Fund Fact API Key").Envar("THSEC_FFACT_API_KEY").String()
	thaiSecFundNames       = fundCommand.Flag("thsec-fund-names", "List of target fund names, used for Thai Sec API").Envar("THSEC_FUND_NAMES").Strings()
)

func main() {
	kingpin.Version(version)
	ctx := context.Background()
	var quoteItems []oracle.QuoteItem
	var err error

	switch kingpin.Parse() {
	case cryptoCommand.FullCommand():
		log.Print("Updating Crypto price")

		var cryptoOracle oracle.Oracle
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

		quoteItems, err = cryptoOracle.GetQuoteItems(ctx, targetCryptos)
		if err != nil {
			log.Fatalf("Couldn't retrieve quote data from oracle: %s", err.Error())
		}
	case fundCommand.FullCommand():
		log.Print("Updating mutual fund price")
		fundOracle := oracle.ThaiSec{
			FundFactAPIKey:      *thaiSecFundFactAPIKey,
			FundDailyInfoAPIKey: *thaiSecFundDailyAPIKey,
		}

		quoteItems, err = fundOracle.GetQuoteItems(ctx, *thaiSecFundNames)
		if err != nil {
			log.Fatalf("Couldn't retrieve quote data from oracle: %s", err.Error())
		}
	}

	tradingPairs := createTradingPairs(quoteItems)

	priceUpdater := getPriceUpdater()
	if err := priceUpdater.UpdatePrice(ctx, tradingPairs); err != nil {
		log.Fatalf("Couldn't update price: %s", err.Error())
	}

	log.Print("Finish updating price")
}

func getPriceUpdater() updater.Updater {
	switch *flagUpdater {
	case gsheetUpdaterSa:
		return updater.NewGoogleSheet(
			*googleSheetSAPath,
			*googleSheetID,
			*googleSheetRange,
		)
	case gsheetUpdaterOauth:
		priceUpdater, err := updater.NewGoogleSheetOAuth(
			*googleSheetOauthCredPath,
			*googleSheetOauthTokPath,
			*googleSheetID,
			*googleSheetRange,
		)
		if err != nil {
			log.Fatalf("Couldn't initialize updater: %s", err.Error())
		}
		return priceUpdater
	default:
		log.Fatalf("Unmatched updater %s\n", *flagUpdater)
	}

	return nil
}

func createTradingPairs(quoteItems []oracle.QuoteItem) []updater.TradingPair {
	var out []updater.TradingPair
	for _, v := range quoteItems {
		out = append(out, updater.TradingPair{
			BaseSymbol:  v.Symbol,
			QuoteSymbol: v.BaseCurrency,
			Price:       v.Price,
			UpdatedTime: v.LastUpdated,
		})
	}
	return out
}
