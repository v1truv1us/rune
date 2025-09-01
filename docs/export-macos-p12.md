# Exporting a macOS Code Signing Certificate as .p12 for GitHub Actions

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

---

## 1. Export the Certificate as a .p12 File (Keychain Access)

1. Open **Keychain Access** (`/Applications/Utilities/Keychain Access.app`).
2. In the left sidebar, select **login** under Keychains and **My Certificates** under Category.
3. Find your certificate (e.g., “Apple Development: Your Name (Team ID)”).
   - Make sure it has a disclosure triangle (▼) next to it—this means it includes the private key.
4. **Right-click** the certificate and select **Export “<certificate name>”...**
5. Choose a filename (e.g., `certificate.p12`), select **Personal Information Exchange (.p12)** as the format, and click **Save**.
6. When prompted, enter a strong password (this will be your `MACOS_CERTIFICATE_PWD`).
7. Enter your Mac user password to allow the export.

---

## 1a. Alternative: Export via Terminal (CLI)

### Quick Bash Script (for automation)

A ready-to-use script is available in this repo:

```bash
./scripts/export-apple-identity.sh "Apple Development: Your Name" output.p12
```

- The first argument is a substring of the certificate's Common Name (CN).
- The second argument is the output .p12 file path.
- The script will prompt for a password to protect the .p12 file.

See the script for more details and usage instructions.


Use these commands when you prefer the terminal. Make sure you are in a GUI session (not over SSH) and that your certificate in “My Certificates” shows an attached private key.

### Quick method: export identities (fastest if you only have one)

1) Unlock your login keychain if needed (enter your keychain password if prompted):

```bash
security unlock-keychain ~/Library/Keychains/login.keychain-db
```

2) Export identities to a `.p12` (you will be prompted to set a password that protects the `.p12`):

```bash
security export -k ~/Library/Keychains/login.keychain-db -t identities -f pkcs12 -o /path/to/certificate.p12
```

### Single identity only (robust selection by CN)

This script exports only the matching identity by comparing the certificate’s public key to its private key, ensuring the final `.p12` contains exactly the identity you want.

Copy-paste the entire block below into Terminal. It will prompt you for: a CN substring to match (e.g., “Rune” or the exact “Apple Development: Your Name”), an output path, and the export password to protect the final `.p12`.

