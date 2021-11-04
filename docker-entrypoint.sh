#!/usr/bin/env bash

set -e

echo $GSHEET_SA > /app/sa.json

exec /app/priceupdater "$@"
