#!/usr/bin/env bash

set -e

# Ensure resources directory
mkdir -p resources

# Generate private key (Ed25519)
openssl genpkey -algorithm Ed25519 -out resources/vayload-jwt-priv.pem

# Extract public key
openssl pkey -in resources/vayload-jwt-priv.pem -pubout -out resources/vayload-jwt-pub.pem

echo "Keys generated:"
echo " - resources/vayload-jwt-priv.pem"
echo " - resources/vayload-jwt-pub.pem"

echo "Base64 public key"
echo $(cat resources/vayload-jwt-pub.pem | base64 | tr -d '\n')
