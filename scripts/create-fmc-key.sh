#!/usr/bin/env bash

SECRET="$1"

if [ -z "$SECRET" ]; then
  echo "Use: $0 <HMAC_SECRET>"
  exit 1
fi

# Header and payload (base64url without padding)
HEADER=$(echo -n '{"alg":"HS256","typ":"JWT"}' | openssl base64 -A | tr '+/' '-_' | tr -d '=')
PAYLOAD=$(echo -n "{\"ts\":$(date +%s)}" | openssl base64 -A | tr '+/' '-_' | tr -d '=')

DATA="$HEADER.$PAYLOAD"

# Sign with HMAC-SHA256 and encode base64url
SIGNATURE=$(echo -n "$DATA" | openssl dgst -sha256 -hmac "$SECRET" -binary | openssl base64 -A | tr '+/' '-_' | tr -d '=')

TOKEN="$DATA.$SIGNATURE"
echo "$TOKEN"
