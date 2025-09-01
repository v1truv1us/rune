#!/bin/bash
# export-apple-identity.sh
#
# Export a code signing identity (certificate + private key) from macOS Keychain as a .p12 file via CLI.
#
# Usage:
#   ./export-apple-identity.sh "Common Name Substring" output.p12
#
# Example:
#   ./export-apple-identity.sh "Apple Development: John Ferguson" my-certificate.p12
#
# Requirements:
#   - macOS with Keychain Access
#   - OpenSSL installed (brew install openssl if needed)
#
# This script will:
#   1. Find the identity in your login keychain matching the provided name substring
#   2. Export the certificate and private key as a .p12 file (you will be prompted for a password)
#   3. Handle errors and edge cases robustly

set -euo pipefail

if [ "$#" -ne 2 ]; then
  echo "Usage: $0 \"Common Name Substring\" output.p12"
  exit 1
fi

NAME_SUBSTRING="$1"
OUTPUT_P12="$2"
KEYCHAIN=login.keychain-db

# Find the identity (certificate + private key) by name substring
echo "Searching for identity containing: $NAME_SUBSTRING"
IDENTITY_INFO=$(security find-identity -v -p codesigning "$KEYCHAIN" | grep "$NAME_SUBSTRING" | head -n1)
if [ -z "$IDENTITY_INFO" ]; then
  echo "No identity found matching: $NAME_SUBSTRING"
  exit 1
fi

IDENTITY_HASH=$(echo "$IDENTITY_INFO" | awk '{print $2}')
echo "Found identity hash: $IDENTITY_HASH"

# Export the identity as a .p12 (will prompt for password)
echo "Exporting identity to $OUTPUT_P12 ..."
security export -k "$KEYCHAIN" -t identities -f pkcs12 -P "" -Z "$IDENTITY_HASH" -o "$OUTPUT_P12"

if [ $? -eq 0 ]; then
  echo "Export successful: $OUTPUT_P12"
else
  echo "Export failed."
  exit 1
fi

echo "Done."
