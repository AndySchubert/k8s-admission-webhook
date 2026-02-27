# Build stage
FROM golang:1.22 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /out/webhook ./cmd/webhook

# Runtime stage
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=build /out/webhook /webhook
USER nonroot:nonroot
EXPOSE 8443
ENTRYPOINT ["/webhook"]