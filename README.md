# k8s-admission-webhook

A Kubernetes **Validating Admission Webhook** written in Go that enforces basic workload security and reliability policies.

## Enforced Policies

- ❌ Disallow `:latest` image tags  
- ❌ Require CPU and memory requests & limits  
- ✅ Allow properly defined workloads  

---

## What This Project Demonstrates

- Extending the Kubernetes control plane with a custom admission controller
- Handling `AdmissionReview` objects in Go
- TLS setup with a self-signed Certificate Authority
- Injecting `caBundle` into a `ValidatingWebhookConfiguration`
- Deployment to k3s (Rancher Desktop)
- Publishing container images to GHCR
- HTTPS readiness & liveness probes

---

## Architecture

```
Kubernetes API Server
        ↓
ValidatingWebhookConfiguration
        ↓
Service (ClusterIP, 443)
        ↓
Admission Webhook Pod (Go server, TLS on :8443)
```

The API server sends an `AdmissionReview` request to `/validate`.  
The webhook inspects the Pod spec and returns an allow/deny response.

---

## Project Structure

```
cmd/webhook/               # Application entrypoint
internal/admission/        # Validation logic
deploy/                    # Kubernetes manifests (templated)
certs/                     # Generated CA + server certs
scripts/
  gen-certs.sh             # Generate self-signed CA + server cert
  render-manifests.sh      # Inject CA bundle into webhook config
```

---

## Deployment (k3s / Rancher Desktop)

### 1. Generate Certificates

```bash
./scripts/gen-certs.sh
```

### 2. Create Namespace and TLS Secret

```bash
kubectl create namespace platform-system 2>/dev/null || true

kubectl create secret tls webhook-tls \
  -n platform-system \
  --cert=certs/server/tls.crt \
  --key=certs/server/tls.key
```

### 3. Render Manifests (Inject CA Bundle)

```bash
./scripts/render-manifests.sh
```

### 4. Deploy

```bash
kubectl apply -f deploy/manifests.yaml
kubectl rollout status deploy/k8s-admission-webhook -n platform-system
```

---

## Container Image

Published to GHCR:

```
ghcr.io/andyschubert/k8s-admission-webhook:dev
```

Build and push:

```bash
docker build -t ghcr.io/andyschubert/k8s-admission-webhook:dev .
docker push ghcr.io/andyschubert/k8s-admission-webhook:dev
```

---

## Health Endpoints

The webhook exposes:

- `GET /healthz`
- `GET /readyz`

Both served over HTTPS on port `8443`.

---

## Testing

Create a test namespace:

```bash
kubectl create namespace test
```

---

### ✅ Allowed Pod

```bash
cat <<'EOF' | kubectl apply -n test -f -
apiVersion: v1
kind: Pod
metadata:
  name: good-pod
spec:
  containers:
    - name: app
      image: nginx:1.25
      resources:
        requests:
          cpu: "100m"
          memory: "128Mi"
        limits:
          cpu: "200m"
          memory: "256Mi"
EOF
```

---

### ❌ Rejected: latest tag

```bash
cat <<'EOF' | kubectl apply -n test -f -
apiVersion: v1
kind: Pod
metadata:
  name: bad-pod-latest
spec:
  containers:
    - name: app
      image: nginx:latest
      resources:
        requests:
          cpu: "100m"
          memory: "128Mi"
        limits:
          cpu: "200m"
          memory: "256Mi"
EOF
```

Expected error:

```
denied the request: container app uses disallowed image tag (latest or no tag)
```

---

### ❌ Rejected: missing resources

```bash
cat <<'EOF' | kubectl apply -n test -f -
apiVersion: v1
kind: Pod
metadata:
  name: bad-pod-no-resources
spec:
  containers:
    - name: app
      image: nginx:1.25
EOF
```

Expected error:

```
denied the request: container app must define cpu and memory requests and limits
```

---

## Observability

Tail webhook logs:

```bash
kubectl logs -n platform-system -l app=k8s-admission-webhook -f
```

Example:

```
Admission request UID=... Namespace=test Name=bad-pod-latest Operation=CREATE
```

---

## Future Improvements

- Configurable policies via ConfigMap
- Support for Deployments & StatefulSets
- Prometheus metrics endpoint
- Mutation webhook support
- Integration tests using kind
- Replace self-signed certificates with cert-manager

---

## Why This Project?

This repository demonstrates:

- Deep understanding of Kubernetes control plane extensions
- Secure TLS bootstrapping
- Cluster-level policy enforcement
- Practical platform engineering patterns in Go

---

## License

MIT