#!/usr/bin/env bash
set -euo pipefail

CERTS_DIR="${CERTS_DIR:-certs}"
CA_B64="${CERTS_DIR}/bundle/ca.crt.b64"

if [[ ! -f "${CA_B64}" ]]; then
  echo "Missing ${CA_B64}. Run: ./scripts/gen-certs.sh" >&2
  exit 1
fi

cat "${CA_B64}"
