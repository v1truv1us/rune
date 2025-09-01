# Apple Code Signing and Notarization Setup

This document explains how to set up Apple Developer ID code signing and notarization for the Rune CLI project.

## Prerequisites

1. **Apple Developer Account**: You need an active Apple Developer account with Developer ID capabilities
2. **Apple Developer ID Certificate**: A valid Developer ID Application certificate from Apple
3. **App Store Connect API Key**: For notarization service

---

## Understanding .cer vs .p12 Files

When working with Apple code signing, you will encounter both `.cer` and `.p12` files. They serve different purposes:

### What is a `.cer` file?
- **.cer** stands for certificate.
- Contains the **public key** and certificate holder info.
- Downloaded from the Apple Developer portal after creating a certificate.
- **Does not contain the private key** (cannot be used for signing on its own).

### What is a `.p12` file?
- **.p12** (PKCS#12) is a password-protected file containing both the **public certificate** and the **private key**.
- Required for code signing and exporting certificates to other machines or CI systems.
- Created by exporting your certificate and private key from Keychain Access.

### Why does Apple only give you a `.cer` file?
- The **private key** is generated and stored locally on your Mac when you create a Certificate Signing Request (CSR).
- Apple only issues the public certificate (`.cer`).
- You must combine the `.cer` with your private key (in Keychain) to export a `.p12`.

### Summary Table

| File Type | Contains           | Use Case                | How to Get It                |
|-----------|--------------------|-------------------------|------------------------------|
| `.cer`    | Public certificate | Verifying identity      | Download from Apple portal   |
| `.p12`    | Cert + Private key | Code signing, exporting | Export from Keychain Access  |

### How to Export a .p12 File from Keychain Access

1. **Double-click your downloaded `.cer` file** to add it to Keychain Access (if not already present).
2. Open **Keychain Access** (`/Applications/Utilities/Keychain Access.app`).
3. In the left sidebar, select **"login"** and **"My Certificates"**.
4. Find your certificate (should show a disclosure triangle; expanding it reveals the private key).
5. **Right-click the certificate** and choose **Export**.
6. Select **.p12** as the file format, choose a filename, and click **Save**.
7. Set a strong password when prompted (you’ll need this for CI/CD or GitHub secrets).

You now have a `.p12` file containing both your certificate and private key, suitable for code signing and automation workflows.

---

## Required GitHub Secrets

Add the following secrets to your GitHub repository (Settings � Secrets and variables � Actions):

### Code Signing Secrets

- `MACOS_CERTIFICATE`: Base64-encoded .p12 certificate file
- `MACOS_CERTIFICATE_PWD`: Password for the .p12 certificate
- `MACOS_SIGN_IDENTITY`: The certificate identity (e.g., "Developer ID Application: Your Name (TEAM_ID)")
- `KEYCHAIN_PASSWORD`: A secure password for the temporary keychain

### Notarization Secrets

- `MACOS_NOTARY_ISSUER_ID`: App Store Connect Issuer ID
- `MACOS_NOTARY_KEY_ID`: App Store Connect Key ID  
- `MACOS_NOTARY_KEY`: Base64-encoded .p8 private key content

## Setup Steps

### 1. Create Developer ID Certificate

1. Log into [Apple Developer Portal](https://developer.apple.com/)
2. Go to Certificates, Identifiers & Profiles
3. Create a new "Developer ID Application" certificate
4. Download the certificate and install it in Keychain Access
5. Export as .p12 file with a strong password

### 2. Convert Certificate to Base64

```bash
# Convert .p12 to base64 for GitHub secrets
base64 -i certificate.p12 | pbcopy
```

### 3. Create App Store Connect API Key

1. Go to [App Store Connect](https://appstoreconnect.apple.com/)
2. Navigate to Users and Access � Keys
3. Create a new API key with "Developer" role
4. Download the .p8 key file
5. Note the Key ID and Issuer ID

### 4. Convert API Key to Base64

```bash
# Convert .p8 key to base64
base64 -i AuthKey_XXXXXXXXXX.p8 | pbcopy
```

### 5. Find Certificate Identity

```bash
# List available code signing identities
security find-identity -v -p codesigning
```

Look for an identity like: `Developer ID Application: Your Name (TEAM_ID)`

## GitHub Secrets Configuration

Set the following secrets in your GitHub repository:

```
MACOS_CERTIFICATE=<base64-encoded-p12-content>
MACOS_CERTIFICATE_PWD=<p12-password>
MACOS_SIGN_IDENTITY=Developer ID Application: Your Name (TEAM_ID)
KEYCHAIN_PASSWORD=<secure-random-password>
MACOS_NOTARY_ISSUER_ID=<issuer-id-from-app-store-connect>
MACOS_NOTARY_KEY_ID=<key-id-from-app-store-connect>
MACOS_NOTARY_KEY=<base64-encoded-p8-content>
```

## Verification

Once configured, the release workflow will:

1. Import the Developer ID certificate into a temporary keychain
2. Sign all macOS binaries with the certificate and entitlements
3. Submit signed binaries to Apple's notarization service
4. Wait for notarization to complete
5. Clean up temporary files and keychain

## Troubleshooting

### Certificate Issues

- Ensure the certificate is valid and not expired
- Verify the certificate identity matches exactly
- Check that the certificate is a "Developer ID Application" type

### Notarization Issues

- Verify the App Store Connect API key has proper permissions
- Ensure the Issuer ID and Key ID are correct
- Check that the .p8 key file is valid and properly encoded

### Common Errors

- `errSecInternalComponent`: Certificate import failed - check password
- `Invalid credentials`: API key issues - verify IDs and key content
- `Notarization failed`: Binary may have issues - check entitlements

## Testing Locally

You can test code signing locally (requires macOS with certificates installed):

```bash
# Test code signing
codesign --sign "Developer ID Application: Your Name (TEAM_ID)" \
  --timestamp --options=runtime \
  --entitlements ./entitlements.plist \
  ./bin/rune

# Verify signature
codesign --verify --verbose ./bin/rune

# Test notarization (requires API key setup)
ditto -c -k --keepParent ./bin/rune rune.zip
xcrun notarytool submit rune.zip \
  --issuer "$ISSUER_ID" \
  --key-id "$KEY_ID" \
  --key "$KEY_FILE" \
  --wait
```

## Security Notes

- Never commit certificates or private keys to the repository
- Use strong passwords for certificate protection
- Rotate API keys periodically for security
- The temporary keychain is automatically cleaned up after each build
- All sensitive data is handled through GitHub encrypted secrets

## References

- [Apple Code Signing Guide](https://developer.apple.com/library/archive/documentation/Security/Conceptual/CodeSigningGuide/)
- [Apple Notarization Service](https://developer.apple.com/documentation/security/notarizing_macos_software_before_distribution)
- [App Store Connect API Keys](https://developer.apple.com/documentation/appstoreconnectapi/creating_api_keys_for_app_store_connect_api)