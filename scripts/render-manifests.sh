#!/usr/bin/env bash
set -euo pipefail

CERTS_DIR="${CERTS_DIR:-certs}"
CA_B64_FILE="${CERTS_DIR}/bundle/ca.crt.b64"

TEMPLATE="deploy/manifests.tmpl.yaml"
OUT="deploy/manifests.yaml"

if [[ ! -f "${CA_B64_FILE}" ]]; then
  echo "Missing ${CA_B64_FILE}. Run: ./scripts/gen-certs.sh" >&2
  exit 1
fi

if [[ ! -f "${TEMPLATE}" ]]; then
  echo "Missing ${TEMPLATE}" >&2
  exit 1
fi

CA_BUNDLE="$(cat "${CA_B64_FILE}")"

# Replace placeholder with the CA bundle
# (Use sed with a safe delimiter.)
sed "s|__CA_BUNDLE__|${CA_BUNDLE}|g" "${TEMPLATE}" > "${OUT}"

echo "Wrote ${OUT}"
