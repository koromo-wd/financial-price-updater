#!/usr/bin/env bash

echo $GSHEET_SA > /app/sa.json

exec "$@"
