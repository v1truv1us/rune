#!/bin/bash

BINARY_PATH="$1"
OS="$2"

if [[ "$OS" == "darwin" && -n "$MACOS_SIGN_IDENTITY" ]]; then
  echo "Signing macOS binary: $BINARY_PATH"
  codesign --sign "$MACOS_SIGN_IDENTITY" --timestamp --options=runtime --entitlements ./entitlements.plist "$BINARY_PATH"
  
  # Notarize if we have notarization credentials (skip for local testing)
  if [[ -n "$MACOS_NOTARY_ISSUER_ID" && -n "$MACOS_NOTARY_KEY_ID" && "$SKIP_NOTARIZATION" != "true" ]]; then
    echo "Notarizing macOS binary: $BINARY_PATH"
    
    # Handle both file path (local) and file content (GitHub Actions)
    if [[ -f "$MACOS_NOTARY_KEY" ]]; then
      # Local development - MACOS_NOTARY_KEY is a file path
      NOTARY_KEY_PATH="$MACOS_NOTARY_KEY"
    else
      # GitHub Actions - MACOS_NOTARY_KEY contains the file content
      NOTARY_KEY_PATH="/tmp/AuthKey_$(date +%s).p8"
      echo "$MACOS_NOTARY_KEY" > "$NOTARY_KEY_PATH"
      echo "Created temporary key file: $NOTARY_KEY_PATH"
    fi
    
    # Create a temporary zip for notarization
    NOTARY_ZIP="/tmp/$(basename $BINARY_PATH)_notarization.zip"
    ditto -c -k --keepParent "$BINARY_PATH" "$NOTARY_ZIP"
    
    # Submit for notarization
    xcrun notarytool submit "$NOTARY_ZIP" \
      --issuer "$MACOS_NOTARY_ISSUER_ID" \
      --key-id "$MACOS_NOTARY_KEY_ID" \
      --key "$NOTARY_KEY_PATH" \
      --wait \
      --timeout 10m
    
    # Clean up temporary files
    rm -f "$NOTARY_ZIP"
    if [[ "$NOTARY_KEY_PATH" != "$MACOS_NOTARY_KEY" ]]; then
      rm -f "$NOTARY_KEY_PATH"
    fi
    echo "Notarization completed for: $BINARY_PATH"
  else
    echo "Skipping notarization (SKIP_NOTARIZATION=true or missing credentials)"
  fi
fi