```bash
bash -c '
set -euo pipefail

# Prompt for the identity CN substring to match (e.g., Apple Development: Your Name or a unique part like "Rune")
read -r -p "Enter part of the certificate Common Name to match (CN substring): " CN_FILTER

# Where to write the final .p12
read -r -p "Enter output .p12 path (e.g., ~/Desktop/rune-cert.p12): " OUTPUT_P12
OUTPUT_P12="${OUTPUT_P12/#\~/$HOME}"

# Read export password (to protect the final .p12)
read -r -s -p "Enter password to protect the exported .p12: " OUT_PASS; echo

# Resolve login keychain path
KEYCHAIN_PATH=$(security login-keychain | sed -e "s/[\"\']//g")

# Unlock keychain if locked (will prompt for keychain password)
if ! security show-keychain-info "$KEYCHAIN_PATH" >/dev/null 2>&1; then
  echo "Unlocking login keychain..."
  security unlock-keychain "$KEYCHAIN_PATH"
fi

# Temp working directory
WORKDIR="$(mktemp -d)"
trap "rm -rf \"$WORKDIR\"" EXIT

TEMP_P12="$WORKDIR/all_identities.p12"
TEMP_PASS="$(openssl rand -base64 24)"

echo "Exporting all identities from $KEYCHAIN_PATH to a temporary PKCS#12..."
security export -k "$KEYCHAIN_PATH" -t identities -f pkcs12 -P "$TEMP_PASS" -o "$TEMP_P12"

# Extract all leaf certs and all private keys
ALL_CERTS="$WORKDIR/allcerts.pem"
ALL_KEYS="$WORKDIR/allkeys.pem"

openssl pkcs12 -in "$TEMP_P12" -passin pass:"$TEMP_PASS" -clcerts -nokeys -out "$ALL_CERTS" >/dev/null 2>&1
openssl pkcs12 -in "$TEMP_P12" -passin pass:"$TEMP_PASS" -nocerts -out "$ALL_KEYS" -nodes >/dev/null 2>&1

# Split multi-PEM into individual files
split_pem() {
  local infile="$1" prefix="$2" type_marker="$3"
  awk -v prefix="$prefix" -v type="$type_marker" "
    /BEGIN/ && $0 ~ type {inblk=1; fn=sprintf(\"%s-%03d.pem\", prefix, ++i)}
    inblk {print > fn}
    /END/ && $0 ~ type {inblk=0}
  " "$infile"
}

split_pem "$ALL_CERTS" "$WORKDIR/cert" "BEGIN CERTIFICATE"
split_pem "$ALL_KEYS"  "$WORKDIR/key"  "BEGIN (ENCRYPTED )?PRIVATE KEY"

# Find matching cert by CN substring (case-insensitive)
MATCH_CERT=""
for cert in "$WORKDIR"/cert-*.pem; do
  subj="$(openssl x509 -noout -subject -in "$cert" 2>/dev/null || true)"
  if printf "%s" "$subj" | grep -qi -- "$CN_FILTER"; then
    MATCH_CERT="$cert"
    break
  fi
done

if [ -z "$MATCH_CERT" ]; then
  echo "Error: No certificate subject matched CN substring: $CN_FILTER"
  echo "Tip: Run: security find-identity -p codesigning -v | sed \"s/.*\\(\".*\"\\).*/\\1/\" to see available CNs"
  exit 1
fi

# Compute the public key of the matched cert
CERT_PUB="$WORKDIR/cert.pub"
openssl x509 -pubkey -noout -in "$MATCH_CERT" > "$CERT_PUB"

# Find the corresponding private key by comparing public keys
MATCH_KEY=""
for key in "$WORKDIR"/key-*.pem; do
  KEY_PUB="$WORKDIR/key.pub"
  if openssl pkey -in "$key" -pubout > "$KEY_PUB" 2>/dev/null; then
    if diff -q "$CERT_PUB" "$KEY_PUB" >/dev/null 2>&1; then
      MATCH_KEY="$key"
      break
    fi
  fi
done

if [ -z "$MATCH_KEY" ]; then
  echo "Error: Found a certificate match, but no corresponding private key."
  echo "This can happen if the private key is non-extractable or not present on this Mac."
  exit 1
fi

# Extract CN for friendly name
FRIENDLY_NAME="$(openssl x509 -noout -subject -in "$MATCH_CERT" | sed -n "s/^subject=.*CN=\([^,/]*\).*/\1/p" | head -1)"
[ -z "$FRIENDLY_NAME" ] && FRIENDLY_NAME="Exported Identity"

# Build final single-identity PKCS#12
openssl pkcs12 -export \
  -inkey "$MATCH_KEY" \
  -in "$MATCH_CERT" \
  -name "$FRIENDLY_NAME" \
  -out "$OUTPUT_P12" \
  -passout pass:"$OUT_PASS"

echo
echo "Success!"
echo "Exported single identity:"
echo "  Subject: $(openssl x509 -noout -subject -in "$MATCH_CERT")"
echo "  Output:  $OUTPUT_P12"
echo "Note: The .p12 is protected with the password you provided."
'
```

Tip: To see available identities and their names, you can run:

```bash
security find-identity -p codesigning -v
```

After the export finishes, continue with step 2 below to base64-encode for GitHub Actions.

---

## 2. Base64-Encode the .p12 File

```bash
base64 -i certificate.p12 | pbcopy
```
- This command copies the base64-encoded contents of your `.p12` file to your clipboard.

---

## 3. Add the Secrets to GitHub

1. Go to your repository on GitHub.
2. Click **Settings** → **Secrets and variables** → **Actions**.
3. Click **New repository secret** and add:
   - **Name:** `MACOS_CERTIFICATE`
     - **Value:** Paste the base64-encoded string from your clipboard.
   - **Name:** `MACOS_CERTIFICATE_PWD`
     - **Value:** The password you set when exporting the `.p12`.
   - **Name:** `KEYCHAIN_PASSWORD`
     - **Value:** A random strong password (generate with the command below or use a password manager).

Generate a strong password for `KEYCHAIN_PASSWORD`:
```bash
openssl rand -base64 32 | pbcopy
```

---

## 4. Re-run the Release Workflow

- Go to the **Actions** tab in your repo.
- Find the failed Release workflow for your tag.
- Click **Re-run jobs** (or push a new tag if you want to test a new release).

---

## Troubleshooting

- If you see “User interaction is not allowed” or errors when exporting over SSH: run in a GUI session and unlock your keychain first:
  ```bash
  security unlock-keychain ~/Library/Keychains/login.keychain-db
  ```
- If your Keychain item does not show a private key under the certificate: you cannot create a `.p12` on this Mac. Export it on the machine where the CSR/private key was originally created, or regenerate a new key and certificate.
- If multiple identities exist and the quick export produces a multi-identity `.p12`: use the single-identity script above with a specific CN substring.
- If passphrase prompts behave oddly: supply `-P "temp-pass"` with `security export` to avoid GUI prompts, or ensure the keychain is unlocked.

---

**You’re done!**

If you follow these steps, your workflow should be able to import the certificate and proceed with code signing.
