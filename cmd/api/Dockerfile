# --- build stage ---
FROM golang:1.27.4 AS build
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /flow ./cmd/api

# --- runtime stage ---
FROM gcr.io/distroless/static-debian13:nonroot

WORKDIR /app
COPY --from=build /flow /app/flow

USER nonroot:nonroot
EXPOSE 8000
ENTRYPOINT ["/app/flow"]
