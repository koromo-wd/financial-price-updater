#!/usr/bin/env bash

set -e

echo $GSHEET_SA > /app/sa.json
echo $GSHEET_OAUTH_CRED > /app/oauth-cred.json

exec /app/priceupdater "$@"
