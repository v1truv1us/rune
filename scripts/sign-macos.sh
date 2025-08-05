#!/bin/bash

BINARY_PATH="$1"
OS="$2"

if [[ "$OS" == "darwin" && -n "$MACOS_SIGN_IDENTITY" ]]; then
  echo "Signing macOS binary: $BINARY_PATH"
  codesign --sign "$MACOS_SIGN_IDENTITY" --timestamp --options=runtime --entitlements ./entitlements.plist "$BINARY_PATH"
  
  # Notarize if we have notarization credentials (skip for local testing)
  if [[ -n "$MACOS_NOTARY_ISSUER_ID" && -n "$MACOS_NOTARY_KEY_ID" && "$SKIP_NOTARIZATION" != "true" ]]; then
    echo "Notarizing macOS binary: $BINARY_PATH"
    
    # Create a temporary zip for notarization
    NOTARY_ZIP="/tmp/$(basename $BINARY_PATH)_notarization.zip"
    ditto -c -k --keepParent "$BINARY_PATH" "$NOTARY_ZIP"
    
    # Submit for notarization
    xcrun notarytool submit "$NOTARY_ZIP" \
      --issuer "$MACOS_NOTARY_ISSUER_ID" \
      --key-id "$MACOS_NOTARY_KEY_ID" \
      --key "$MACOS_NOTARY_KEY" \
      --wait \
      --timeout 10m
    
    # Clean up temporary zip
    rm -f "$NOTARY_ZIP"
    echo "Notarization completed for: $BINARY_PATH"
  else
    echo "Skipping notarization (SKIP_NOTARIZATION=true or missing credentials)"
  fi
fi