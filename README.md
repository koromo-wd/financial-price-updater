# Financial Price Updater

## Current use case

- Cryptocurrency
- Thai Mutual Fund

## Supported Oracle

- CoinGecko
- CoinMarketCap
- Thai SEC Open API

## Supported Update destination

- Google Sheet (can be authenticated using oauth or service account)

## How to run

Run a binary and specify flags or set ENV needed for the use case.

(You can either build the binary from source or using docker image)

### Example

Updating crypto price

```bash
./priceupdater crypto --coingecko-crypto-ids=bitcoin,ethereum,cardano,polkadot --gsheet-oauth-cred-path={yourOauthCredentialPath} --gsheet-oauth-token-path={pathToStoreOauthToken} --gsheet-id={yourGSheetID}
```

or

```bash
export GSHEET_OAUTH_CRED_PATH={yourOauthCredentialPath}
export GSHEET_OAUTH_TOKEN_PATH={pathToStoreOauthToken}
export GSHEET_ID={yourGSheetID}
export COINGECKO_CRYPTO_IDS=bitcoin,ethereum,cardano,polkadot

./priceupdater crypto
```

Updating mutualfund price

```bash
export GSHEET_OAUTH_CRED_PATH={yourOauthCredentialPath}
export GSHEET_OAUTH_TOKEN_PATH={pathToStoreOauthToken}
export GSHEET_ID={yourGSheetID}
export THSEC_FFACT_API_KEY={apiKey1}
export THSEC_FDAILY_API_KEY={apiKey2}
export THSEC_FUND_NAMES=SCBNK225,SCBEUEQ

./priceupdater fund
```

## TODO

- Add stock support
