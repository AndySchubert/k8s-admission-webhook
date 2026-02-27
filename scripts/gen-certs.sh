#!/usr/bin/env bash
set -euo pipefail

# Defaults (override by env vars, e.g. NAMESPACE=platform ./scripts/gen-certs.sh)
SERVICE_NAME="${SERVICE_NAME:-k8s-admission-webhook}"
NAMESPACE="${NAMESPACE:-default}"
DAYS="${DAYS:-3650}"

CERTS_DIR="${CERTS_DIR:-certs}"
CA_DIR="${CERTS_DIR}/ca"
SERVER_DIR="${CERTS_DIR}/server"
BUNDLE_DIR="${CERTS_DIR}/bundle"

mkdir -p "${CA_DIR}" "${SERVER_DIR}" "${BUNDLE_DIR}"

CA_KEY="${CA_DIR}/ca.key"
CA_CRT="${CA_DIR}/ca.crt"

SERVER_KEY="${SERVER_DIR}/tls.key"
SERVER_CSR="${SERVER_DIR}/tls.csr"
SERVER_CRT="${SERVER_DIR}/tls.crt"
CSR_CONF="${SERVER_DIR}/csr.conf"

CA_B64="${BUNDLE_DIR}/ca.crt.b64"

FQDN1="${SERVICE_NAME}"
FQDN2="${SERVICE_NAME}.${NAMESPACE}"
FQDN3="${SERVICE_NAME}.${NAMESPACE}.svc"

echo "==> Generating CA (if missing): ${CA_CRT}"
if [[ ! -f "${CA_KEY}" || ! -f "${CA_CRT}" ]]; then
  openssl genrsa -out "${CA_KEY}" 2048
  openssl req -x509 -new -nodes \
    -key "${CA_KEY}" \
    -subj "/CN=${SERVICE_NAME}-ca" \
    -days "${DAYS}" \
    -out "${CA_CRT}"
else
  echo "    CA already exists, skipping."
fi

echo "==> Writing CSR config with SANs: ${CSR_CONF}"
cat > "${CSR_CONF}" <<CONF
[req]
default_bits = 2048
prompt = no
default_md = sha256
req_extensions = req_ext
distinguished_name = dn

[dn]
CN = ${FQDN3}

[req_ext]
subjectAltName = @alt_names

[alt_names]
DNS.1 = ${FQDN1}
DNS.2 = ${FQDN2}
DNS.3 = ${FQDN3}
CONF

echo "==> Generating server key: ${SERVER_KEY}"
openssl genrsa -out "${SERVER_KEY}" 2048

echo "==> Generating CSR: ${SERVER_CSR}"
openssl req -new -key "${SERVER_KEY}" -out "${SERVER_CSR}" -config "${CSR_CONF}"

echo "==> Signing server certificate: ${SERVER_CRT}"
# CAcreateserial writes a .srl next to the CA cert/key; keep it in CA_DIR if present
# We run in CA_DIR to keep ca.srl there.
(
  cd "${CA_DIR}"
  openssl x509 -req \
    -in "../server/$(basename "${SERVER_CSR}")" \
    -CA "$(basename "${CA_CRT}")" \
    -CAkey "$(basename "${CA_KEY}")" \
    -CAcreateserial \
    -out "../server/$(basename "${SERVER_CRT}")" \
    -days "${DAYS}" \
    -extensions req_ext \
    -extfile "../server/$(basename "${CSR_CONF}")"
)

echo "==> Writing base64 CA bundle (single line): ${CA_B64}"
base64 < "${CA_CRT}" | tr -d '\n' > "${CA_B64}"
echo >> "${CA_B64}"

echo "==> Done."
echo "    CA:     ${CA_CRT}"
echo "    Server: ${SERVER_CRT}"
echo "    Key:    ${SERVER_KEY}"
echo "    Bundle: ${CA_B64}"
echo
echo "==> Verify SANs:"
openssl x509 -in "${SERVER_CRT}" -noout -text | sed -n '/Subject:/p;/Subject Alternative Name/,+1p'
