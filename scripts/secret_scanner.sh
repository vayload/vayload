#!/usr/bin/env bash
set -euo pipefail

# Colors
RED="\033[0;31m"
YELLOW="\033[1;33m"
GREEN="\033[0;32m"
NC="\033[0m"

echo "Scanning files..."
echo

# Common secret patterns
PATTERNS=(
  # Generic API keys / tokens
  'api[_-]?key[[:space:]]*[:=][[:space:]]*["'\'']?[A-Za-z0-9_\-]{16,}'
  'secret[[:space:]]*[:=][[:space:]]*["'\'']?[A-Za-z0-9_\-]{16,}'
  'token[[:space:]]*[:=][[:space:]]*["'\'']?[A-Za-z0-9_\-]{16,}'

  # AWS
  'AKIA[0-9A-Z]{16}'
  'aws_secret_access_key[[:space:]]*[:=][[:space:]]*["'\'']?[A-Za-z0-9/+=]{40}'

  # Private keys
  '-----BEGIN (RSA|EC|OPENSSH|DSA)? ?PRIVATE KEY-----'

  # JWT
  'eyJ[A-Za-z0-9_\-]+\.[A-Za-z0-9_\-]+\.[A-Za-z0-9_\-]+'

  # URLs with embedded credentials
  'https?:\/\/[^[:space:]]+:[^[:space:]]+@'

  # Passwords
  'password[[:space:]]*[:=][[:space:]]*["'\'']?.{6,}'
  'pass[[:space:]]*[:=][[:space:]]*["'\'']?.{6,}'

  # Base64 blobs
  '(?:[A-Za-z0-9+\/]{32,}={0,2})'

  # URL-safe Base64
  '(?:[A-Za-z0-9_\-]{32,}={0,2})'

  # Hex secrets
  '\b[a-fA-F0-9]{32}\b'
  '\b[a-fA-F0-9]{40}\b'
  '\b[a-fA-F0-9]{64}\b'
  '\b[a-fA-F0-9]{128}\b'

  # Crypto material
  '-----BEGIN CERTIFICATE-----'
  '-----BEGIN PUBLIC KEY-----'
  '-----BEGIN ENCRYPTED PRIVATE KEY-----'
  '-----BEGIN PGP PRIVATE KEY BLOCK-----'

  # SSH keys
  'ssh-rsa[[:space:]]+[A-Za-z0-9+\/=]+'
  'ssh-ed25519[[:space:]]+[A-Za-z0-9+\/=]+'

  # Binary blobs
  '\\x[a-fA-F0-9]{2}(\\x[a-fA-F0-9]{2}){8,}'

  # OpenSSL-style keys
  '(?i)(aes|rsa|ecdsa|hmac|sha256|sha512)?[_-]?(key|secret|iv)[[:space:]]*[:=][[:space:]]*["'\'']?[a-fA-F0-9]{16,}'

  # PEM blocks
  '-----BEGIN [A-Z ]+-----'

  # Bearer tokens
  'Bearer[[:space:]]+[A-Za-z0-9\-._~+/]+=*'

  # Stripe
  'sk_live_[A-Za-z0-9]{24,}'
  'sk_test_[A-Za-z0-9]{24,}'

  # GitHub
  'gh[pousr]_[A-Za-z0-9]{20,}'

  # Slack
  'xox[baprs]-[A-Za-z0-9-]{10,}'

  # Google
  'AIza[0-9A-Za-z\-_]{35}'

  # Client secrets
  'client[_-]?secret[[:space:]]*[:=][[:space:]]*["'\'']?[A-Za-z0-9_\-\/+=]{16,}'
)

# Files to exclude
EXCLUDE_FILES=(
  "go.sum"
  "go.mod"
  "package-lock.json"
  "yarn.lock"
)

# Build exclude args
EXCLUDE_ARGS=()
for file in "${EXCLUDE_FILES[@]}"; do
  EXCLUDE_ARGS+=(--exclude="$file")
done

# Get tracked + untracked files
mapfile -t FILES < <(git ls-files -co --exclude-standard)

FOUND=0

for pattern in "${PATTERNS[@]}"; do
  echo -e "${YELLOW}Searching pattern:${NC} $pattern"
  echo "----------------------------------------"

  MATCHES=$(
    grep -RInE \
      --binary-files=without-match \
      --color=never \
      "${EXCLUDE_ARGS[@]}" \
      -- \
      "$pattern" \
      "${FILES[@]}" \
      2>/dev/null || true
  )

  if [[ -n "$MATCHES" ]]; then
    echo -e "${RED}Potential secrets found:${NC}"
    echo "$MATCHES"
    echo
    FOUND=1
  else
    echo "No matches."
    echo
  fi
done

echo "========================================"

if [[ $FOUND -eq 1 ]]; then
  echo -e "${RED}⚠ Potential secrets detected. Review before committing.${NC}"
  exit 1
else
  echo -e "${GREEN}✅ No obvious secrets detected.${NC}"
fi
