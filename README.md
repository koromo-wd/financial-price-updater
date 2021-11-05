# Financial Price Updater

## Current use case

- Cryptocurrency

## Supported Oracle

- CoinGecko
- CoinMarketCap

## Supported Update destination

- Google Sheet (can be authenticated using oauth or service account)

## How to run

Run a binary and specify flags or set ENV needed for the use case.

(You can either build the binary from source or using docker image)

### Example

```bash
./priceupdater --coingecko-crypto-ids=bitcoin,ethereum,cardano,polkadot --gsheet-oauth-cred-path={yourOauthCredentialPath} --gsheet-oauth-token-path={pathToStoreOauthToken} --gsheet-id={yourGSheetID}
```

or

```bash
export GSHEET_OAUTH_CRED_PATH=/app/oauth-cred.json
export GSHEET_OAUTH_TOKEN_PATH=/app/token.json
export GSHEET_ID={yourGSheetID}
export COINGECKO_CRYPTO_IDS=bitcoin,ethereum,cardano,polkadot

./priceupdater
```

## TODO

- Add stock support
- Add mutual fund support
