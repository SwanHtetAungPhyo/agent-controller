FROM golang:1.25.3-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata
WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download && go mod verify

COPY .. .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/main cmd/main.go


FROM gcr.io/distroless/static-debian12

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8081

LABEL maintainer="kanios.it.com"
LABEL version="VERSION-1.10"
LABEL description="Kainos"
LABEL pembroke.com='VERSION-1.10'

USER nonroot:nonroot

CMD ["./main"]